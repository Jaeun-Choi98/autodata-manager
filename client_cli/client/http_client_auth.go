package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func (hc *HttpClient) UpdateUser(filePath string) (map[string]interface{}, error) {

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", filePath)
	}
	defer file.Close()

	// formData: { "file" : file}

	// 파일 필드 이름은 "file"로 지정
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %s", filePath)
	}

	// 파일 내용을 formData에 복사
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// 요청 본문 마무리
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// e.g., http://localhost:8080/service/create/csv
	url := fmt.Sprintf("%s/user/update", hc.baseUrl)
	// HTTP POST 요청 보내기
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Content-Type 헤더 설정
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 클라이언트로 요청 보내기
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("404 Not Found: the requested URL %s does not exist", url)
	}

	var response map[string]interface{}
	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)
	err = json.Unmarshal(buf.Bytes(), &response)
	if err != nil {
		return response, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// 응답 코드 확인
	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("received non-OK response: %v(%v)", resp.Status, response["error"])
	}

	return response, nil
}

func (hc *HttpClient) RegisterUser(filePath string) (map[string]interface{}, error) {

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", filePath)
	}
	defer file.Close()

	// formData: { "file" : file}

	// 파일 필드 이름은 "file"로 지정
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %s", filePath)
	}

	// 파일 내용을 formData에 복사
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// 요청 본문 마무리
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// e.g., http://localhost:8080/service/create/csv
	url := fmt.Sprintf("%s/user/register", hc.baseUrl)
	// HTTP POST 요청 보내기
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Content-Type 헤더 설정
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 클라이언트로 요청 보내기
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("404 Not Found: the requested URL %s does not exist", url)
	}

	var response map[string]interface{}
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	//io.Copy(&buf, resp.Body)
	err = json.Unmarshal(buf.Bytes(), &response)
	if err != nil {
		return response, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// 응답 코드 확인
	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("received non-OK response: %v(%v)", resp.Status, response["error"])
	}

	return response, nil
}
