package service

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestCreateTableFromCSV(t *testing.T) {
	godotenv.Load("../.env")
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	service, err := NewService(fmt.Sprintf("user=postgres dbname=%s password=%s host=%s port=%s sslmode=disable",
		dbName, dbPwd, dbHost, dbPort))
	if err != nil {
		log.Println("NewService method err")
		return
	}
	err = service.CreateTableFromCSV("../data.csv", "testtable")
	if err != nil {
		log.Println("Service.CreateTableFromCSV err")
	}
}

func TestReadAllRecordByTableName(t *testing.T) {
	godotenv.Load("../.env")
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	service, err := NewService(fmt.Sprintf("user=postgres dbname=%s password=%s host=%s port=%s sslmode=disable",
		dbName, dbPwd, dbHost, dbPort))
	if err != nil {
		log.Println("NewService method err")
		return
	}
	ret, _ := service.ReadAllRecordByTableName("testtable")
	fmt.Printf("return val: %v", ret)
}
