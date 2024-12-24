package client

type ClientInterface interface {
	MakeTable(url, fileName, tableName, extension string) error
	NormalizeTable(url, filePath, extension string) error
	DropTable(url, tableName string) error
	ExportTable(url, tableName, extension string) error
}
