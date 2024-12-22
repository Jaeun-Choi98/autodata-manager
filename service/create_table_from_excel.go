package service

import (
	"fmt"
	"log"

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
	err = CreateTableFromStringArr(s, &tableName, &rows[0], &records, false)
	if err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}
	log.Printf("Table '%s' created successfully!", tableName)

	// 레코드 추가
	if len(records) == 0 {
		log.Printf("nothing records")
	} else {
		err = AddStringArrRecord(s, &tableName, &rows[0], &records)
		if err != nil {
			log.Printf("failed to add records: %v", err)
			return err
		}
		log.Println("Records added successfully!")
	}
	return nil
}
