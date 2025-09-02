package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"cliscore/internal/config"
	"cliscore/internal/models"
)

type APIClient struct {
	config *config.Config
}

func New(cfg *config.Config) *APIClient {
	return &APIClient{config: cfg}
}

func (c *APIClient) Search(req *models.SearchRequest, apiKey string) (*models.SearchResponse, error) {
	return makeRequest[models.SearchResponse]("POST", "/search", req, apiKey)
}

func (c *APIClient) Count(req *models.CountRequest, apiKey string) (*models.DetailedCountResponse, error) {
	return makeRequest[models.DetailedCountResponse]("POST", "/count/detailed", req, apiKey)
}

func (c *APIClient) ValidateAPIKey(apiKey string) error {
	validationReq := map[string]string{
		"apiKey": apiKey,
	}
	
	cfg := c.config
	url := cfg.BaseURL + "/validate"
	
	jsonData, err := json.Marshal(validationReq)
	if err != nil {
		return fmt.Errorf("error marshaling request: %v", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	
	// Success is 200 OK with no response body
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	
	// For errors, read the response body
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid API key")
	}
	
	if resp.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("invalid request: %s", string(body))
	}
	
	return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
}

func (c *APIClient) GetMachineInfo(uuid string, apiKey string) (*models.MachineInfoResponse, error) {
	cfg := c.config
	url := fmt.Sprintf("%s/machineinfo?uuid=%s", cfg.BaseURL, uuid)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	
	var result models.MachineInfoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}
	
	return &result, nil
}

func (c *APIClient) DownloadFile(uuid, filePath, apiKey string, outputPath string) error {
	cfg := c.config
	
	// Build URL with query parameters
	url := fmt.Sprintf("%s/download?uuid=%s", cfg.BaseURL, uuid)
	if filePath != "" {
		url += fmt.Sprintf("&file=%s", filePath)
	}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	
	// Determine output filename
	var fileName string
	if filePath != "" {
		fileName = filepath.Base(filePath)
	} else {
		fileName = fmt.Sprintf("%s.zip", uuid)
	}
	
	// Use provided output path or generate one
	var finalPath string
	if outputPath != "" {
		finalPath = outputPath
	} else {
		finalPath = fileName
	}
	
	// Create the file
	outFile, err := os.Create(finalPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outFile.Close()
	
	// Copy the response body to the file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}
	
	fmt.Printf("File downloaded successfully: %s\n", finalPath)
	return nil
}

func (c *APIClient) GetCredits(apiKey string) (*models.CreditsResponse, error) {
	req := &models.ApiKeyValidation{
		ApiKey: apiKey,
	}
	return makeRequest[models.CreditsResponse]("POST", "/credits", req, "")
}

func makeRequest[T any](method, endpoint string, data interface{}, apiKey string) (*T, error) {
	cfg := config.Load()
	url := cfg.BaseURL + endpoint

	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request: %v", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	return &result, nil
}