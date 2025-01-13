package dao_test

import (
	"cju/dao"

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
