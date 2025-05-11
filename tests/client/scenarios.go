package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// TestScenario represents a scenario to test
type TestScenario struct {
	Name        string
	Description string
	Execute     func(*Client) error
}

// GetAvailableScenarios returns all available test scenarios
func GetAvailableScenarios() []TestScenario {
	return []TestScenario{
		{
			Name:        "health",
			Description: "Tests the API health endpoint",
			Execute:     testHealthCheck,
		},
		{
			Name:        "auth",
			Description: "Tests user registration, login, and JWT authentication",
			Execute:     testUserRegistrationAndAuth,
		},
		{
			Name:        "account",
			Description: "Tests retrieving account info and balance",
			Execute:     testAccountManagement,
		},
		{
			Name:        "transaction",
			Description: "Tests creating movements and listing transactions",
			Execute:     testTransactionOperations,
		},
		{
			Name:        "transfer",
			Description: "Tests creating multiple users and transferring funds between accounts",
			Execute:     testMultipleUsersAndTransfers,
		},
	}
}

// RunScenarios executes test scenarios based on command-line flags
func RunScenarios() {
	// Define command-line flags
	var (
		listFlag     bool
		scenarioFlag string
		allFlag      bool
	)

	flag.BoolVar(&listFlag, "list", false, "List all available test scenarios")
	flag.StringVar(&scenarioFlag, "test", "", "Specify which test scenario to run (comma-separated)")
	flag.BoolVar(&allFlag, "all", false, "Run all test scenarios")
	flag.Parse()

	// Get all scenarios
	scenarios := GetAvailableScenarios()

	// Handle -list flag: list all available scenarios and exit
	if listFlag {
		fmt.Println("Available test scenarios:")
		for i, scenario := range scenarios {
			fmt.Printf("%d. %s - %s\n", i+1, scenario.Name, scenario.Description)
		}
		os.Exit(0)
	}

	// Create client
	client := NewClient()

	fmt.Println("Starting API test scenarios...")

	// If no flags are specified, default to running all scenarios
	if !allFlag && scenarioFlag == "" {
		allFlag = true
	}

	if allFlag {
		// Run all scenarios
		for i, scenario := range scenarios {
			runSingleScenario(client, scenario, i, len(scenarios))
		}
	} else {
		// Run selected scenarios
		requestedScenarios := strings.Split(scenarioFlag, ",")
		for _, requestedName := range requestedScenarios {
			requestedName = strings.TrimSpace(requestedName)
			scenarioFound := false

			for i, scenario := range scenarios {
				if scenario.Name == requestedName {
					runSingleScenario(client, scenario, i, len(requestedScenarios))
					scenarioFound = true
					break
				}
			}

			if !scenarioFound {
				fmt.Printf("❌ Scenario '%s' not found\n", requestedName)
			}
		}
	}

	// Save test results to JSON file
	if err := client.SaveResults(); err != nil {
		fmt.Printf("Failed to save test results: %v\n", err)
	} else {
		fmt.Printf("\nTest results saved to %s and %s\n", logFile, jsonFile)
	}

	// Print summary
	client.PrintSummary()
}

// runSingleScenario runs a single test scenario
func runSingleScenario(client *Client, scenario TestScenario, index int, total int) {
	fmt.Printf("\n=== Scenario %d/%d: %s ===\n", index+1, total, scenario.Name)
	fmt.Printf("Description: %s\n\n", scenario.Description)

	// Execute scenario
	startTime := time.Now()
	err := scenario.Execute(client)
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ Scenario failed: %v\n", err)
	} else {
		fmt.Printf("✅ Scenario completed successfully in %.2f seconds\n", duration.Seconds())
	}

	// Add delay between scenarios
	if index < total-1 {
		time.Sleep(1 * time.Second)
	}
}

// testHealthCheck tests the health check endpoint
func testHealthCheck(client *Client) error {
	success, err := client.HealthCheck()
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	if !success {
		return fmt.Errorf("health check returned unsuccessful status")
	}

	fmt.Printf("✓ API health check passed\n")
	return nil
}

