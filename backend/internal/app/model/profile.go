package model

// Profile of user
type Profile struct {
	UserID    int    `json:"user_id" gorm:"association_foreignkey:id"`
	UserEmail string `json:"user_email"`
	About     string `json:"about,omitempty"`
}
