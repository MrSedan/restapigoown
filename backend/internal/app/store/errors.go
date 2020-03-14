package store

import "errors"

var (
	// ErrRecordNotFound ...
	ErrRecordNotFound = errors.New("Record not found")
	// ErrNotValidToken ...
	ErrNotValidToken = errors.New("this token is not valid")
	//ErrNotMessages ...
	ErrNotMessages = errors.New("No messages in db")
)
