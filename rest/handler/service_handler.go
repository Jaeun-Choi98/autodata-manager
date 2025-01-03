package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ReadAllSchemas(c *gin.Context) {
	schemas, err := h.myService.ReadAllSchemas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"schema_list": schemas,
	})
}

func (h *Handler) DeleteSchema(c *gin.Context) {
	schemaName := c.PostForm("schema_name")
	if schemaName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "schema_name is required"})
		return
	}
	err := h.myService.DeleteSchema(schemaName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "schema is deleted successfully",
		"schema_name": schemaName,
	})
}

func (h *Handler) CreateSchema(c *gin.Context) {
	schemaName := c.PostForm("schema_name")
	if schemaName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "schema_name is required"})
		return
	}
	err := h.myService.CreateSchema(schemaName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "schema is created successfully",
		"schema_name": schemaName,
	})
}

func (h *Handler) BackupDB(c *gin.Context) {
	dbName := c.PostForm("db_name")
	err := h.myService.BackupDatabase(dbName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}

func (h *Handler) CronBackupDB(c *gin.Context) {
	dbName, respQuery := c.PostForm("db_name"), c.PostForm("query")
	var mapQuery map[string]interface{}
	json.Unmarshal([]byte(respQuery), &mapQuery)
	var query []string
	for _, str := range mapQuery["data"].([]interface{}) {
		query = append(query, str.(string))
	}
	err := h.myService.CronBackupDataBase(dbName, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}

func (h *Handler) RemoveCronJob(c *gin.Context) {
	id := c.PostForm("job_id")
	h.myService.RemoveCronJob(id)
	c.JSON(http.StatusOK, gin.H{"message": "successful"})
}

func (h *Handler) GetJobList(c *gin.Context) {
	jobs := h.myService.GetJobList()
	c.JSON(http.StatusOK, jobs)
}

func (h *Handler) CronStart(c *gin.Context) {
	h.myService.CronStart()
	c.Done()
}

func (h *Handler) CronStop(c *gin.Context) {
	h.myService.CronStop()
	c.Done()
}

// 권한 검증이 필요할 수도 있음.
func (h *Handler) Listen(c *gin.Context) {
	err := h.myService.StartListenerManager()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "listen successfully"})
}

// 권한 검증이 필요할 수도 있음.
func (h *Handler) Unlisten(c *gin.Context) {
	err := h.myService.StopListenerManager()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unlisten successfully"})
}

func (h *Handler) ReadAllTablesBySchema(c *gin.Context) {
	schema_name := c.PostForm("schema_name")
	if schema_name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_name is required"})
		return
	}
	tables, err := h.myService.ReadAllTablesBySchemaNamd(schema_name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": map[string][]string{schema_name: tables}})
}

func (h *Handler) ReadAllRecordByTableName(c *gin.Context) {
	tableName := c.PostForm("table_name")
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_name is required"})
		return
	}
	records, err := h.myService.ReadAllRecordByTableName(tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": records})
}

func (h *Handler) CreateTableCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	tableName := c.PostForm("table_name")
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_name is required"})
		return
	}
	savePath := fmt.Sprintf("./resource/%s", file.Filename)
	c.SaveUploadedFile(file, savePath)

	err = h.myService.CreateTableFromCSV(savePath, tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "table is created successfully",
		"saved_file": file.Filename,
		"table_name": tableName,
	})
}

func (h *Handler) CreateTableJSON(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	tableName := c.PostForm("table_name")
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_name is required"})
		return
	}
	savePath := fmt.Sprintf("./resource/%s", file.Filename)
	c.SaveUploadedFile(file, savePath)

	err = h.myService.CreateTableFromJSON(savePath, tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "table is created successfully",
		"saved_file": file.Filename,
		"table_name": tableName,
	})
}

func (h *Handler) CreateTableExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	tableName := c.PostForm("table_name")
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_name is required"})
		return
	}
	savePath := fmt.Sprintf("./resource/%s", file.Filename)
	c.SaveUploadedFile(file, savePath)

	err = h.myService.CreateTableFromExcel(savePath, tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "table is created successfully",
		"saved_file": file.Filename,
		"table_name": tableName,
	})
}

func (h *Handler) CreateNormalizeTableCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	savePath := fmt.Sprintf("./resource/%s", file.Filename)
	c.SaveUploadedFile(file, savePath)

	retStr, err := h.myService.CreateNormalizeTableFromCSV(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "table is normalized successfully",
		"saved_file": file.Filename,
		"info":       retStr,
	})
}

func (h *Handler) DeleteTable(c *gin.Context) {

	tableName := c.PostForm("table_name")
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_name is required"})
		return
	}

	err := h.myService.DropTableByTableName(tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "table is deleted successfully",
		"table_name": tableName,
	})
}

func (h *Handler) ExportTableJSON(c *gin.Context) {

	tableName := c.PostForm("table_name")
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_name is required"})
		return
	}

	err, jsonFilePath := h.myService.ExportTableToJSON(tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.FileAttachment(jsonFilePath, fmt.Sprintf("%s_table.json", tableName))
}

func (h *Handler) ExportTableCSV(c *gin.Context) {

	tableName := c.PostForm("table_name")
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_name is required"})
		return
	}

	err, csvFilePath := h.myService.ExportTableToCSV(tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.FileAttachment(csvFilePath, fmt.Sprintf("%s_table.csv", tableName))
}
