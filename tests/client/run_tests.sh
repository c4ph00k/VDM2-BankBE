#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}==========================================${NC}"
echo -e "${BLUE}  VDM2-Bank API Test Client Runner       ${NC}"
echo -e "${BLUE}==========================================${NC}"

# Parse command-line parameters
ALL_FLAG=""
SCENARIO_FLAG=""
LIST_FLAG=""

print_usage() {
    echo -e "${YELLOW}Usage: $0 [options]${NC}"
    echo -e "Options:"
    echo -e "  -a, --all         Run all test scenarios (default if no options provided)"
    echo -e "  -t, --test TEST   Run specific test scenario(s) (comma-separated)"
    echo -e "  -l, --list        List all available test scenarios and exit"
    echo -e "  -h, --help        Display this help message"
}

# Process command line arguments
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -h|--help)
            print_usage
            exit 0
            ;;
        -a|--all)
            ALL_FLAG="--all"
            shift
            ;;
        -t|--test)
            SCENARIO_FLAG="--test $2"
            shift
            shift
            ;;
        -l|--list)
            LIST_FLAG="--list"
            shift
            ;;
        *)
            # Unknown option
            echo -e "${RED}Unknown option: $1${NC}"
            print_usage
            exit 1
            ;;
    esac
done

# Check if the API is running
echo -e "${BLUE}Checking if API is running...${NC}"
if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health; then
    echo -e "${GREEN}API is running!${NC}"
else
    echo -e "${RED}Error: API is not running at http://localhost:8080${NC}"
    echo -e "${RED}Please start the API server first.${NC}"
    exit 1
fi

# Run the tests
echo -e "${BLUE}Running API tests...${NC}"
cd "$(dirname "$0")"

# Combine all flags
ALL_ARGS="${LIST_FLAG} ${ALL_FLAG} ${SCENARIO_FLAG}"
if [ -z "${ALL_ARGS}" ]; then
    # Default to run all tests if no args provided
    ALL_ARGS="--all"
fi

go run . ${ALL_ARGS}

# Check if tests were successful
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Tests completed!${NC}"
    echo -e "${BLUE}Test results saved to:${NC}"
    echo -e "  - $(pwd)/test_results.log"
    echo -e "  - $(pwd)/test_results.json"
else
    echo -e "${RED}Tests failed!${NC}"
    echo -e "${RED}Please check the log files for details.${NC}"
    exit 1
fi