package model

// Profile of user
type Profile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserID    int    `json:"user_id" gorm:"association_foreignkey:id"`
	UserEmail string `json:"user_email"`
	About     string `json:"about,omitempty"`
}
