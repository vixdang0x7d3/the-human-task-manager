package http

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/keighl/postmark"
	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/models"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/components"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/pages"
)

// registerAuthRoutes handle unauthenticated routes such as signup, login and oauth2.
// unauthenticated routes has prefix '/u'
func (s *Server) registerAuthRoutes(r *echo.Group) {
	r.GET("/login", s.handleLoginShow)
	r.GET("/signup", s.handleSignupShow)

	r.POST("/signup", s.handleSignup)
	r.POST("/login/email", s.handleLoginEmail)
	r.POST("/login/password", s.handleLoginPassword)
	r.GET("/login/success", s.handleLoginSuccess)

	r.GET("/password-reset-send", s.handleResetPasswordSendEmailShow)
	r.POST("/password-reset-send", s.handleResetPasswordSendEmail)

	r.GET("/password-reset/:token", s.handlePasswordResetShow)
	r.POST("/password-reset/:token", s.handlePasswordReset)
}

func (s *Server) handleLoginShow(c echo.Context) error {
	return render(c, http.StatusOK, pages.LoginEmail("/u/login/email", "/u/signup"))
}

func (s *Server) handleSignupShow(c echo.Context) error {
	return render(c, http.StatusOK, pages.Signup("/u/signup"))
}

func (s *Server) handleSignup(c echo.Context) error {
	type formValues struct {
		Username  string `form:"username" validate:"required"`
		FirstName string `form:"first_name" validate:"required"`
		LastName  string `form:"last_name" validate:"required"`
		Email     string `form:"email" validate:"required,email"`
		Password  string `form:"password" validate:"required"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return render(c, http.StatusBadRequest, components.AlertError("internal error"))
	}

	if err := c.Validate(form); err != nil {
		c.Logger().Error("form validation error", err)
		return render(c, http.StatusBadRequest, components.AlertError("invalid signup info"))
	}

	user, err := s.UserService.Create(c.Request().Context(), domain.CreateUserCmd(form))
	if err != nil {
		return render(c, http.StatusBadRequest, components.AlertError(domain.ErrorMessage(err)))
	}

	m := models.UserView{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return render(c, http.StatusOK, components.UserInfoPostSignup(m))
}

func (s *Server) handleLoginEmail(c echo.Context) error {
	type formValues struct {
		Email string `form:"email" validate:"required,email"`
	}

	form := formValues{}

	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return render(c, http.StatusBadRequest, components.AlertError("invalid form data"))
	}

	if err := c.Validate(form); err != nil {
		c.Logger().Error("form validation error", err)
		return c.HTML(http.StatusBadRequest, "invalid email")
	}

	user, err := s.UserService.ByEmail(c.Request().Context(), form.Email)
	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
	}

	m := models.UserView{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return render(c, http.StatusOK, components.LoginPassword(m, "/u/login/password", "/u/password-reset-send"))
}

func (s *Server) handleLoginPassword(c echo.Context) error {

	type formValues struct {
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid form data")
	}

	if err := c.Validate(form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid password")
	}

	user, err := s.UserService.ByEmailWithPassword(c.Request().Context(), form.Email, form.Password)
	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
	}

	if err := s.sessions.RenewToken(c.Request().Context()); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusInternalServerError, "internal error")
	}
	s.sessions.Put(c.Request().Context(), "userID", user.ID.String())

	return render(c, http.StatusOK, components.LoginAlertRedirect("login successful!", "/u/login/success"))
}

func (s *Server) handleLoginSuccess(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *Server) handleResetPasswordSendEmailShow(c echo.Context) error {
	return render(c, http.StatusOK, pages.PasswordResetEmail("/u/password-reset-send"))
}

func (s *Server) handleResetPasswordSendEmail(c echo.Context) error {

	emailClient := postmark.NewClient(
		os.Getenv("POSTMARK_SERVER_TOKEN"),
		os.Getenv("POSTMARK_ACCOUNT_TOKEN"),
	)

	type formValues struct {
		Email string `form:"email" validate:"required,email"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid form data")
	}

	if err := c.Validate(form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid email format")
	}

	user, err := s.UserService.ByEmail(c.Request().Context(), form.Email)
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID:
			return render(c, http.StatusBadRequest, components.AlertError("invalid email"))
		case domain.ENOTFOUND:
			return render(c, http.StatusBadRequest, components.AlertError("email not found"))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return render(c, http.StatusBadRequest, components.AlertError("internal error"))
		}
	}

	token, err := generateToken()
	if err != nil {
		return render(c, http.StatusBadRequest, components.AlertError("internal error"))
	}

	resetLink := fmt.Sprintf("%s/password-reset/%s", os.Getenv("BASE_URL"), token)

	email := postmark.Email{
		From:     "n21dcat065@student.ptithcm.edu.vn",
		To:       user.Email,
		Subject:  "Password Reset Request",
		HtmlBody: getPasswordResetEmailHTML(user.Username, resetLink),
		TextBody: getPasswordResetEmailText(user.Username, resetLink),
		Tag:      "password-reset",
	}

	_, err = emailClient.SendEmail(email)
	if err != nil {
		c.Logger().Error(err)
		return render(c, http.StatusBadRequest, components.AlertError("send email failed"))
	}

	expiresAt := time.Now().Add(1 * time.Hour)

	s.sessions.Put(c.Request().Context(), "password_reset_expires", expiresAt.Format(time.RFC3339))
	s.sessions.Put(c.Request().Context(), "password_reset_user_id", user.ID.String())

	return render(c, http.StatusOK, components.AlertSuccess("check email for password reset"))
}

