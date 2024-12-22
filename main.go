package main

import "cju/service"

func main() {

	// var mydb dao.DBLayerInterface
	con := "user=postgres dbname=test password=cjswo123 host=localhost port=5432 sslmode=disable"
	// mydb, _ = dao.NewPostgreSQL(con)
	// defer mydb.ClosePostgreSQL()

	// mydb.AutoMigrateUesr()
	// err := mydb.AutoMigrateJob()
	// if err != nil {
	// 	log.Println(err)
	// }

	var sv service.ServiceInterface
	sv, _ = service.NewService(con)
	defer sv.CloseService()
	//sv.CreateTableFromCSV("data_nothing_records.csv", "datatest")
	//sv.DropTableByTableName("data_excel")
	//sv.CreateTableFromCSV("data.csv", "data_csv")
	//sv.CreateTableFromExcel("data.xlsx", "data_excel")
	//sv.CreateTableFromJSON("data.json", "data_json")
	//sv.ExportTableToJsonAndCSV("data", "export_data")
	sv.CreateNormalizeTableFromCSV("data.csv")
}
