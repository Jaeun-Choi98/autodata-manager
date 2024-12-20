package service

import (
	"cju/dao"
	"log"
	"strconv"
)

type ServiceInterface interface {
	CloseService() error
	CreateTableFromCSV(filePath, tableName string) error
	DropTableByTableName(tableName string) error
	CreateTableFromExcel(filePath, tableName string) error
	CreateTableFromJSON(filePath, tableName string) error
	ExportTableToJsonAndCSV(tableName, filePath string) error
}

type Service struct {
	mydb dao.DaoInterface
}

func NewService(dbInfo string) (ServiceInterface, error) {
	db, err := dao.NewPostgreSQL(dbInfo)
	if err != nil {
		return nil, err
	}
	return &Service{mydb: db}, nil
}

func (s *Service) CloseService() error {
	err := s.mydb.CloseDB()
	if err != nil {
		log.Println("failed to close service")
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
