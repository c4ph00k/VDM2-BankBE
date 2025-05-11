# VDM2-Bank API Test Client

This is a test client for the VDM2-Bank API. It performs a series of automated tests to verify that the API is working correctly.

## Features

- Tests all major API endpoints or individual APIs as needed
- Health check for API availability
- Runs predefined test scenarios
- Logs all test results with timing information
- Generates a JSON report of test results
- Handles authentication and token management
- Creates test users and accounts automatically
- Tests transfers between accounts

## Test Scenarios

The client can run the following test scenarios:

1. **Health Check** - Tests the API health endpoint
2. **User Registration and Authentication** - Tests creating users and authentication
3. **Account Management** - Tests retrieving account info and balance
4. **Transaction Operations** - Tests creating movements (deposits/withdrawals) and listing transactions
5. **Multiple Users and Transfers** - Tests creating multiple users and transferring funds between accounts

## Usage

First, make sure the VDM2-Bank API is running locally. Then, run the test client with the provided shell script:

```bash
# Run all tests (default behavior)
./run_tests.sh

# List available test scenarios
./run_tests.sh --list
./run_tests.sh -l

# Run specific test scenario(s)
./run_tests.sh --test health
./run_tests.sh -t auth
./run_tests.sh --test account,transaction

# Explicitly run all test scenarios
./run_tests.sh --all
./run_tests.sh -a

# Display help message
./run_tests.sh --help
./run_tests.sh -h
```

Alternatively, you can run the Go code directly:

```bash
# From the repository root:
go run ./tests/client --test health,auth

# Or from the tests/client directory:
go run . --all
go run . --list
go run . --test health
```

## Output

The client will output test progress to the console, and also create:

- `test_results.log` - A text log of all test operations with timestamps
- `test_results.json` - A structured JSON file with detailed test results

## Example Output

```
Starting API test scenarios...

=== Scenario 1/5: health ===
Description: Tests the API health endpoint

✓ API health check passed
✅ Scenario completed successfully in 0.05 seconds

=== Scenario 2/5: auth ===
Description: Tests user registration, login, and JWT authentication

✓ Created user with ID 3fa85f64-5717-4562-b3fc-2c963f66afa6
✓ Login successful, received token
✓ Invalid login correctly rejected
✓ Duplicate email registration correctly rejected
✅ Scenario completed successfully in 0.42 seconds

=== Scenario 3/5: account ===
Description: Tests retrieving account info and balance

✓ Account balance: 0.00 EUR
✓ Unauthenticated balance check correctly rejected
✅ Scenario completed successfully in 0.15 seconds

...

======= TEST SUMMARY =======
Total tests:  24
Successful:   24 (100.0%)
Failed:       0 (0.0%)
Total time:   1532.56 ms
Average time: 63.86 ms
===========================
```

## Troubleshooting

If tests fail, check:

1. The API server is running at `http://localhost:8080/api/v1`
2. The API health endpoint is accessible at `http://localhost:8080/health`
3. The database is running and properly initialized
4. The Redis server is running
5. The configuration is correct

You can modify the `baseURL` constant in `main.go` if your server runs at a different location.