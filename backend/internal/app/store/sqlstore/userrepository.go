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

// CreateProfile ctreating a profile
func (r *UserRepository) CreateProfile(u *model.User) error {
	return r.store.db.DB().QueryRow(
		"INSERT INTO profiles (user_email , user_id) VALUES ((SELECT email FROM users WHERE email=$1), (SELECT id FROM users WHERE email=$1)) RETURNING user_email",
		u.Email,
	).Scan(&u.Email)
}

// FindByEmail finding user in DB by email-addr
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.DB().QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
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
func (r *UserRepository) GetToken(token string) (string, error) {
	var (
		tok string
		em  string
	)
	if err := r.store.db.DB().QueryRow(
		"SELECT jwt_token, email FROM users WHERE jwt_token=$1",
		token,
	).Scan(&tok, &em); err != nil || tok == "" {
		return "", store.ErrNotValidToken
	}
	return em, nil
}

// GetProfile Getting profile
func (r *UserRepository) GetProfile(email string) *model.Profile {
	type sc struct {
		ID        int
		UserEmail string
		About     string
	}
	u := &model.User{}
	sk := &sc{}
	r.store.db.DB().QueryRow(
		"SELECT id FROM users WHERE email=$1",
		email,
	).Scan(&u.ID)
	r.store.db.DB().QueryRow(
		"SELECT user_email, about FROM profiles WHERE user_email=$1",
		email,
	).Scan(&sk.UserEmail, &sk.About)
	pr := &model.Profile{
		UserID:    u.ID,
		UserEmail: sk.UserEmail,
		About:     sk.About,
	}
	return pr
}

//EditAbout editing about
func (r *UserRepository) EditAbout(id int, about string) error {
	_, err := r.store.db.DB().Exec(
		"UPDATE profiles SET about=$1 WHERE user_id=$2",
		about,
		id,
	)
	return err
}

//EditPass changing a password to new
func (r *UserRepository) EditPass(u *model.User) error {
	if err := u.BeforeCreate(); err != nil {
		return err
	}
	_, err := r.store.db.DB().Exec(
		"UPDATE users SET encrypted_password=$1 WHERE email=$2",
		u.EncryptedPassword,
		u.Email,
	)
	return err
}
