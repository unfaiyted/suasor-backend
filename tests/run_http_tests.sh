#!/bin/bash
# HTTP API Test Runner for Suasor
# Uses a minimal Neovim config to run Kulala.nvim in headless mode

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

# Path to Neovim binary
NVIM_BIN="/home/faiyt/Applications/nvim-linux-arm64/bin/nvim"

# Create a temporary directory for the minimal config
TEMP_CONFIG_DIR=$(mktemp -d)
mkdir -p "$TEMP_CONFIG_DIR/lua/kulala_runner"

# Create a simple plugin to run the tests
cat > "$TEMP_CONFIG_DIR/lua/kulala_runner/init.lua" << 'EOF'
local M = {}

function M.run_tests(file)
  vim.cmd('edit ' .. file)
  
  -- Wait for file to load
  vim.defer_fn(function()
    print("Attempting to run requests in " .. file)
    
    -- Try to run requests directly
    local ok, result = pcall(function()
      -- First approach: direct function call
      local kulala = require("kulala")
      kulala.send_all_requests()
    end)
    
    if not ok then
      print("First approach failed, trying alternative method")
      -- Second approach: use Vim command if available
      pcall(vim.cmd, "KuLalaSendAllRequests")
    end
    
    -- Exit after tests complete
    vim.defer_fn(function()
      print("Tests completed")
      vim.cmd('qa!')
    end, 8000)  -- Wait 8 seconds for requests to complete
  end, 2000)
end

return M
EOF

# Create a minimal init.lua
cat > "$TEMP_CONFIG_DIR/init.lua" << 'EOF'
-- Basic Neovim settings
vim.opt.compatible = false
vim.opt.loadplugins = true
vim.opt.runtimepath:append('~/.local/share/nvim-dane/lazy/kulala.nvim')
vim.opt.runtimepath:append(vim.fn.expand('~/.local/share/nvim-dane/lazy/plenary.nvim'))
vim.opt.runtimepath:append(vim.fn.expand('~/.local/share/nvim-dane/lazy/nvim-treesitter'))
vim.opt.runtimepath:append(vim.fn.expand('$TEMP_CONFIG_DIR'))

-- Add our test runner to rtp
vim.opt.runtimepath:append(vim.fn.expand('$TEMP_CONFIG_DIR'))

-- Enable filetype detection
vim.cmd('filetype plugin on')
vim.cmd('syntax enable')

-- Access test runner
local runner = require("kulala_runner")

-- Expose the run_tests function globally
_G.run_tests = runner.run_tests
EOF

# Function to run a test suite
run_test_suite() {
    local suite=$1
    local file=$2
    echo -e "${YELLOW}Running $suite test suite...${NC}"
    
    # Run Neovim with our minimal config
    TEMP_CONFIG_DIR="$TEMP_CONFIG_DIR" NVIM_APPNAME=nvim-dane eval "$NVIM_BIN --headless -u $TEMP_CONFIG_DIR/init.lua -c 'lua run_tests(\"$file\")'"
    
    # Check exit status
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}$suite tests completed successfully${NC}"
    else
        echo -e "${RED}$suite tests failed${NC}"
    fi
    echo
}

# Run tests
run_test_suite "Core" "/home/faiyt/codebase/suasor/backend/tests/_runners/core.http"
run_test_suite "Client" "/home/faiyt/codebase/suasor/backend/tests/_runners/clients.http"
run_test_suite "All" "/home/faiyt/codebase/suasor/backend/tests/_runners/all.http"

# Clean up
rm -rf "$TEMP_CONFIG_DIR"

echo -e "${BLUE}===== All tests completed =====${NC}"
echo "Tests finished at $(date)"
