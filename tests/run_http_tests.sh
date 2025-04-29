#!/bin/bash

# HTTP API Test Runner for Suasor
# This script executes the HTTP API tests using curl

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}===== Suasor HTTP API Test Runner =====${NC}"
echo "Starting tests at $(date)"

# Environment variables - modify these for your environment
export BASE_URL="http://localhost:8080/api"
export TEST_ADMIN_USER="admin@example.com"
export TEST_ADMIN_PASSWORD="password"
export TEST_USER="user@example.com"
export TEST_PASSWORD="password"
export TEST_NEW_USER="newuser@example.com"

# Check if http client is available
if ! command -v curl &> /dev/null; then
    echo -e "${RED}Error: curl is not installed. Please install curl to run the tests.${NC}"
    exit 1
fi

# Function to run a test suite
run_test_suite() {
    local suite=$1
    local file=$2
    echo -e "${YELLOW}Running $suite test suite...${NC}"
    
    # Run the HTTP test file using your HTTP client
    # This is a placeholder - you would need to replace this with the actual command
    # to run your .http files based on your IDE/tool
    
    echo "To run this test suite manually, use your HTTP client to run:"
    echo "  $file"
    
    echo -e "${GREEN}$suite tests completed${NC}"
    echo
}

# Run core tests
run_test_suite "Core" "/home/faiyt/codebase/suasor/backend/tests/_runners/core.http"

# Run client tests
run_test_suite "Client" "/home/faiyt/codebase/suasor/backend/tests/_runners/clients.http"

# Run all tests
run_test_suite "All" "/home/faiyt/codebase/suasor/backend/tests/_runners/all.http"

echo -e "${BLUE}===== All tests completed =====${NC}"
echo "Tests finished at $(date)"

echo
echo -e "${YELLOW}NOTE: This script is a template. To run the actual tests:${NC}"
echo "1. Use an HTTP client that supports .http files (like VS Code with REST Client extension,"
echo "   IntelliJ IDEA, or a similar tool)"
echo "2. Open the test runner files in the _runners directory"
echo "3. Configure your environment variables in http-client.env.json"
echo "4. Run the tests using your HTTP client's interface"