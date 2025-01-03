package service

import (
	"fmt"
	"log"
)

func (s *Service) ReadAllSchemas() ([]string, error) {
	schemas, err := s.mydb.ReadAllSchemas()
	if err != nil {
		return nil, err
	}
	return schemas, nil
}

func (s *Service) CreateSchema(schemaName string) error {
	exist, err := s.mydb.ExistSchema(schemaName)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("schema '%s' exists", schemaName)
	}
	err = s.mydb.ExecQuery(fmt.Sprintf("CREATE SCHEMA %s", schemaName))
	if err != nil {
		return nil
	}
	log.Printf("schema '%s' is created successfully!", schemaName)
	return nil
}

func (s *Service) DeleteSchema(schemaName string) error {
	exist, err := s.mydb.ExistSchema(schemaName)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("schema '%s' does not exist", schemaName)
	}
	err = s.mydb.ExecQuery(fmt.Sprintf("DROP SCHEMA %s", schemaName))
	if err != nil {
		return nil
	}
	log.Printf("schema '%s' is deleted successfully!", schemaName)
	return nil
}
