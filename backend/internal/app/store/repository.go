package store

import "github.com/MrSedan/restapigoown/backend/internal/app/model"

// UserRepository is the interface for db
type UserRepository interface {
	Create(*model.User) error
	GetProfile(string) (*model.Profile, error)
	CreateProfile(*model.User) error
	// Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	FindByID(string) (*model.User, error)
	ClaimToken(*model.User, string)
	GetToken(string) (string, error)
	EditAbout(int, string) error
	EditPass(*model.User) error
}
