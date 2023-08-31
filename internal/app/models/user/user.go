package user

import "fmt"

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

func (u *User) String() string {

	return fmt.Sprintf("Username: %s, Password: %s, FirstName: %s, LastName: %s, Email: %s", u.Username, u.Password, u.FirstName, u.LastName, u.Email)

}

func (u *User) Validate() error {

	if u.Username == "" {
		return fmt.Errorf("username is required")
	}

	if u.Password == "" {
		return fmt.Errorf("password is required")
	}

	if u.FirstName == "" {
		return fmt.Errorf("firstName is required")
	}

	if u.LastName == "" {
		return fmt.Errorf("lastName is required")
	}

	if u.Email == "" {
		return fmt.Errorf("email is required")
	}

	return nil
}
