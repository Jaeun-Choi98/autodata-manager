package client

type ClientInterface interface {
	MakeTable(url, filePath, tableName, extension string) (map[string]interface{}, error)
	NormalizeTable(url, filePath, extension string) (map[string]interface{}, error)
	DropTable(url, tableName string) (map[string]interface{}, error)
	ExportTable(url, tableName, extension string) error
}
