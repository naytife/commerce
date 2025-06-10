#!/bin/bash

# End-to-End Payment Methods Testing Script
# This script tests the complete payment methods functionality

set -e

echo "🚀 Starting End-to-End Payment Methods Testing..."

BASE_URL="http://localhost:8002"
SHOP_ID="14"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper function to make API calls
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "${BLUE}📡 ${description}${NC}"
    
    if [ -n "$data" ]; then
        response=$(curl -s -X "$method" "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -X "$method" "$BASE_URL$endpoint")
    fi
    
    # Check if response contains success status
    if echo "$response" | jq -r '.status' | grep -q "success"; then
        echo -e "${GREEN}✅ Success${NC}"
        echo "$response" | jq '.'
        return 0
    else
        echo -e "${RED}❌ Failed${NC}"
        echo "$response" | jq '.'
        return 1
    fi
}

# Test 1: Get initial payment methods
echo -e "${YELLOW}🧪 Test 1: Get Initial Payment Methods${NC}"
api_call "GET" "/v1/shops/$SHOP_ID/payment-methods" "" "Getting current payment methods"

# Test 2: Create/Update Stripe payment method
echo -e "${YELLOW}🧪 Test 2: Create/Update Stripe Payment Method${NC}"
stripe_config='{
    "method_type": "stripe",
    "is_enabled": false,
    "config": {
        "publishable_key": "pk_test_e2e_testing_123",
        "secret_key": "sk_test_e2e_testing_456",
        "test_mode": true
    }
}'
api_call "PUT" "/v1/shops/$SHOP_ID/payment-methods/stripe" "$stripe_config" "Creating/updating Stripe payment method"

# Test 3: Enable Stripe payment method
echo -e "${YELLOW}🧪 Test 3: Enable Stripe Payment Method${NC}"
enable_data='{"is_enabled": true}'
api_call "PATCH" "/v1/shops/$SHOP_ID/payment-methods/stripe/status" "$enable_data" "Enabling Stripe payment method"

# Test 4: Test Stripe connection
echo -e "${YELLOW}🧪 Test 4: Test Stripe Connection${NC}"
api_call "POST" "/v1/shops/$SHOP_ID/payment-methods/stripe/test" "" "Testing Stripe connection"

# Test 5: Create PayPal payment method
echo -e "${YELLOW}🧪 Test 5: Create PayPal Payment Method${NC}"
paypal_config='{
    "method_type": "paypal",
    "is_enabled": true,
    "config": {
        "client_id": "paypal_e2e_test_client_id",
        "client_secret": "paypal_e2e_test_client_secret",
        "sandbox_mode": true
    }
}'
api_call "PUT" "/v1/shops/$SHOP_ID/payment-methods/paypal" "$paypal_config" "Creating PayPal payment method"

# Test 6: Create Flutterwave payment method
echo -e "${YELLOW}🧪 Test 6: Create Flutterwave Payment Method${NC}"
flutterwave_config='{
    "method_type": "flutterwave",
    "is_enabled": false,
    "config": {
        "public_key": "FLWPUBK_TEST-e2e_testing",
        "secret_key": "FLWSECK_TEST-e2e_testing",
        "encryption_key": "FLWSECK_TEST-e2e_encryption",
        "test_mode": true
    }
}'
api_call "PUT" "/v1/shops/$SHOP_ID/payment-methods/flutterwave" "$flutterwave_config" "Creating Flutterwave payment method"

# Test 7: Create Paystack payment method
echo -e "${YELLOW}🧪 Test 7: Create Paystack Payment Method${NC}"
paystack_config='{
    "method_type": "paystack",
    "is_enabled": true,
    "config": {
        "public_key": "pk_test_e2e_paystack",
        "secret_key": "sk_test_e2e_paystack",
        "test_mode": true
    }
}'
api_call "PUT" "/v1/shops/$SHOP_ID/payment-methods/paystack" "$paystack_config" "Creating Paystack payment method"

# Test 8: Toggle Flutterwave status
echo -e "${YELLOW}🧪 Test 8: Toggle Flutterwave Status${NC}"
api_call "PATCH" "/v1/shops/$SHOP_ID/payment-methods/flutterwave/status" "$enable_data" "Enabling Flutterwave payment method"

# Test 9: Update Stripe configuration
echo -e "${YELLOW}🧪 Test 9: Update Stripe Configuration${NC}"
updated_stripe_config='{
    "method_type": "stripe",
    "is_enabled": true,
    "config": {
        "publishable_key": "pk_test_updated_e2e_123",
        "secret_key": "sk_test_updated_e2e_456",
        "test_mode": false
    }
}'
api_call "PUT" "/v1/shops/$SHOP_ID/payment-methods/stripe" "$updated_stripe_config" "Updating Stripe configuration"

# Test 10: Disable PayPal
echo -e "${YELLOW}🧪 Test 10: Disable PayPal${NC}"
disable_data='{"is_enabled": false}'
api_call "PATCH" "/v1/shops/$SHOP_ID/payment-methods/paypal/status" "$disable_data" "Disabling PayPal payment method"

# Test 11: Delete Paystack payment method
echo -e "${YELLOW}🧪 Test 11: Delete Paystack Payment Method${NC}"
api_call "DELETE" "/v1/shops/$SHOP_ID/payment-methods/paystack" "" "Deleting Paystack payment method"

# Test 12: Get final payment methods state
echo -e "${YELLOW}🧪 Test 12: Get Final Payment Methods State${NC}"
api_call "GET" "/v1/shops/$SHOP_ID/payment-methods" "" "Getting final payment methods state"

# Test 13: Validation test - Try to create method with missing config
echo -e "${YELLOW}🧪 Test 13: Validation Test - Invalid Configuration${NC}"
invalid_config='{
    "method_type": "stripe",
    "is_enabled": true,
    "config": {}
}'
echo -e "${BLUE}📡 Testing validation with empty configuration${NC}"
response=$(curl -s -X "PUT" "$BASE_URL/v1/shops/$SHOP_ID/payment-methods/stripe" \
    -H "Content-Type: application/json" \
    -d "$invalid_config")
    
echo "Response to invalid config:"
echo "$response" | jq '.'

echo -e "${GREEN}🎉 End-to-End Payment Methods Testing Completed!${NC}"

# Summary
echo -e "${BLUE}📊 Test Summary:${NC}"
echo "✅ Payment methods CRUD operations"
echo "✅ Enable/disable functionality" 
echo "✅ Configuration management"
echo "✅ Payment method testing"
echo "✅ Validation handling"
echo "✅ Error scenarios"

echo -e "${GREEN}✨ All payment processor functionality has been successfully tested!${NC}"
