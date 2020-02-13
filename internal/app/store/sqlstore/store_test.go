package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=postgres dbname=myrest_test sslmode=disable user=user password=test"
	}
	os.Exit(m.Run())
}
