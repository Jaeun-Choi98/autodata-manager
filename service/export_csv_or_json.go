package service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func (s *Service) ExportTableToJsonAndCSV(tableName, filePath string) error {

	if !s.mydb.ExistTable(tableName) {
		log.Printf("Table '%s' doesn't exist", tableName)
		return fmt.Errorf("table '%s' doesn't exist", tableName)
	}

	rows, err := s.mydb.ReadAllTableData(tableName)
	if err != nil {
		return err
	}

	var headers []string
	if len(rows) > 0 {
		for col := range rows[0] {
			headers = append(headers, col)
		}
	}

	var records [][]string
	for _, row := range rows {
		var record []string
		for _, col := range headers {
			val := row[col]
			if val == nil {
				record = append(record, "")
			} else {
				record = append(record, fmt.Sprintf("%v", val))
			}
		}
		records = append(records, record)
	}

	err = SaveToJson(filePath, &headers, &records)
	if err != nil {
		return err
	}
	log.Printf("'%s' file saved successfully", filePath)

	err = SaveToCSV(filePath, &headers, &records)
	if err != nil {
		return err
	}
	log.Printf("'%s' file saved successfully", filePath)

	return nil
}

func SaveToCSV(filePath string, headers *[]string, records *[][]string) error {

	file, err := os.Create(fmt.Sprintf("%s.csv", filePath))
	if err != nil {
		log.Printf("failed to open '%s' file.", filePath)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(*headers); err != nil {
		return err
	}

	for _, row := range *records {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func SaveToJson(filePath string, headers *[]string, records *[][]string) error {

	file, err := os.Create(fmt.Sprintf("%s.json", filePath))
	if err != nil {
		log.Printf("failed to open '%s' file.", filePath)
		return err
	}
	defer file.Close()

	// 데이터를 JSON 형식으로 변환
	var jsonData []map[string]string
	for _, row := range *records {
		record := make(map[string]string)
		for i, header := range *headers {
			record[header] = row[i]
		}
		jsonData = append(jsonData, record)
	}

	// JSON 파일로 저장
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(jsonData); err != nil {
		log.Printf("failed to encode json '%s' file", filePath)
		return err
	}

	return nil
}
