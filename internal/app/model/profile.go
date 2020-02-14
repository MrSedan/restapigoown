package model

// Profile of user
type Profile struct {
	User      *User  `gorm:"association_foreignkey:Email"`
	UserEmail string `json:"user_email"`
	About     string `json:"about,omitempty"`
}
