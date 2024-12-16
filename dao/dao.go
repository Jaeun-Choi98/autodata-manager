package dao

type DaoInterface interface {
	Init() error
	CloseDB() error
	ExecQuery(query string) error
	ExistTable(tableName string) bool
}
