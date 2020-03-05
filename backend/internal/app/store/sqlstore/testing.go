package sqlstore

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	// GORM driver for Postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// TestDB ...
func TestDB(t *testing.T, databaseURL string) (*gorm.DB, func(...string)) {
	t.Helper()

	db, err := gorm.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.DB().Ping(); err != nil {
		t.Fatal(err)
	}
	return db, func(tables ...string) {
		if len(tables) > 0 {
			db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
		}

		db.Close()
	}

}
