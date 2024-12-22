package service

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func (s *Service) CreateTableFromCSV(filePath, tableName string) error {

	if exists := s.mydb.ExistTable(tableName); exists {
		log.Printf("existed '%s' table", tableName)
		return fmt.Errorf("existed '%s' table", tableName)
	}

	//read csv file
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("failed to open CSV file: %v", err)
		return err
	}
	defer file.Close()

	//read headers and records
	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		log.Printf("failed to read CSV header: %v", err)
		return err
	}
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("failed to read records")
		return err
	}

	err = CreateTableFromStringArr(s, &tableName, &headers, &records)
	if err != nil {
		log.Printf("failed to create table: %v", err)
		return err
	}
	log.Printf("Table '%s' created successfully!", tableName)

	// add records
	if len(records) == 0 {
		log.Printf("nothing records")
	} else {
		err = AddStringArrRecord(s, &tableName, &headers, &records)
		if err != nil {
			log.Printf("failed to add records")
			return err
		}
		log.Printf("Add records!")
	}

	return nil
}
