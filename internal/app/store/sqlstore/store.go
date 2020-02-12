package sqlstore

import (
	"database/sql"

	"github.com/MrSedan/restapigoown/internal/app/store"
)

// Store is the struct of db
type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

// New is Creating new Store
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User returning UserRepository inteface
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}
	return s.userRepository
}
