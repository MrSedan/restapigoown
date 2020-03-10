package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

// User a struct of user info
type User struct {
	ID                int    `json:"id" gorm:"not null;primary key"`
	Password          string `sql:"-" json:"-"`
	UserName          string `json:"user_name" sql:"user_name" gorm:"not null;unique"`
	Email             string `gorm:"not null;unique" json:"email,omitempty"`
	EncryptedPassword string `json:"encrypted_password,omitempty"`
	JwtToken          string `sql:"jwt_token" json:"-"`
}

// Validate a validate of user
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(8, 64)),
	)
}

// Sanitize make passw empty
func (u *User) Sanitize() {
	u.Password = ""
}

// ComparePassword ...
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

// BeforeCreate a func for create hashed passw
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = enc
	}
	return nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
