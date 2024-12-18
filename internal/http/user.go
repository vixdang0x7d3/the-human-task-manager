package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/domain"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/generic"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/models"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/components"
	"github.com/vixdang0x7d3/the-human-task-manager/internal/http/templates/pages"
)

func (s *Server) registerUserRoutes(r *echo.Group) {
	r.GET("/index", s.handleIndexShow)
	r.GET("/profile", s.handleProfileShow)

	r.POST("/change-info", s.handleChangeProfileInfo)
	r.POST("/change-info-save", s.handleSaveProfileInfoChange)
	r.POST("/change-email", s.handleChangeProfileEmail)
	r.POST("/change-email-save", s.handleSaveProfileEmailChange)
	r.POST("/change-password", s.handleChangeProfilePassword)
	r.POST("/change-password-save", s.handleSaveprofilePasswordChange)

	r.DELETE("/delete-account", s.handleDeleteUser)
	r.DELETE("/logout", s.handleLogout)
}

func (s *Server) handleIndexShow(c echo.Context) error {

	state := "started"
	taskItems, total, err := s.TaskItemService.Find(c.Request().Context(), domain.TaskItemFilter{
		State:  &state,
		Offset: 0,
		Limit:  3,
	})
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EINVALID, domain.EUNAUTHORIZED:
			return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
		case domain.EINTERNAL:
			c.Logger().Error(domain.ErrorMessage(err))
			return c.HTML(http.StatusInternalServerError, "internal error")
		}
	}
	_ = total

	user := domain.UserFromContext(c.Request().Context())
	if user == nil {
		c.Logger().Error("no user in current session")
		return echo.NewHTTPError(http.StatusBadRequest, "no user in login session")
	}

	m := models.UserView{
		FirstName: user.FirstName,
	}

	return render(c, http.StatusOK, pages.Index(m,
		generic.Map(taskItems, toTaskItemView),
		"/logout"))
}

func (s *Server) handleProfileShow(c echo.Context) error {

	user := domain.UserFromContext(c.Request().Context())

	m := models.UserView{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		ID:        user.ID.String(),
	}

	return render(c, http.StatusOK, pages.Profile(m, "/change-info", "/change-email", "/change-password", "/delete-account", "/logout"))
}

func (s *Server) handleChangeProfileInfo(c echo.Context) error {

	type formValues struct {
		Username  string `form:"username" validate:"required"`
		FirstName string `form:"first_name" validate:"required"`
		LastName  string `form:"last_name" validate:"required"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid form data")
	}

	if err := c.Validate(form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid value")
	}

	m := models.UserView{
		Username:  form.Username,
		FirstName: form.FirstName,
		LastName:  form.LastName,
	}

	return render(c, http.StatusOK, components.ChangeInfoForm(m, "/change-info-save"))
}

func (s *Server) handleSaveProfileInfoChange(c echo.Context) error {
	type formValues struct {
		Username  string `form:"username" validate:"required"`
		FirstName string `form:"first_name" validate:"required"`
		LastName  string `form:"last_name" validate:"required"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid form data")
	}

	if err := c.Validate(&form); err != nil {
		var userNameMessage, firstNameMessage, lastNameMessage string
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			for _, fieldErr := range validationErrors {
				switch fieldErr.Field() {
				case "Username":
					userNameMessage = "User Name is required."
				case "FirstName":
					firstNameMessage = "First Name is required."
				case "LastName":
					lastNameMessage = "Last Name is required."
				}
			}
		}
		return render(c, http.StatusBadRequest, components.InfoErrorMessage(userNameMessage, firstNameMessage, lastNameMessage))
	}

	// save info into database_______________________!!!
	cmd := domain.UpdateUserCmd{
		Username:  &form.Username,
		FirstName: &form.FirstName,
		LastName:  &form.LastName,
	}

	user, err := s.UserService.Update(c.Request().Context(), cmd)
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EUNAUTHORIZED:
			return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return c.HTML(http.StatusInternalServerError, "internal error")
		}
	}

	c.Logger().Info("update user successful ", user.ID)
	s.sessions.Put(c.Request().Context(), "userID", user.ID.String())

	m := models.UserView{
		Username:  form.Username,
		FirstName: form.FirstName,
		LastName:  form.LastName,
	}

	return render(c, http.StatusOK, components.SavedInfoForm(m, "change-info"))
}

