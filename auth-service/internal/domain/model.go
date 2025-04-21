package domain

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRepository interface {
	Create(user User) error
	GetByUsername(username string) (*User, error)
}
