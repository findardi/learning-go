package user

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	FirstName string
	LastName  string
	UserName  string
	CreatedAt time.Time
}

// constructor method
func New(firstname, lastname, username string) (*User, error) {
	if firstname == "" || lastname == "" || username == "" {
		return nil, errors.New("firstname or lastname or username cant be empty")
	}
	return &User{
		FirstName: firstname,
		LastName:  lastname,
		UserName:  username,
		CreatedAt: time.Now(),
	}, nil
}

func (u *User) PrintOut() {
	fmt.Println(u.FirstName, u.LastName, u.UserName, u.CreatedAt)
}

func (u *User) ClearUsername() {
	u.UserName = ""
}
