package sqlstore_test

import (
	"testing"

	"github.com/MrSedan/restapigoown/internal/app/model"
	"github.com/MrSedan/restapigoown/internal/app/store/sqlstore"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)

	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}
