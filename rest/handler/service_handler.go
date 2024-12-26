package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ReadAllTablesBySchema(c *gin.Context) {
	schema_name := c.PostForm("schema_name")
	if schema_name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_name is required"})
		return
	}
	tables, err := h.myService.ReadAllTablesBySchemaNamd(schema_name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tables})
}

func (h *Handler) ReadAllRecordByTableName(c *gin.Context) {
	tableName := c.PostForm("table_name")
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_name is required"})
		return
	}
	records, err := h.myService.ReadAllRecordByTableName(tableName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