// testUserRegistrationAndAuth tests user registration and authentication
func testUserRegistrationAndAuth(client *Client) error {
	// Generate unique test user
	timestamp := time.Now().Unix()
	testEmail := fmt.Sprintf("test%d@example.com", timestamp)
	testUsername := fmt.Sprintf("user%d", timestamp)
	testPassword := "SecurePass123!"

	// Test signup
	user, err := client.Signup(
		testEmail,
		testPassword,
		testUsername,
		"Test",
		"User",
		fmt.Sprintf("FC%d", timestamp),
	)
	if err != nil {
		return fmt.Errorf("signup failed: %w", err)
	}
	fmt.Printf("✓ Created user with ID %s\n", user["id"])

	// Test login
	token, err := client.Login(testEmail, testPassword)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}
	fmt.Printf("✓ Login successful, received token\n")

	// Test invalid login
	_, err = client.Login(testEmail, "WrongPassword")
	if err == nil {
		return fmt.Errorf("invalid login should fail but succeeded")
	}
	fmt.Printf("✓ Invalid login correctly rejected\n")

	// Reset token for other tests
	client.SetToken(token)

	// Test invalid registration with existing email
	_, err = client.Signup(
		testEmail,
		testPassword,
		fmt.Sprintf("another%d", timestamp),
		"Another",
		"User",
		fmt.Sprintf("FC%d_2", timestamp),
	)
	if err == nil {
		return fmt.Errorf("signup with existing email should fail but succeeded")
	}
	fmt.Printf("✓ Duplicate email registration correctly rejected\n")

	return nil
}

// testAccountManagement tests account management operations
func testAccountManagement(client *Client) error {
	// Test get balance
	balance, err := client.GetBalance()
	if err != nil {
		return fmt.Errorf("get balance failed: %w", err)
	}

	fmt.Printf("✓ Account balance: %s %s\n", balance["balance"], balance["currency"])

	// Get balance without authentication
	unauthClient := NewClient()
	_, err = unauthClient.GetBalance()
	if err == nil {
		return fmt.Errorf("unauthenticated balance check should fail but succeeded")
	}
	fmt.Printf("✓ Unauthenticated balance check correctly rejected\n")

	return nil
}

// testTransactionOperations tests transaction operations
func testTransactionOperations(client *Client) error {
	// Test create credit movement
	creditAmount := "100.00"
	creditMovement, err := client.CreateMovement(creditAmount, "credit", "Test deposit")
	if err != nil {
		return fmt.Errorf("create credit movement failed: %w", err)
	}
	fmt.Printf("✓ Created credit amount of %s\n", creditAmount)
	fmt.Printf("✓ Created credit movement of %s\n", creditMovement)

	// Check updated balance
	balance, err := client.GetBalance()
	if err != nil {
		return fmt.Errorf("get balance after credit failed: %w", err)
	}
	fmt.Printf("✓ Balance after credit: %s\n", balance["balance"])

	// Test create debit movement
	debitAmount := "30.00"
	debitMovement, err := client.CreateMovement(debitAmount, "debit", "Test withdrawal")
	if err != nil {
		return fmt.Errorf("create debit movement failed: %w", err)
	}
	fmt.Printf("✓ Created debit amount of %s\n", debitAmount)
	fmt.Printf("✓ Created debit movement of %s\n", debitMovement)

	// Check updated balance
	balance, err = client.GetBalance()
	if err != nil {
		return fmt.Errorf("get balance after debit failed: %w", err)
	}
	fmt.Printf("✓ Balance after debit: %s\n", balance["balance"])

	// Test get movements
	movements, err := client.GetMovements(1, 10)
	if err != nil {
		return fmt.Errorf("get movements failed: %w", err)
	}

	// Check movement list
	data, ok := movements["data"].([]interface{})
	if !ok {
		return fmt.Errorf("movements data is not an array")
	}
	fmt.Printf("✓ Retrieved %d movements\n", len(data))

	// Try to create a debit movement exceeding the balance
	largeAmount := "1000.00"
	_, err = client.CreateMovement(largeAmount, "debit", "Test excessive withdrawal")
	if err == nil {
		return fmt.Errorf("excessive withdrawal should fail but succeeded")
	}
	fmt.Printf("✓ Excessive withdrawal correctly rejected\n")

	// Try to create a movement with invalid type
	_, err = client.CreateMovement("10.00", "invalid", "Test invalid type")
	if err == nil {
		return fmt.Errorf("movement with invalid type should fail but succeeded")
	}
	fmt.Printf("✓ Invalid movement type correctly rejected\n")

	return nil
}

