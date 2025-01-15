package rest

import (
	"cju/rest/handler"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Serve(addr string) error {
	godotenv.Load()
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	h := handler.NewHandler(dbHost, dbPort, dbPwd, dbName)
	return serveWithHandler(addr, h)
}

func serveWithHandler(addr string, h handler.HandlerInterface) error {
	defer h.CloseHandler()
	r := gin.Default()
	r.POST("/service/create/csv", h.CreateTableCSV)
	r.POST("/service/create/json", h.CreateTableJSON)
	r.POST("/service/create/excel", h.CreateTableExcel)
	r.POST("/service/create/normalize/csv", h.CreateNormalizeTableCSV)
	r.POST("/service/delete", h.DeleteTable)
	r.POST("/service/export/json", h.ExportTableJSON)
	r.POST("/service/export/csv", h.ExportTableCSV)
	r.POST("/service/read-table-all", h.ReadAllRecordByTableName)
	r.POST("/service/get-all-tables", h.ReadAllTablesBySchema)
	r.GET("/service/listen", h.Listen)
	r.GET("/service/unlisten", h.Unlisten)
	r.POST("/service/backup/database", h.BackupDB)
	r.POST("/service/backup/cron", h.CronBackupDB)
	r.POST("/service/cron/remove", h.RemoveCronJob)
	r.GET("/service/cron/jobs", h.GetJobList)
	r.GET("/service/cron/start", h.CronStart)
	r.GET("/service/cron/stop", h.CronStop)
	r.GET("/service/schema/list", h.ReadAllSchemas)
	r.POST("/service/schema/create", h.CreateSchema)
	r.POST("/service/schema/delete", h.DeleteSchema)
	r.POST("/service/user/register", h.RegisterUser)
	r.POST("/service/user/update", h.UpdateUser)
	return r.Run(addr)
}
