package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Record map[string]interface{} // 동적으로 키-값을 가지는 레코드

func (s *Service) CreateTableFromJSON(filePath, tableName string) error {
	exists := s.mydb.ExistTable(tableName)
	if exists {
		log.Printf("Table '%s' already exists", tableName)
		return nil
	}

	// JSON 파일 열기
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("failed to open JSON file: %v", err)
		return err
	}
	defer file.Close()

	// JSON 데이터 읽기
	var records []Record
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&records)
	if err != nil {
		log.Printf("failed to decode JSON: %v", err)
		return err
	}

	// headers 동적으로 추출 (첫 번째 레코드의 키를 사용)
	headers := extractHeaders(&records)

	// 테이블 생성
	err = createTableFromJSONHeaders(s, &tableName, &headers)
	if err != nil {
		log.Printf("failed to create table: %v", err)
		return err
	}
	log.Printf("Table '%s' created successfully!", tableName)

	// 레코드 추가
	err = addJSONRecords(s, &tableName, &headers, &records)
	if err != nil {
		return err
	}
	log.Printf("Records added successfully!")

	return nil
}

// headers를 첫 번째 레코드에서 동적으로 추출하는 함수
func extractHeaders(records *[]Record) []string {
	var headers []string
	if len(*records) > 0 {
		// 첫 번째 레코드의 키를 headers로 사용
		for key := range (*records)[0] {
			headers = append(headers, key)
		}
	}
	return headers
}

func addJSONRecords(s *Service, tableName *string, headers *[]string, records *[]Record) error {
	var query strings.Builder
	query.WriteString(fmt.Sprintf("INSERT INTO %s(%s) VALUES ", *tableName, strings.Join(*headers, ", ")))

	leng := len(*records)
	for i, record := range *records {
		recordValues := []string{}
		for _, header := range *headers {
			// 각 레코드에서 header에 해당하는 값을 가져와서 작은따옴표로 감싸기
			val := fmt.Sprintf("'%v'", record[header])
			recordValues = append(recordValues, val)
		}
		query.WriteString(fmt.Sprintf("(%s)", strings.Join(recordValues, ", ")))
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

func createTableFromJSONHeaders(s *Service, tableName *string, headers *[]string) error {
	fields := make([]struct {
		Name     string
		DataType string
	}, 0)

	// 기본적으로 id 필드 추가
	fields = append(fields, struct {
		Name     string
		DataType string
	}{"id", "SERIAL PRIMARY KEY"})

	// JSON 헤더를 기반으로 필드 정의
	for _, header := range *headers {
		// 기본 데이터 타입은 TEXT로 설정
		fields = append(fields, struct {
			Name     string
			DataType string
		}{header, "TEXT"})
	}

	// 테이블 생성 SQL 쿼리 빌드
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

	// SQL 쿼리 실행
	return s.mydb.ExecQuery(query.String())
}
