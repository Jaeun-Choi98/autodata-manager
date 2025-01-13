package dao

import (
	entity "cju/entity/sample"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PGDB struct {
	db *gorm.DB
}

func NewPGDB(con string) (*PGDB, error) {
	db, err := gorm.Open(postgres.Open(con), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println("failed new db")
		return nil, err
	}
	return &PGDB{db}, err
}

func (pg *PGDB) ClosePostgreSQL() error {
	db, err := pg.db.DB()
	if err != nil {
		log.Println("failed close db")
		return err
	}
	return db.Close()
}

func (pg *PGDB) AutoMigrateJob() error {
	err := pg.db.AutoMigrate(&entity.Job{})
	if err != nil {
		log.Println("failed automigrate entity.job")
		return err
	}
	return nil
}

func (pg *PGDB) AutoMigrateUesr() error {
	err := pg.db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Println("failed automigrate entity.user")
		return err
	}
	return nil
}

func (pg *PGDB) AddUser(user entity.User) error {
	result := pg.db.Create(&user)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	return nil
}

func (pg *PGDB) GetUser() (*entity.User, error) {
	var user entity.User
	result := pg.db.First(&user)
	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}
	return &user, nil
}
