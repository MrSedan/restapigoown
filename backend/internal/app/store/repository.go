package store

import "github.com/MrSedan/restapigoown/backend/internal/app/model"

// UserRepository is the interface for db
type UserRepository interface {
	Create(*model.User) error
	GetProfile(string) (*model.Profile, error)
	CreateProfile(*model.User) error
	FindByEmail(string) (*model.User, error)
	FindByID(string) (*model.User, error)
	FindByNick(string) (*model.User, error)
	ClaimToken(*model.User, string)
	CheckToken(string) (string, error)
	CompareToken(*model.User, string) bool
	GetToken(int) (string, error)
	//?Profile Editing
	EditAbout(int, string) error
	EditFirstName(int, string) error
	EditLastName(int, string) error
	//?
	EditPass(*model.User) error
	GetAllUsers() ([]*model.User, error)
}
