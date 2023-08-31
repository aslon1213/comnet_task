package models

import (
	"fmt"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Login     string    `json:"username"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) String() string {

	return fmt.Sprintf("Username: %s, Password: %s, FirstName: %s, LastName: %s, Email: %s", u.
		Login, u.Password, u.Name, u.Email)

}

func (u *User) Validate() error {

	if u.Login == "" {
		return fmt.Errorf("username is required")
	}

	if u.Password == "" {
		return fmt.Errorf("password is required")
	}

	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	if u.Age == 0 {
		return fmt.Errorf("age is required")
	}

	return nil
}
