package domain

type UserCore struct {
	Store UserStore
}

func (core *UserCore) AddUser(username, firstName, lastName, email, password string) (User, error) {
	return User{}, nil
}

func (core *UserCore) GetUserByID(userID string) (User, error) {
	return User{}, nil
}

type UserStore interface {
}
