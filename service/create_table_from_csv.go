package service

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func (s *Service) CreateTableFromCSV(filePath, tableName string) error {

	exists := s.mydb.ExistTable(tableName)
	if exists {
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

	//read headers
	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		log.Printf("failed to read CSV header: %v", err)
		return err
	}

	// createTableFromCSVHeaders 함수: CSV 헤더 기반으로 테이블 생성
	err = createTableFromCSVHeaders(s, &tableName, &headers)
	if err != nil {
		log.Printf("failed to create table: %v", err)
		return err
	}
	log.Printf("Table '%s' created successfully!", tableName)

	//레코드 추가
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("nothing records")
		return err
	}
	err = addCSVRecord(s, &tableName, &headers, &records)
	if err != nil {
		return err
	}
	log.Printf("Add records!")

	return nil
}

func createTableFromCSVHeaders(s *Service, tableName *string, headers *[]string) error {

	fields := make([]struct {
		Name     string
		DataType string
	}, 0)

	// 도메인 설정 로직 필요.
	fields = append(fields, struct {
		Name     string
		DataType string
	}{"id", "SERIAL PRIMARY KEY"})

	for _, header := range *headers {
		// default type: TEXT
		fields = append(fields, struct {
			Name     string
			DataType string
		}{header, "TEXT"})
	}

	// build sql query
	var query strings.Builder
	leng := len(fields)
	query.WriteString(fmt.Sprintf("CREATE TABLE %s (", *tableName))
	for i, field := range fields {
		query.WriteString(fmt.Sprintf("%s %s", field.Name, field.DataType))
		if i < leng-1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(");")

	return s.mydb.ExecQuery(query.String())
}

func addCSVRecord(s *Service, tableName *string, headers *[]string, records *[][]string) error {

	var query strings.Builder
	query.WriteString(fmt.Sprintf("INSERT INTO %s(%s) VALUES ", *tableName, strings.Join(*headers, ", ")))

	leng := len(*records)
	for i, record := range *records {
		for j, field := range record {
			record[j] = fmt.Sprintf("'%s'", field)
		}
		query.WriteString(fmt.Sprintf("(%s)", strings.Join(record, ", ")))
		if i < leng-1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(";")

	err := s.mydb.ExecQuery(query.String())
	if err != nil {
		log.Printf("failed to add records")
		return err
	}
	return nil
}
