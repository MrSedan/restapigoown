package apiserver

import (
	"net/http"

	"github.com/MrSedan/restapigoown/backend/internal/app/store/sqlstore"
	"github.com/jinzhu/gorm"

	// This is driver for PostgresDB
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Start is a function for start server
func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	store := sqlstore.New(db)
	srv := newServer(store)
	srv.jwtKey = config.JwtKey
	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
