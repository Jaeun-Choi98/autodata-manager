package dao

import (
	"fmt"
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

func (pq *PostgreSQL) ReadAllTableData(tableName string) ([]map[string]interface{}, error) {
	var rows []map[string]interface{}
	err := pq.db.Table(tableName).Find(&rows).Error
	if err != nil {
		log.Printf("failed to query table data: %v", err)
		return nil, err
	}
	return rows, nil
}

func (pq *PostgreSQL) ReadAllTables(schemaName string) ([]string, error) {
	var tables []string
	query := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'", schemaName)
	err := pq.db.Raw(query).Scan(&tables).Error
	if err != nil {
		log.Printf("failed to 'ReadAllTables' agr(%s)", schemaName)
		return nil, err
	}
	return tables, err
}

func (pq *PostgreSQL) ExistSchema(schemaName string) (bool, error) {
	var schemaNameResult string
	err := pq.db.Raw("SELECT schema_name FROM information_schema.schemata WHERE schema_name = ?", schemaName).Scan(&schemaNameResult).Error
	if err != nil {
		log.Printf("failed to 'ExistSchema' arg(%s, %v)", schemaName, err)
		return true, err
	} else if schemaNameResult == "" {
		return false, nil
	} else {
		return true, nil
	}
}

func (pq *PostgreSQL) ReadAllSchemas() ([]string, error) {
	var schemas []string
	err := pq.db.Raw("SELECT schema_name FROM information_schema.schemata").Scan(&schemas).Error
	if err != nil {
		log.Printf("failed to 'ReadAllSchemas' (%v)", err)
		return nil, err
	}
	return schemas, nil
}