func (s *Server) handleChangeProfileEmail(c echo.Context) error {
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
		return c.HTML(http.StatusBadRequest, "invalid value")
	}

	m := models.UserView{
		Email: form.Email,
	}

	return render(c, http.StatusOK, components.ChangeEmailForm(m, "/change-email-save"))
}

func (s *Server) handleSaveProfileEmailChange(c echo.Context) error {
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

	// save info into database_______________________!!!
	cmd := domain.UpdateUserCmd{
		Email: &form.Email,
	}

	user, err := s.UserService.Update(c.Request().Context(), cmd)
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EUNAUTHORIZED:
			return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return c.HTML(http.StatusInternalServerError, "internal error")
		}
	}

	c.Logger().Info("update user successful ", user.ID)
	s.sessions.Put(c.Request().Context(), "userID", user.ID.String())

	m := models.UserView{
		Email: form.Email,
	}

	return render(c, http.StatusOK, components.SavedEmailForm(m, "change-email"))
}

func (s *Server) handleChangeProfilePassword(c echo.Context) error {
	return render(c, http.StatusOK, components.ChangePasswordForm("/change-password-save"))
}

func (s *Server) handleSaveprofilePasswordChange(c echo.Context) error {
	type formValues struct {
		CurrentPassword string `form:"current-password" validate:"required"`
		NewPassword     string `form:"new-password" validate:"required"`
	}

	form := formValues{}
	if err := c.Bind(&form); err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusBadRequest, "invalid form data")
	}

	if err := c.Validate(form); err != nil {
		c.Logger().Error(err)
		if err := c.Validate(&form); err != nil {
			var currentPwMessage, newPwMessage string
			validationErrors, ok := err.(validator.ValidationErrors)
			if ok {
				for _, fieldErr := range validationErrors {
					switch fieldErr.Field() {
					case "CurrentPassword":
						currentPwMessage = "Current password is required."
					case "NewPassword":
						newPwMessage = "New password is required."
					}
				}
			}
			return render(c, http.StatusBadRequest, components.PassWordErrorMessage(currentPwMessage, newPwMessage))
		}
		return c.HTML(http.StatusBadRequest, "invalid password")
	}

	// check current password_______________________!!!
	_, err := s.UserService.WithPassword(c.Request().Context(), form.CurrentPassword)
	if err != nil {
		return render(c, http.StatusBadRequest, components.PassWordErrorMessage("wrong current password", ""))
	}

	// save info into database_______________________!!!
	cmd := domain.UpdateUserCmd{
		Password: &form.NewPassword,
	}
	user, err := s.UserService.Update(c.Request().Context(), cmd)
	if err != nil {
		switch domain.ErrorCode(err) {
		case domain.EUNAUTHORIZED:
			return c.HTML(http.StatusBadRequest, domain.ErrorMessage(err))
		case domain.EINTERNAL:
			c.Logger().Error(err)
			return c.HTML(http.StatusInternalServerError, "internal error")
		}
	}
	c.Logger().Info("update user successful ", user.ID)
	s.sessions.Put(c.Request().Context(), "userID", user.ID.String())

	return render(c, http.StatusOK, components.SavedPasswordForm("/change-password"))
}

func (s *Server) handleDeleteUser(c echo.Context) error {
	_, err := s.UserService.Delete(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusInternalServerError, "internal error")
	}

	err = s.sessions.Destroy(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusInternalServerError, "internal error")
	}

	c.Response().Header().Set("HX-Redirect", "/u/login")
	return c.NoContent(http.StatusOK)
}

// TODO: flash message support
func (s *Server) handleLogout(c echo.Context) error {
	err := s.sessions.Destroy(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return c.HTML(http.StatusInternalServerError, "internal error")
	}

	c.Response().Header().Set("HX-Redirect", "/u/login")
	return c.NoContent(http.StatusOK)
}
