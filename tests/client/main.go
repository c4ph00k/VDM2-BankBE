package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	baseURL  = "http://localhost:8080/api/v1"
	logFile  = "test_results.log"
	jsonFile = "test_results.json"
)

// TestResult represents the result of a test
type TestResult struct {
	TestName   string    `json:"test_name"`
	Endpoint   string    `json:"endpoint"`
	Method     string    `json:"method"`
	StatusCode int       `json:"status_code"`
	Success    bool      `json:"success"`
	Message    string    `json:"message"`
	Response   string    `json:"response,omitempty"`
	Duration   float64   `json:"duration_ms"`
	Timestamp  time.Time `json:"timestamp"`
}

// Client is a simple HTTP client for testing the API
type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
	logger     *log.Logger
	results    []TestResult
}

// NewClient creates a new API test client
func NewClient() *Client {
	// Create logger
	file, err := os.Create(logFile)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		logger:  log.New(file, "", log.LstdFlags),
		results: make([]TestResult, 0),
	}
}

// SetToken sets the JWT token for authenticated requests
func (c *Client) SetToken(token string) {
	c.token = token
}

// makeRequest performs an HTTP request and logs the result
func (c *Client) makeRequest(method, endpoint string, body interface{}, result interface{}) TestResult {
	testName := fmt.Sprintf("%s %s", method, endpoint)
	startTime := time.Now()

	// Create URL
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	// Prepare body if needed
	var bodyReader io.Reader
	if body != nil {
		bodyData, err := json.Marshal(body)
		if err != nil {
			return TestResult{
				TestName:   testName,
				Endpoint:   endpoint,
				Method:     method,
				StatusCode: 0,
				Success:    false,
				Message:    fmt.Sprintf("Failed to marshal request body: %v", err),
				Duration:   0,
				Timestamp:  startTime,
			}
		}
		bodyReader = bytes.NewReader(bodyData)
	}

	// Create request
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return TestResult{
			TestName:   testName,
			Endpoint:   endpoint,
			Method:     method,
			StatusCode: 0,
			Success:    false,
			Message:    fmt.Sprintf("Failed to create request: %v", err),
			Duration:   0,
			Timestamp:  startTime,
		}
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	// Perform request
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime).Milliseconds()
	if err != nil {
		return TestResult{
			TestName:   testName,
			Endpoint:   endpoint,
			Method:     method,
			StatusCode: 0,
			Success:    false,
			Message:    fmt.Sprintf("Request failed: %v", err),
			Duration:   float64(duration),
			Timestamp:  startTime,
		}
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{
			TestName:   testName,
			Endpoint:   endpoint,
			Method:     method,
			StatusCode: resp.StatusCode,
			Success:    false,
			Message:    fmt.Sprintf("Failed to read response body: %v", err),
			Duration:   float64(duration),
			Timestamp:  startTime,
		}
	}

	// Check if the status code indicates success
	success := resp.StatusCode >= 200 && resp.StatusCode < 300

	// Create result
	testResult := TestResult{
		TestName:   testName,
		Endpoint:   endpoint,
		Method:     method,
		StatusCode: resp.StatusCode,
		Success:    success,
		Message:    fmt.Sprintf("Response code: %d", resp.StatusCode),
		Response:   string(respBody),
		Duration:   float64(duration),
		Timestamp:  startTime,
	}

	// Unmarshal response body into result if provided
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			testResult.Success = false
			testResult.Message = fmt.Sprintf("Failed to unmarshal response: %v", err)
		}
	}

	// Log result
	logMsg := fmt.Sprintf("[%s] %s %s - Status: %d, Duration: %.2f ms, Success: %v",
		startTime.Format(time.RFC3339),
		method,
		endpoint,
		testResult.StatusCode,
		testResult.Duration,
		testResult.Success,
	)

	if !testResult.Success {
		logMsg += fmt.Sprintf(", Error: %s", testResult.Message)
	}

	c.logger.Println(logMsg)
	fmt.Println(logMsg)

	// Store result
	c.results = append(c.results, testResult)

	return testResult
}

// SaveResults saves test results to a JSON file
func (c *Client) SaveResults() error {
	file, err := os.Create(jsonFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c.results)
}

// Signup tests user registration
func (c *Client) Signup(email, password, username, firstName, lastName, fiscalCode string) (map[string]interface{}, error) {
	requestBody := map[string]string{
		"email":       email,
		"password":    password,
		"username":    username,
		"first_name":  firstName,
		"last_name":   lastName,
		"fiscal_code": fiscalCode,
	}

	var response map[string]interface{}
	result := c.makeRequest("POST", "/auth/signup", requestBody, &response)

	if !result.Success {
		return nil, fmt.Errorf("signup failed: %s", result.Message)
	}

	return response, nil
}

