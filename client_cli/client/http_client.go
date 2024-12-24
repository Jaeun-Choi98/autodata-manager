package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type HttpClient struct {
	client *http.Client
}

func NewClient() ClientInterface {
	return &HttpClient{client: &http.Client{}}
}

func (hc *HttpClient) ExportTable(url, tableName, extension string) error {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// formData: { "table_name" : tableName }
	err := writer.WriteField("table_name", tableName)
	if err != nil {
		return fmt.Errorf("failed to write table name: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	url = fmt.Sprintf("%s/%s", url, extension)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to export request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := hc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("received non-OK response: %s", resp.Status)
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	exportFile, _ := os.Create(tableName + ".json")
	defer exportFile.Close()

	io.Copy(exportFile, resp.Body)

	log.Printf("export '%s' table", tableName)
	return nil
}

func (hc *HttpClient) DropTable(url, tableName string) error {

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// formData: { "table_name" : tableName }
	err := writer.WriteField("table_name", tableName)
	if err != nil {
		return fmt.Errorf("failed to write table name: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to drop request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := hc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		log.Printf("failed to decode json response: %w", err)
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("received non-OK response: %s, error: %s", resp.Status, response["error"])
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	log.Printf("response: %s", response)
	return nil
}

func (hc *HttpClient) NormalizeTable(url, filePath, extension string) error {

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// formData: { "file" : file }
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return fmt.Errorf("failed to normalize form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	url = fmt.Sprintf("%s/%s", url, extension)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := hc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		log.Printf("failed to decode json response: %w", err)
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("received non-OK response: %s, error: %s", resp.Status, response["error"])
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	log.Printf("response: %s", response)
	return nil
}

func (hc *HttpClient) MakeTable(url, filePath, tableName, extension string) error {

	// 요청 본문을 작성할 버퍼
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 파일을 formData에 추가
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// formData: { "file" : file, "table_name" : tableName}

	// 파일 필드 이름은 "file"로 지정
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	// 파일 내용을 formData에 복사
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// 테이블 이름을 formData에 추가
	err = writer.WriteField("table_name", tableName)
	if err != nil {
		return fmt.Errorf("failed to write table name: %w", err)
	}

	// 요청 본문 마무리
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	// e.g., http://localhost:8080/service/create/csv
	url = fmt.Sprintf("%s/%s", url, extension)
	// HTTP POST 요청 보내기
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Content-Type 헤더 설정
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 클라이언트로 요청 보내기
	resp, err := hc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		log.Printf("failed to decode json response: %w", err)
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// 응답 코드 확인
	if resp.StatusCode != http.StatusOK {
		log.Printf("received non-OK response: %s, error: %s", resp.Status, response["error"])
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	// 응답 출력
	log.Printf("Response: %s", response)
	return nil
}
