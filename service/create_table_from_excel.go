package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

func (s *Service) CreateTableFromExcel(filePath, tableName string) error {

	if s.mydb.ExistTable(tableName) {
		log.Printf("existed '%s' table", tableName)
		return fmt.Errorf("existed '%s' table", tableName)
	}

	file, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Printf("failed to open CSV file: %v", err)
		return err
	}
	defer file.Close()

	// 첫 번째 시트 가져오기
	sheetName := file.GetSheetName(0)
	if sheetName == "" {
		log.Println("No sheets found in the Excel file")
		return fmt.Errorf("no sheets found")
	}

	// Excel 데이터를 읽기
	rows, err := file.GetRows(sheetName)
	if err != nil {
		log.Printf("failed to read rows: %v", err)
		return err
	}

	var records [][]string
	if len(rows) > 0 {
		records = rows[1:] // 헤더를 제외한 나머지 행
	} else {
		records = make([][]string, 0)
	}

	// 테이블 생성
	err = createTableFromExcelHeaders(s, &tableName, &rows[0], &records)
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}
	log.Printf("Table '%s' created successfully!", tableName)

	// 레코드 추가
	if len(records) == 0 {
		log.Printf("nothing records")
	} else {
		err = addExcelRecord(s, &tableName, &rows[0], &records)
		if err != nil {
			log.Printf("failed to add records: %v", err)
			return err
		}
		log.Println("Records added successfully!")
	}
	return nil
}

func createTableFromExcelHeaders(s *Service, tableName *string, headers *[]string, records *[][]string) error {

	fields := make([]struct {
		Name     string
		DataType string
	}, 0)

	// add 'ID' field
	fields = append(fields, struct {
		Name     string
		DataType string
	}{"id", "SERIAL PRIMARY KEY"})

	// 샘플링 size 100, 이후 랜덤하게(or규칙적이게) 샘플링 하는 로직 필요할 수도 있음.
	size := len(*records)
	if size > 100 {
		size = 100
	}

	for col, header := range *headers {

		series := make([]string, size)
		for i := 0; i < size; i++ {
			series[i] = (*records)[i][col]
		}

		dataType := inferDataType(&series)

		fields = append(fields, struct {
			Name     string
			DataType string
		}{header, dataType})
	}

	// SQL 쿼리 빌드
	var query strings.Builder
	query.WriteString(fmt.Sprintf("CREATE TABLE %s (", *tableName))
	for i, field := range fields {
		query.WriteString(fmt.Sprintf("%s %s", field.Name, field.DataType))
		if i < len(fields)-1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(");")

	return s.mydb.ExecQuery(query.String())
}

func addExcelRecord(s *Service, tableName *string, headers *[]string, records *[][]string) error {

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
