package dao

type DaoInterface interface {
	Init() error
	CloseDB() error
	ExecQuery(query string) error
	ExistTable(tableName string) bool
	ReadAllTableData(tableName string) ([]map[string]interface{}, error)
}
