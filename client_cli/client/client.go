package client

type ClientInterface interface {
	MakeTable(filePath, tableName, extension string) (map[string]interface{}, error)
	NormalizeTable(filePath, extension string) (map[string]interface{}, error)
	DropTable(tableName string) (map[string]interface{}, error)
	ExportTable(tableName, extension string) error
	ReadAllRecord(tableName string) (map[string]interface{}, error)
	ReadAllTables(schemaName string) (interface{}, error)
	SubscribeDDL() (map[string]interface{}, error)
	UnsubscribeDDL() (map[string]interface{}, error)
	CronCommand(param, jobId string) (map[string]interface{}, error)
	BackupDB(dbName string) (map[string]interface{}, error)
	CronBackupDB(dbName string, query []string) (map[string]interface{}, error)
}
