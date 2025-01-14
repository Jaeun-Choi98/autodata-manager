package auth

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func AutoMigrateAuthSchema() error {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("failed to load .env")
		return err
	}
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	dbInfo := fmt.Sprintf("user=postgres dbname=%s password=%s host=%s port=%s sslmode=disable",
		dbName, dbPwd, dbHost, dbPort)
	db, err := gorm.Open(postgres.Open(dbInfo), &gorm.Config{})
	if err != nil {
		log.Println("failed to connect db")
		return err
	}
	err = db.AutoMigrate(
		&User{},
		&Profile{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database")
		return err
	}
	return nil
}
