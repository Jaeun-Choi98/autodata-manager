package rest

import (
	"cju/rest/handler"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Serve(addr string) error {
	godotenv.Load()
	dbHost, dbName, dbPwd, dbPort := os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT")
	h := handler.NewHandler(fmt.Sprintf("user=postgres dbname=%s password=%s host=%s port=%s sslmode=disable",
		dbName, dbPwd, dbHost, dbPort))
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
	return r.Run(addr)
}
