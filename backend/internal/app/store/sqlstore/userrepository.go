package sqlstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/MrSedan/restapigoown/backend/internal/app/model"
	"github.com/MrSedan/restapigoown/backend/internal/app/store"
)

// UserRepository a struct with store
type UserRepository struct {
	store *Store
}

var ctx context.Context

// Create a user row in db
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	if err := u.BeforeCreate(); err != nil {
		return err
	}
	return r.store.db.DB().QueryRow(
		"INSERT INTO users (user_name, email, encrypted_password) VALUES ($1, $2, $3) RETURNING id",
		u.UserName,
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
		"SELECT user_name, id, email, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(
		&u.UserName,
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
		"SELECT user_name, id, email, encrypted_password FROM users WHERE id = $1",
		id,
	).Scan(
		&u.UserName,
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

//FindByNick searching user by username
func (r *UserRepository) FindByNick(username string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.DB().QueryRow(
		"SELECT user_name, id, email, encrypted_password FROM users WHERE user_name = $1",
		username,
	).Scan(
		&u.UserName,
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

// CheckToken is checking token in db
func (r *UserRepository) CheckToken(token string) (string, error) {
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

// GetToken getting token from db by ID
func (r *UserRepository) GetToken(id int) (string, error) {
	var token string
	err := r.store.db.DB().QueryRow(
		"SELECT jwt_token FROM users WHERE id=$1",
		id,
	).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

// GetAllUsers getting all users from DB
func (r *UserRepository) GetAllUsers() ([]*model.User, error) {
	users := make([]*model.User, 0)
	rows, err := r.store.db.DB().Query("SELECT user_name, id FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u = &model.User{}
		if err := rows.Scan(&u.UserName, &u.ID); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if len(users) > 0 {
		return users, nil
	}
	return nil, errors.New("No users in db")
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

//?Profile Editing

//EditAbout editing about
func (r *UserRepository) EditAbout(id int, about string) error {
	_, err := r.store.db.DB().Exec(
		"UPDATE profiles SET about=$1 WHERE user_id=$2",
		about,
		id,
	)
	return err
}

//EditFirstName editing first name
func (r *UserRepository) EditFirstName(id int, firstName string) error {
	_, err := r.store.db.DB().Exec(
		"UPDATE profiles SET first_name=$1 WHERE user_id=$2",
		firstName,
		id,
	)
	return err
}

//EditLastName editing last name
func (r *UserRepository) EditLastName(id int, lastName string) error {
	_, err := r.store.db.DB().Exec(
		"UPDATE profiles SET last_name=$1 WHERE user_id=$2",
		lastName,
		id,
	)
	return err
}

//?

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

// CompareToken comparing given token with user token
func (r *UserRepository) CompareToken(u *model.User, token string) bool {
	usrToken, _ := r.store.User().GetToken(u.ID)
	if token == usrToken {
		return true
	}
	return false
}

//NewMessage storing new message in db
func (r *UserRepository) NewMessage(from int, to int, body string, timestamp int64) error {
	msg := &model.Message{
		FromID: from,
		ToID:   to,
		Body:   body,
		Time:   timestamp,
	}
	if err := r.store.db.Model(&model.Message{}).Create(&msg).Error; err != nil {
		return err
	}
	return nil
}

//GetMessageHistory giving messagehistory
func (r *UserRepository) GetMessageHistory(p1 int, p2 int) ([]*model.Message, error) {
	messages := make([]*model.Message, 0)
	rows1, err := r.store.db.DB().Query("SELECT id, from_id, to_id, body, time FROM messages WHERE from_id=$1 AND to_id=$2",
		p1,
		p2)
	if err != nil {
		return nil, err
	}
	defer rows1.Close()
	for rows1.Next() {
		var msg = &model.Message{}
		if err := rows1.Scan(&msg.ID, &msg.FromID, &msg.ToID, &msg.Body, &msg.Time); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if p1 != p2 {
		rows2, err := r.store.db.DB().Query("SELECT id, from_id, to_id, body, time FROM messages WHERE from_id=$1 AND to_id=$2",
			p2,
			p1)
		if err != nil {
			return nil, err
		}
		defer rows2.Close()
		for rows2.Next() {
			var msg = &model.Message{}
			if err := rows2.Scan(&msg.ID, &msg.FromID, &msg.ToID, &msg.Body, &msg.Time); err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		}
	}
	if len(messages) > 0 {
		return messages, nil
	}
	return nil, store.ErrNotMessages
}
