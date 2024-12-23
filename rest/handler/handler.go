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
	CloseHandler()
}

type Handler struct {
	myService service.ServiceInterface
}

func NewHandler(dbInfo string) HandlerInterface {
	sv, _ := service.NewService(dbInfo)
	return &Handler{myService: sv}
}

func (h *Handler) CloseHandler() {
	h.myService.CloseService()
}
