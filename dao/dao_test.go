package dao_test

import (
	"cju/dao"
	"cju/entity"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

const testConnectionString = "host=%s user=postgres password=%s dbname=%s port=%s sslmode=disable"

func TestReadAllSchemas(t *testing.T) {
	godotenv.Load("../.env")
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	con := fmt.Sprintf(testConnectionString, dbHost, dbPwd, dbName, dbPort)
	pgdb, _ := dao.NewPostgreSQL(con)
	ret, _ := pgdb.ReadAllSchemas()
	fmt.Println(ret)
}

func TestExistSchema(t *testing.T) {
	godotenv.Load("../.env")
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	con := fmt.Sprintf(testConnectionString, dbHost, dbPwd, dbName, dbPort)
	pgdb, _ := dao.NewPostgreSQL(con)
	ret, _ := pgdb.ExistSchema("public")
	fmt.Println(ret)
	assert.Equal(t, true, ret)
	ret, _ = pgdb.ExistSchema("pbblic")
	fmt.Println(ret)
	assert.Equal(t, false, ret)
}

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
