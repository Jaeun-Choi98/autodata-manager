package service

import (
	"fmt"
	"log"
)

func (s *Service) DropTableByTableName(tableName string) error {
	err := s.mydb.ExecQuery(fmt.Sprintf("DROP TABLE %s", tableName))
	if err != nil {
		log.Println("failed to drop table by table_name")
		return err
	}

	log.Printf("Table '%s' deleted successfully!", tableName)

	return nil
}