// testMultipleUsersAndTransfers tests operations with multiple users
func testMultipleUsersAndTransfers(client *Client) error {
	// Store the first user's token
	firstUserToken := client.token

	// Create a second user
	timestamp := time.Now().Unix()
	secondEmail := fmt.Sprintf("second%d@example.com", timestamp)
	secondUsername := fmt.Sprintf("second%d", timestamp)
	secondPassword := "SecurePass456!"

	// Register second user
	secondUser, err := client.Signup(
		secondEmail,
		secondPassword,
		secondUsername,
		"Second",
		"User",
		fmt.Sprintf("FC2%d", timestamp),
	)
	if err != nil {
		return fmt.Errorf("second user signup failed: %w", err)
	}
	fmt.Printf("✓ Created second user with ID %s\n", secondUser["id"])

	// Login as second user
	_, err = client.Login(secondEmail, secondPassword)
	if err != nil {
		return fmt.Errorf("second user login failed: %w", err)
	}
	fmt.Printf("✓ Second user login successful\n")

	// Add credit to second user's account
	_, err = client.CreateMovement("200.00", "credit", "Initial deposit")
	if err != nil {
		return fmt.Errorf("second user deposit failed: %w", err)
	}
	fmt.Printf("✓ Added funds to second user's account\n")

	// Get second user's account details to get account ID
	secondBalance, err := client.GetBalance()
	if err != nil {
		return fmt.Errorf("get second user balance failed: %w", err)
	}
	secondAccountID, ok := secondBalance["account_id"].(string)
	if !ok {
		return fmt.Errorf("couldn't get second user's account ID")
	}
	fmt.Printf("✓ Second user account ID: %s\n", secondAccountID)

	// Switch back to first user
	client.SetToken(firstUserToken)

	// Transfer from first user to second user
	transferAmount := "25.00"
	_, err = client.CreateTransfer(secondAccountID, transferAmount, "Test transfer")
	if err != nil {
		return fmt.Errorf("transfer failed: %w", err)
	}
	fmt.Printf("✓ Transferred %s to second user\n", transferAmount)

	// Check first user balance after transfer
	firstBalance, err := client.GetBalance()
	if err != nil {
		return fmt.Errorf("get first user balance after transfer failed: %w", err)
	}
	fmt.Printf("✓ First user balance after transfer: %s\n", firstBalance["balance"])

	// Get transfer history
	transfers, err := client.GetTransfers(1, 10)
	if err != nil {
		return fmt.Errorf("get transfers failed: %w", err)
	}

	// Check transfer list
	transferData, ok := transfers["data"].([]interface{})
	if !ok {
		return fmt.Errorf("transfers data is not an array")
	}
	fmt.Printf("✓ Retrieved %d transfers\n", len(transferData))

	// Try transfer with invalid account ID
	_, err = client.CreateTransfer("invalid-uuid", "10.00", "Invalid transfer")
	if err == nil {
		return fmt.Errorf("transfer with invalid account ID should fail but succeeded")
	}
	fmt.Printf("✓ Transfer with invalid account ID correctly rejected\n")

	// Try transfer with excessive amount
	_, err = client.CreateTransfer(secondAccountID, "1000.00", "Excessive transfer")
	if err == nil {
		return fmt.Errorf("transfer with excessive amount should fail but succeeded")
	}
	fmt.Printf("✓ Excessive transfer correctly rejected\n")

	return nil
}
