package domain

// represent core business logics for user domain
type UserCore struct {
	Store UserStore
}

// TODO:
// - field validation especially email, only allows certain characters for password
// - encrypt password
// At the moment the method will just encrypt the password
// and save the user to db
func (core *UserCore) AddUser(username, firstName, lastName, email, password string) (User, error) {
	return User{}, nil
}

// TODO:
// - return
func (core *UserCore) GetUserByID(userID string) (User, error) {
	return User{}, nil
}

type UserStore interface {
}
