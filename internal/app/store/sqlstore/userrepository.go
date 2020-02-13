package sqlstore

import (
	"database/sql"

	"github.com/MrSedan/restapigoown/internal/app/model"
	"github.com/MrSedan/restapigoown/internal/app/store"
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
	return r.store.db.DB().QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

// FindByEmail finding user in DB by email-addr
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.DB().QueryRow(
		"SELECT email, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

// ClaimToken set a new token to db
func (r *UserRepository) ClaimToken(u *model.User, token string) {
	r.store.db.Model(&model.User{}).Where("email=?", u.Email).Update("jwt_token", token)
}

// GetToken is checking token in db
func (r *UserRepository) GetToken(token string) error {
	var tok = ""
	if err := r.store.db.DB().QueryRow(
		"SELECT jwt_token FROM users WHERE jwt_token=$1",
		token,
	).Scan(&tok); err != nil || token == "" {
		return store.ErrNotValidToken
	}
	return nil
}
