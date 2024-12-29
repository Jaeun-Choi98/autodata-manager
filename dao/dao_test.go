package dao_test

import (
	"cju/dao"
	"cju/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testConnectionString = "host=localhost user=postgres password=cjswo123 dbname=test port=5432 sslmode=disable"

func TestNewPGDB(t *testing.T) {
	pgdb, err := dao.NewPGDB(testConnectionString)
	assert.NoError(t, err)
	assert.NotNil(t, pgdb)

	err = pgdb.ClosePostgreSQL()
	assert.NoError(t, err)
}

func TestAutoMigrateJob(t *testing.T) {
	pgdb, err := dao.NewPGDB(testConnectionString)
	assert.NoError(t, err)
	defer pgdb.ClosePostgreSQL()

	err = pgdb.AutoMigrateJob()
	assert.NoError(t, err)
}

func TestAutoMigrateUser(t *testing.T) {
	pgdb, err := dao.NewPGDB(testConnectionString)
	assert.NoError(t, err)
	defer pgdb.ClosePostgreSQL()

	err = pgdb.AutoMigrateUesr()
	assert.NoError(t, err)
}

func TestAutoMigrateJobAndUser(t *testing.T) {
	pgdb, err := dao.NewPGDB(testConnectionString)
	assert.NoError(t, err)
	defer pgdb.ClosePostgreSQL()
	err = pgdb.AutoMigrateUesr()
	pgdb.AutoMigrateJob()

	assert.NoError(t, err)
}

func TestAddUser(t *testing.T) {
	pgdb, err := dao.NewPGDB(testConnectionString)
	assert.NoError(t, err)
	defer pgdb.ClosePostgreSQL()

	err = pgdb.AutoMigrateUesr()
	assert.NoError(t, err)

	user := entity.User{
		Age:     30,
		Name:    "Test User",
		Hobbies: []string{"reading", "coding"},
		Jobs: []entity.Job{
			{Name: "Software Engineer"},
		},
	}

	err = pgdb.AddUser(user)
	assert.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	pgdb, err := dao.NewPGDB(testConnectionString)
	assert.NoError(t, err)
	defer pgdb.ClosePostgreSQL()

	err = pgdb.AutoMigrateUesr()
	assert.NoError(t, err)

	user := entity.User{
		Age:     30,
		Name:    "Test User",
		Hobbies: []string{"reading", "coding"},
		Jobs: []entity.Job{
			{Name: "Software Engineer"},
		},
	}

	err = pgdb.AddUser(user)
	assert.NoError(t, err)

	retrievedUser, err := pgdb.GetUser()
	assert.NoError(t, err)
	assert.NotNil(t, retrievedUser)
	assert.Equal(t, user.Name, retrievedUser.Name)
	assert.Equal(t, user.Age, retrievedUser.Age)
	assert.ElementsMatch(t, user.Hobbies, retrievedUser.Hobbies)
}
