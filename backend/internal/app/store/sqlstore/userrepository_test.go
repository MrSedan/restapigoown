package sqlstore_test

import (
	"testing"

	"github.com/MrSedan/restapigoown/backend/internal/app/model"
	"github.com/MrSedan/restapigoown/backend/internal/app/store"
	"github.com/MrSedan/restapigoown/backend/internal/app/store/sqlstore"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")
	db.AutoMigrate(&model.User{})

	s := sqlstore.New(db)
	u := model.TestUser(t)

	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")
	db.AutoMigrate(&model.User{})
	s := sqlstore.New(db)
	email := "user@example.org"
	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = email
	s.User().Create(u)
	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_GetTokenAndClaimToken(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")
	db.AutoMigrate(&model.User{})
	s := sqlstore.New(db)
	_, err := s.User().CheckToken("test")
	assert.EqualError(t, err, store.ErrNotValidToken.Error())

	u := model.TestUser(t)
	s.User().Create(u)
	s.User().ClaimToken(u, "test")
	_, err = s.User().CheckToken("test")
	assert.NoError(t, err)
}
