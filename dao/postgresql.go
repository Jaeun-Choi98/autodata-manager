package dao

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQL struct {
	db *gorm.DB
}

func NewPostgreSQL(con string) (*PostgreSQL, error) {
	db, err := gorm.Open(postgres.Open(con), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println("failed new db")
		return nil, err
	}
	return &PostgreSQL{db}, err
}

func (pg *PostgreSQL) CloseDB() error {
	db, err := pg.db.DB()
	if err != nil {
		log.Println("failed to close db")
		return err
	}
	return db.Close()
}

func (pq *PostgreSQL) ExecQuery(query string) error {
	err := pq.db.Exec(query).Error
	if err != nil {
		log.Println("failed to Exec Query")
		return err
	}
	return nil
}

func (pq *PostgreSQL) Init() error {
	return nil
}

func (pq *PostgreSQL) ExistTable(tableName string) bool {
	exists := pq.db.Migrator().HasTable(tableName)
	return exists
}
