package store

import "github.com/MrSedan/restapigoown/internal/app/model"

// UserRepository is the interface for db
type UserRepository interface {
	Create(*model.User) error
	// Find(int) (*model.User, error)
	// FindByEmail(string) (*model.User, error)
}