// Login tests user authentication
func (c *Client) Login(email, password string) (string, error) {
	requestBody := map[string]string{
		"email":    email,
		"password": password,
	}

	var response map[string]interface{}
	result := c.makeRequest("POST", "/auth/login", requestBody, &response)

	if !result.Success {
		return "", fmt.Errorf("login failed: %s", result.Message)
	}

	token, ok := response["token"].(string)
	if !ok {
		return "", fmt.Errorf("no token in response")
	}

	c.SetToken(token)
	return token, nil
}

// GetBalance tests retrieving account balance
func (c *Client) GetBalance() (map[string]interface{}, error) {
	var response map[string]interface{}
	result := c.makeRequest("GET", "/accounts/balance", nil, &response)

	if !result.Success {
		return nil, fmt.Errorf("get balance failed: %s", result.Message)
	}

	return response, nil
}

// CreateMovement tests creating a movement
func (c *Client) CreateMovement(amount, movementType, description string) (map[string]interface{}, error) {
	requestBody := map[string]string{
		"amount":      amount,
		"type":        movementType,
		"description": description,
	}

	var response map[string]interface{}
	result := c.makeRequest("POST", "/accounts/movements", requestBody, &response)

	if !result.Success {
		return nil, fmt.Errorf("create movement failed: %s", result.Message)
	}

	return response, nil
}

// GetMovements tests retrieving account movements
func (c *Client) GetMovements(page, limit int) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/accounts/movements?page=%d&limit=%d", page, limit)

	var response map[string]interface{}
	result := c.makeRequest("GET", endpoint, nil, &response)

	if !result.Success {
		return nil, fmt.Errorf("get movements failed: %s", result.Message)
	}

	return response, nil
}

// CreateTransfer tests creating a transfer
func (c *Client) CreateTransfer(toAccount, amount, description string) (map[string]interface{}, error) {
	requestBody := map[string]string{
		"to_account":  toAccount,
		"amount":      amount,
		"description": description,
	}

	var response map[string]interface{}
	result := c.makeRequest("POST", "/transfers", requestBody, &response)

	if !result.Success {
		return nil, fmt.Errorf("create transfer failed: %s", result.Message)
	}

	return response, nil
}

// GetTransfers tests retrieving transfers
func (c *Client) GetTransfers(page, limit int) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/transfers?page=%d&limit=%d", page, limit)

	var response map[string]interface{}
	result := c.makeRequest("GET", endpoint, nil, &response)

	if !result.Success {
		return nil, fmt.Errorf("get transfers failed: %s", result.Message)
	}

	return response, nil
}

// HealthCheck tests the API health endpoint
func (c *Client) HealthCheck() (bool, error) {
	healthURL := strings.TrimSuffix(c.baseURL, "/api/v1") + "/health"

	req, err := http.NewRequest("GET", healthURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create health check request: %w", err)
	}

	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		c.logger.Printf("[%s] Health check failed: %v", time.Now().Format(time.RFC3339), err)
		return false, err
	}
	defer resp.Body.Close()

	success := resp.StatusCode == http.StatusOK

	testResult := TestResult{
		TestName:   "Health Check",
		Endpoint:   "/health",
		Method:     "GET",
		StatusCode: resp.StatusCode,
		Success:    success,
		Message:    fmt.Sprintf("Health check response code: %d", resp.StatusCode),
		Duration:   float64(duration),
		Timestamp:  startTime,
	}

	// Log result
	logMsg := fmt.Sprintf("[%s] Health Check - Status: %d, Duration: %.2f ms, Success: %v",
		startTime.Format(time.RFC3339),
		testResult.StatusCode,
		testResult.Duration,
		testResult.Success,
	)

	c.logger.Println(logMsg)
	fmt.Println(logMsg)

	// Store result
	c.results = append(c.results, testResult)

	return success, nil
}

// PrintSummary prints a summary of test results
func (c *Client) PrintSummary() {
	// Count success and failures
	var totalTests, successCount, failureCount int
	totalDuration := 0.0

	for _, result := range c.results {
		totalTests++
		totalDuration += result.Duration
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	// Print summary
	fmt.Printf("\n======= TEST SUMMARY =======\n")
	fmt.Printf("Total tests:  %d\n", totalTests)
	fmt.Printf("Successful:   %d (%.1f%%)\n", successCount, float64(successCount)/float64(totalTests)*100)
	fmt.Printf("Failed:       %d (%.1f%%)\n", failureCount, float64(failureCount)/float64(totalTests)*100)
	fmt.Printf("Total time:   %.2f ms\n", totalDuration)
	if totalTests > 0 {
		fmt.Printf("Average time: %.2f ms\n", totalDuration/float64(totalTests))
	}
	fmt.Printf("===========================\n")
}

func main() {
	// Run the tests using command line parameters
	RunScenarios()
}
