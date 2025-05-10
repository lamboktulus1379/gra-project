package user

import "time"

// User represents the user entity with all its attributes
type User struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user instance with current time for created/updated fields
func NewUser(firstName, lastName, email, password string) *User {
	now := time.Now()
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate returns true if the user data is valid
func (u *User) Validate() bool {
	return u.FirstName != "" && u.LastName != "" && u.Email != "" && u.Password != ""
}
