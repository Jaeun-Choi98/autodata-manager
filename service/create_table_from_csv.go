package service

import (
	"cju/service/grpc"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func (s *Service) CreateNormalizeTableFromCSV(filePath string) error {

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("failed to open CSV file: %v", err)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("failed to read records")
		return err
	}

	var reqData strings.Builder
	for _, record := range records {
		reqData.WriteString(strings.Join(record, ","))
		reqData.WriteString("\n")
	}

	/*
		비정규화된 테이블 -> 정규화된 테이블
		1. ai를 사용해서 정규화된 스키마를 얻음. + 알고리즘을 사용해서 데이터를 옮김( 외래키에 대한 처리 )
		2. ai를 사용하여 정규화된 스키마와 데이터를 얻음.
		어떤 것이 좋을지 생각해 봐야함.
	*/
	resp := grpc.NormalizeByOpenAI(reqData.String())
	if resp == "" {
		log.Printf("failed to normalize '%s' file by using grpc", filePath)
		return fmt.Errorf("failed to normalize '%s' file by using grpc", filePath)
	}
	// 이후 실패했을 때, 에러 처리 필요.
	normalizedTables := ParseNormalizationData(resp)

	// 도중에 에러가 생겼을 때, 계속해서 마이그레이션 할 것인지 아니면 에러처리를 해서 멈출 것인지 정할 필요가 있음.
	for tableName, rows := range *normalizedTables {
		tableName = strings.ToLower(tableName)
		if exists := s.mydb.ExistTable(tableName); exists {
			log.Printf("existed '%s' table", tableName)
			return fmt.Errorf("existed '%s' table", tableName)
		}
		headers := rows[0]
		normalizedRecords := rows[1:]
		err = CreateTableFromStringArr(s, &tableName, &headers, &normalizedRecords, true)
		if err != nil {
			log.Printf("failed to create table: %v", tableName)
			return err
		}
		log.Printf("Table '%s' created successfully!", tableName)
		err = AddStringArrRecord(s, &tableName, &headers, &normalizedRecords)
		if err != nil {
			log.Printf("failed to add records: '%v' table", tableName)
			return err
		}
		log.Printf("Add '%v' table records!", tableName)
	}
	return nil
}

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

	err = CreateTableFromStringArr(s, &tableName, &headers, &records, false)
	if err != nil {
		log.Printf("failed to create table: %v", tableName)
		return err
	}
	log.Printf("Table '%s' created successfully!", tableName)

	// add records
	if len(records) == 0 {
		log.Printf("nothing records")
	} else {
		err = AddStringArrRecord(s, &tableName, &headers, &records)
		if err != nil {
			log.Printf("failed to add records '%v' table", tableName)
			return err
		}
		log.Printf("Add records!")
	}

	return nil
}
