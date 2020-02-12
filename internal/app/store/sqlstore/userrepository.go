package sqlstore

import (
	"github.com/MrSedan/restapigoown/internal/app/model"
)

// UserRepository a struct with store
type UserRepository struct {
	store *Store
}

// Create a user row in db
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	if err := u.BeforeCreate(); err != nil {
		return err
	}
	return r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}
