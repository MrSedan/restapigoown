package apiserver

import (
	"database/sql"
	"net/http"

	"github.com/MrSedan/restapigoown/internal/app/store/sqlstore"
	// This is driver for PostgresDB
	_ "github.com/lib/pq"
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

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
