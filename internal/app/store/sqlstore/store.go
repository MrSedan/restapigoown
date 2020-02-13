package sqlstore

import (
	"github.com/MrSedan/restapigoown/internal/app/model"
	"github.com/MrSedan/restapigoown/internal/app/store"
	"github.com/jinzhu/gorm"
)

// Store is the struct of db
type Store struct {
	db             *gorm.DB
	userRepository *UserRepository
}

// New is Creating new Store
func New(db *gorm.DB) *Store {
	db.AutoMigrate(&model.User{})
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
