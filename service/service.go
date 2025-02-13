package service

import (
	"cju/dao"
	"cju/entity/auth"
	"fmt"
	"log"
)

type ServiceInterface interface {
	CloseService() error
	CreateTableFromCSV(filePath, tableName string) error
	DropTableByTableName(tableName string) error
	CreateTableFromExcel(filePath, tableName string) error
	CreateTableFromJSON(filePath, tableName string) error
	ExportTableToJSON(tableName string) (error, string)
	ExportTableToCSV(tableName string) (error, string)
	CreateNormalizeTableFromCSV(filePath string) (string, error)
	ReadAllRecordByTableName(tableName string) ([]map[string]interface{}, error)
	ReadAllTablesBySchemaNamd(schemaName string) ([]string, error)
	StartListenerManager() error
	StopListenerManager() error
	BackupDatabase(dbName string) error
	CronBackupDataBase(dbName string, query []string) error
	RemoveCronJob(id string)
	GetJobList() map[string]string
	CronStart()
	CronStop()
	ReadAllSchemas() ([]string, error)
	DeleteSchema(schemaName string) error
	CreateSchema(schemaName string) error
	AddUserFromCSV(filePath string) error
	UpdateUserFromCSV(filePath string) error
	Login(email, pwd string) (string, error)
	ReadUserByEmail(email string) (*auth.User, error)
}

type Service struct {
	mydb   dao.DaoInterface
	mylm   ListenerManagerInterface
	mycron *CronManager
}

func NewService(dbHost, dbPort, dbPwd, dbName string) (ServiceInterface, error) {
	dbInfo := fmt.Sprintf("user=postgres dbname=%s password=%s host=%s port=%s sslmode=disable",
		dbName, dbPwd, dbHost, dbPort)
	db, err := dao.NewPostgreSQL(dbInfo)
	db.Init()
	if err != nil {
		return nil, err
	}
	lmCon := fmt.Sprintf("postgres://postgres:%s@%s:%s/%s", dbPwd, dbHost, dbPort, dbName)
	lm, err := NewListenManager(lmCon)
	if err != nil {
		return &Service{mydb: db, mylm: nil, mycron: NewCronInstance()}, nil
	}
	return &Service{mydb: db, mylm: lm, mycron: NewCronInstance()}, nil
}

func (s *Service) CloseService() error {
	err := s.mydb.CloseDB()
	if err != nil {
		log.Println("failed to close service")
		return err
	}
	err = s.mylm.Close()
	if err != nil {
		return err
	}
	s.mycron.StopCron()
	return nil
}
