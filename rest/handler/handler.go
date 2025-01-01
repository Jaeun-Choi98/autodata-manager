package handler

import (
	"cju/service"

	"github.com/gin-gonic/gin"
)

type HandlerInterface interface {
	CreateTableCSV(c *gin.Context)
	CreateTableJSON(c *gin.Context)
	CreateTableExcel(c *gin.Context)
	CreateNormalizeTableCSV(c *gin.Context)
	DeleteTable(c *gin.Context)
	ExportTableJSON(c *gin.Context)
	ExportTableCSV(c *gin.Context)
	ReadAllRecordByTableName(c *gin.Context)
	ReadAllTablesBySchema(c *gin.Context)
	CloseHandler()
	Listen(c *gin.Context)
	Unlisten(c *gin.Context)

	BackupDB(c *gin.Context)
	CronBackupDB(c *gin.Context)
	RemoveCronJob(c *gin.Context)
	GetJobList(c *gin.Context)
	CronStart(c *gin.Context)
	CronStop(c *gin.Context)
}

type Handler struct {
	myService service.ServiceInterface
}

func NewHandler(dbHost, dbPort, dbPwd, dbName string) HandlerInterface {
	sv, _ := service.NewService(dbHost, dbPort, dbPwd, dbName)
	return &Handler{myService: sv}
}

func (h *Handler) CloseHandler() {
	h.myService.CloseService()
}
