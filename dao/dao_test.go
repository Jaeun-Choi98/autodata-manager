package dao_test

import (
	"cju/dao"
	"cju/entity/auth"
	"log"

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

func TestAddUser(t *testing.T) {
	godotenv.Load("../.env")
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	con := fmt.Sprintf(testConnectionString, dbHost, dbPwd, dbName, dbPort)
	pgdb, _ := dao.NewPostgreSQL(con)
	pgdb.Init()
	roles := []auth.Role{
		{
			RoleName:    "Admin",
			Description: "with full access",
		},
		{
			RoleName:    "Employee",
			Description: "with restricted access",
		},
	}
	var users []*auth.User
	user1 := &auth.User{
		Username: "cju",
		Email:    "cju@aaa.com",
		Password: "hashed_password",
		IsActive: true,
		Profile: auth.Profile{
			UserID:         1,
			FirstName:      "choi",
			LastName:       "ju",
			PhoneNumber:    "111-1234-1424",
			Address:        "Busan",
			ProfilePicture: "link",
		},
		Roles: roles,
	}

	user2 := &auth.User{
		Username: "cju2",
		Email:    "cju2@aaa.com",
		Password: "hashed_password2",
		IsActive: true,
		Roles:    roles,
	}
	users = append(users, user1)
	users = append(users, user2)
	err := pgdb.AddUser(users)
	assert.NoError(t, err, "failed to 'AddRecordd'")
}

func TestUpdateUser(t *testing.T) {
	godotenv.Load("../.env")
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	con := fmt.Sprintf(testConnectionString, dbHost, dbPwd, dbName, dbPort)
	pgdb, _ := dao.NewPostgreSQL(con)
	pgdb.Init()
	roles := []auth.Role{
		{
			RoleName:    "Admin",
			Description: "with full access",
		},
		{
			RoleName:    "Employee",
			Description: "with restricted access",
		},
	}
	var users []*auth.User
	user1 := &auth.User{
		Username: "cju",
		Email:    "cju@aaa.com",
		Password: "hashed_password",
		IsActive: true,
		Profile: auth.Profile{
			UserID:         1,
			FirstName:      "choi3",
			LastName:       "ju23",
			PhoneNumber:    "111-1234-1424556",
			Address:        "Busannn",
			ProfilePicture: "link123",
		},
		Roles: roles,
	}

	user2 := &auth.User{
		Username: "cju2245",
		Email:    "cju2@aaa.com",
		Password: "hashed_password2",
		IsActive: true,
		Roles:    roles,
	}
	users = append(users, user1)
	users = append(users, user2)
	err := pgdb.UpdateUser(users)
	assert.NoError(t, err, "failed to 'UpdateRecord'")
}

func TestReadUserByEmail(t *testing.T) {
	godotenv.Load("../.env")
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	con := fmt.Sprintf(testConnectionString, dbHost, dbPwd, dbName, dbPort)
	pgdb, _ := dao.NewPostgreSQL(con)
	pgdb.Init()
	user, err := pgdb.ReadUserByEmail("cju@aaa.com")
	assert.NoError(t, err, "failed to read user by email")
	log.Println(user)
}
