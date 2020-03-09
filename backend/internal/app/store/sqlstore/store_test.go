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
		//databaseURL = "host=postgres dbname=myrest_test sslmode=disable user=tester password=test"
		databaseURL = "host=localhost dbname=myrest_test sslmode=disable user=postgres password=postgres"
	}
	os.Exit(m.Run())
}
