package dao

import "cju/entity/auth"

type DaoInterface interface {
	Init() error
	CloseDB() error
	ExecQuery(query string) error
	ExistTable(tableName string) bool
	ReadAllTableData(tableName string) ([]map[string]interface{}, error)
	ReadAllTables(schemaName string) ([]string, error)
	ExistSchema(schemaName string) (bool, error)
	ReadAllSchemas() ([]string, error)
	AddUser(users []*auth.User) error
	UpdateUser(users []*auth.User) error
}
