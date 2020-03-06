package sqlstore

import (
	"database/sql"

	"github.com/MrSedan/restapigoown/backend/internal/app/model"
	"github.com/MrSedan/restapigoown/backend/internal/app/store"
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
		"INSERT INTO users (first_name, last_name, email, encrypted_password) VALUES ($1, $2, $3, $4) RETURNING id",
		u.FirstName,
		u.LastName,
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

// CreateProfile ctreating a profile
func (r *UserRepository) CreateProfile(u *model.User) error {
	return r.store.db.DB().QueryRow(
		"INSERT INTO profiles (first_name, last_name, user_email , user_id) VALUES ((SELECT first_name FROM users WHERE email=$1), (SELECT last_name FROM users WHERE email=$1), (SELECT email FROM users WHERE email=$1), (SELECT id FROM users WHERE email=$1)) RETURNING user_email",
		u.Email,
	).Scan(&u.Email)
}

// FindByEmail finding user in DB by email-addr
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.DB().QueryRow(
		"SELECT first_name, last_name, id, email, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(
		&u.FirstName,
		&u.LastName,
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

//FindByID searching user by id
func (r *UserRepository) FindByID(id string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.DB().QueryRow(
		"SELECT first_name, last_name, id, email, encrypted_password FROM users WHERE id = $1",
		id,
	).Scan(
		&u.FirstName,
		&u.LastName,
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
		id  string
	)
	if err := r.store.db.DB().QueryRow(
		"SELECT jwt_token, id FROM users WHERE jwt_token=$1",
		token,
	).Scan(&tok, &id); err != nil || tok == "" {
		return "", store.ErrNotValidToken
	}
	return id, nil
}

// GetProfile Getting profile
func (r *UserRepository) GetProfile(id string) (*model.Profile, error) {
	type sc struct {
		FirstName string
		LastName  string
		ID        int
		UserEmail string
		About     string
	}
	u := &model.User{}
	sk := &sc{}
	if err := r.store.db.DB().QueryRow(
		"SELECT id FROM users WHERE id=$1",
		id,
	).Scan(&u.ID); err == sql.ErrNoRows {
		return nil, store.ErrRecordNotFound
	}
	r.store.db.DB().QueryRow(
		"SELECT first_name, last_name, about FROM profiles WHERE user_id=$1",
		id,
	).Scan(&sk.FirstName, &sk.LastName, &sk.About)
	pr := &model.Profile{
		FirstName: sk.FirstName,
		LastName:  sk.LastName,
		UserID:    u.ID,
		About:     sk.About,
	}
	return pr, nil
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
		"UPDATE users SET encrypted_password=$1 WHERE id=$2",
		u.EncryptedPassword,
		u.ID,
	)
	return err
}