func getPasswordResetEmailHTML(name, resetLink string) string {
	return fmt.Sprintf(`
        <h2>Hello %s,</h2>
        <p>We received a request to reset your password. If you didn't make this request, you can ignore this email.</p>
        <p>To reset your password, click the link below:</p>
        <p><a href="%s">Reset Password</a></p>
        <p>This link will expire in 1 hour.</p>
        <p>If you're having trouble clicking the link, copy and paste this URL into your browser:</p>
        <p>%s</p>
    `, name, resetLink, resetLink)
}

func getPasswordResetEmailText(name, resetLink string) string {
	return fmt.Sprintf(`
        Hello %s,

        We received a request to reset your password. If you didn't make this request, you can ignore this email.

        To reset your password, copy and paste this URL into your browser:
        %s

        This link will expire in 1 hour.
    `, name, resetLink)
}

func (s *Server) handlePasswordResetShow(c echo.Context) error {

	token := c.Param("token")

	expiresAtStr := s.sessions.GetString(c.Request().Context(), "password_reset_expires")

	expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil {

		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "expired or invalid reset link")
	}

	if time.Now().After(expiresAt) || expiresAt.IsZero() {
		return c.HTML(http.StatusBadRequest, "expired or invalid reset link")
	}

	url := fmt.Sprintf("/u/password-reset/%s", token)

	return render(c, http.StatusOK, pages.PasswordResetForm(url))
}

func (s *Server) handlePasswordReset(c echo.Context) error {

	type formValues struct {
		Password string `form:"password" validate:"required"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid password")
	}

	if err := c.Validate(form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "password can't be empty")
	}

	expiresAtStr := s.sessions.GetString(c.Request().Context(), "password_reset_expires")
	expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil {

		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "expired or invalid reset link")
	}

	if time.Now().After(expiresAt) || expiresAt.IsZero() {
		return c.HTML(http.StatusBadRequest, "expired or invalid reset link")
	}

	userID := s.sessions.GetString(c.Request().Context(), "password_reset_user_id")
	if userID == "" {
		return c.HTML(http.StatusBadRequest, "invalid reset link")
	}

	user, err := s.UserService.ByID(c.Request().Context(), userID)
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID:
			return render(c, http.StatusBadRequest, components.AlertError("invalid reset link"))
		case domain.ENOTFOUND:
			return render(c, http.StatusBadRequest, components.AlertError("invalid reset link"))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return render(c, http.StatusBadRequest, components.AlertError("internal error"))
		}
	}

	_, err = s.UserService.Update(domain.NewContextWithUser(c.Request().Context(), &user), domain.UpdateUserCmd{
		Password: &form.Password,
	})

	if err != nil {
		return render(c, http.StatusInternalServerError, components.AlertError("internal error"))
	}

	return render(c, http.StatusOK, components.AlertSuccess("new password saved"))
}
