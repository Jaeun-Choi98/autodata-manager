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
	public := r.Group("/service")
	{
		public.POST("/user/login", h.Login)
	}

	employee := r.Group("/service")
	employee.Use(handler.JWTMiddlewareEmployee())
	{
		employee.POST("/create/csv", h.CreateTableCSV)
		employee.POST("/create/json", h.CreateTableJSON)
		employee.POST("/create/excel", h.CreateTableExcel)
		employee.POST("/create/normalize/csv", h.CreateNormalizeTableCSV)
		employee.POST("/delete", h.DeleteTable)
		employee.POST("/export/json", h.ExportTableJSON)
		employee.POST("/export/csv", h.ExportTableCSV)
		employee.POST("/read-table-all", h.ReadAllRecordByTableName)
		employee.POST("/get-all-tables", h.ReadAllTablesBySchema)
		employee.GET("/listen", h.Listen)
		employee.GET("/unlisten", h.Unlisten)
		employee.POST("/backup/database", h.BackupDB)
		employee.POST("/backup/cron", h.CronBackupDB)
		employee.POST("/cron/remove", h.RemoveCronJob)
		employee.GET("/cron/jobs", h.GetJobList)
		employee.GET("/cron/start", h.CronStart)
		employee.GET("/cron/stop", h.CronStop)
		employee.GET("/schema/list", h.ReadAllSchemas)
	}

	admin := r.Group("/service")
	admin.Use(handler.JWTMiddlewareAdmin())
	{
		admin.POST("/schema/create", h.CreateSchema)
		admin.POST("/schema/delete", h.DeleteSchema)
		admin.POST("/user/register", h.RegisterUser)
		admin.POST("/user/update", h.UpdateUser)
	}
	return r.Run(addr)
}
