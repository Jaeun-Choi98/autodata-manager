package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// 이후 normalized_data가 복잡해지면 필요.
func ParseNormalizationData(data string) *map[string][][]string {
	var result map[string][][]string
	_ = json.Unmarshal([]byte(data), &result)
	return &result
}

func CreateTableFromStringArr(s *Service, tableName *string, headers *[]string, records *[][]string, hasPrimaryKey bool) error {

	fields := make([]struct {
		Name     string
		DataType string
	}, 0)

	if !hasPrimaryKey {
		fields = append(fields, struct {
			Name     string
			DataType string
		}{"id", "SERIAL PRIMARY KEY"})
	}

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

		// default type: TEXT
		fields = append(fields, struct {
			Name     string
			DataType string
		}{header, dataType})
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

func AddStringArrRecord(s *Service, tableName *string, headers *[]string, records *[][]string) error {

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

// utils, 공통적으로 사용되는 함수. 이후 파일 구조 변경 필요.
func inferDataType(columnData *[]string) string {

	if len(*columnData) == 0 {
		return "TEXT"
	}

	isInt := true
	isFloat := true
	isBool := true
	for _, value := range *columnData {
		if _, err := strconv.Atoi(value); err != nil {
			isInt = false
		}
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			isFloat = false
		}
		if value != "true" && value != "false" {
			isBool = false
		}
	}

	switch {
	case isInt:
		return "INTEGER"
	case isFloat:
		return "FLOAT"
	case isBool:
		return "BOOLEAN"
	default:
		return "TEXT"
	}
}
