package store

import "github.com/MrSedan/restapigoown/internal/app/model"

// UserRepository is the interface for db
type UserRepository interface {
	Create(*model.User) error
	GetProfile(string) *model.Profile
	CreateProfile(*model.User) error
	// Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	ClaimToken(*model.User, string)
	GetToken(string) (string, error)
	EditAbout(int, string) error
}
