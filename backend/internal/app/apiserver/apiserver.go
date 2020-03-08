package apiserver

import (
	"context"
	"net"
	"net/http"

	"github.com/MrSedan/restapigoown/backend/internal/app/store/sqlstore"
	"github.com/jinzhu/gorm"

	// This is driver for PostgresDB
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type contextKey struct {
	key string
}

var connContextKey = &contextKey{"http-conn"}

func saveConnInContext(ctx context.Context, c net.Conn) context.Context {
	return context.WithValue(ctx, connContextKey, c)
}

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
	server := http.Server{
		Addr:        config.BindAddr,
		ConnContext: saveConnInContext,
		Handler:     srv,
	}
	return server.ListenAndServe()
}

func newDB(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
