package user

// Repository defines the interface for user data access
type Repository interface {
	Save(user *User) error
	FindByEmail(email string) (*User, error)
}
