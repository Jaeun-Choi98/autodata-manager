package client

type ClientInterface interface {
	MakeTable(filePath, tableName, extension string) (map[string]interface{}, error)
	NormalizeTable(filePath, extension string) (map[string]interface{}, error)
	DropTable(tableName string) (map[string]interface{}, error)
	ExportTable(tableName, extension string) error
	ReadAllRecord(tableName string) (map[string]interface{}, error)
	ReadAllTables(schemaName string) (interface{}, error)
	Listen() (map[string]interface{}, error)
	Unlisten() (map[string]interface{}, error)
	CronCommand(param, jobId string) (map[string]interface{}, error)
	BackupDB(dbName string) (map[string]interface{}, error)
	CronBackupDB(dbName string, query []string) (map[string]interface{}, error)
	ReadAllSchemas() (map[string]interface{}, error)
	DropSchema(schemaName string) (map[string]interface{}, error)
	MakeSchema(schemaName string) (map[string]interface{}, error)
	UpdateUser(filePath string) (map[string]interface{}, error)
	RegisterUser(filePath string) (map[string]interface{}, error)
	Logout() (map[string]interface{}, error)
	Login(email, pwd string) (map[string]interface{}, error)
	ReadUserInfo(email string) (map[string]interface{}, error)
	GetToken() string
	GetEmail() string
}
