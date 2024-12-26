package service

/*
클라이언트에서 요청하는 데이터 형식을 맞추는 작업을 service계층에서 해야할 지
controller/handler 계층에서 해야할 지 생각해봐야함.
*/
func (s *Service) ReadAllRecordByTableName(tableName string) ([]map[string]interface{}, error) {
	records, err := s.mydb.ReadAllTableData(tableName)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (s *Service) ReadAllTablesBySchemaNamd(schemaName string) ([]string, error) {
	tables, err := s.mydb.ReadAllTables(schemaName)
	if err != nil {
		return nil, err
	}
	return tables, nil
}
