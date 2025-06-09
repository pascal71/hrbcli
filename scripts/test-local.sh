#!/bin/bash

# Test script for Harbor CLI
# This script helps test the CLI against a local Harbor instance

set -e

HARBOR_URL="${HARBOR_URL:-https://localhost}"
HARBOR_USERNAME="${HARBOR_USERNAME:-admin}"
HARBOR_PASSWORD="${HARBOR_PASSWORD:-Harbor12345}"

echo "🧪 Testing Harbor CLI..."
echo "Harbor URL: $HARBOR_URL"
echo ""

# Build the CLI
echo "📦 Building hrbcli..."
make build

# Path to the binary
HRBCLI="./bin/hrbcli"

# Test version command
echo "📋 Testing version command..."
$HRBCLI version

# Configure Harbor connection
echo ""
echo "🔧 Configuring Harbor connection..."
$HRBCLI config set harbor_url "$HARBOR_URL"
$HRBCLI config set username "$HARBOR_USERNAME"
$HRBCLI config set insecure true
export HARBOR_PASSWORD="$HARBOR_PASSWORD"

# Test configuration
echo ""
echo "📋 Current configuration:"
$HRBCLI config list

# Test project commands
echo ""
echo "🏗️  Testing project commands..."

# List projects
echo "- Listing projects..."
$HRBCLI project list

# Create a test project
TEST_PROJECT="test-project-$(date +%s)"
echo "- Creating test project: $TEST_PROJECT"
$HRBCLI project create "$TEST_PROJECT" --public

# Get project details
echo "- Getting project details..."
$HRBCLI project get "$TEST_PROJECT"

# Update project
echo "- Updating project..."
$HRBCLI project update "$TEST_PROJECT" --auto-scan=true

# Check if project exists
echo "- Checking if project exists..."
$HRBCLI project exists "$TEST_PROJECT"

# List projects with the new one
echo "- Listing projects again..."
$HRBCLI project list

# Delete the test project
echo "- Deleting test project..."
$HRBCLI project delete "$TEST_PROJECT" --force

echo ""
echo "✅ All tests passed!"
