package model

//Message is a struct of message
type Message struct {
	ID     int    `json:"id" gorm:"not null;primary"`
	FromID int    `json:"from"`
	ToID   int    `json:"to"`
	Time   int64  `json:"time"`
	Body   string `json:"body"`
	Token  string `json:"token,omitempty" sql:"-"`
}
