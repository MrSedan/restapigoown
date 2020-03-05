package store

// Store is a main store
type Store interface {
	User() UserRepository
}
