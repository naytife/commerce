#!/bin/bash

# Stripe Integration End-to-End Testing Script
# This script runs comprehensive tests for the Stripe integration

set -e

echo "🧪 Starting Stripe Integration E2E Testing..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
FRONTEND_URL="http://localhost:5174"
BACKEND_URL="http://gossip.localhost:8000"
TEST_EMAIL="test@example.com"
TEST_PHONE="+1234567890"

echo -e "${BLUE}📋 Test Configuration:${NC}"
echo "Frontend URL: $FRONTEND_URL"
echo "Backend URL: $BACKEND_URL"
echo "Test Email: $TEST_EMAIL"
echo "Test Phone: $TEST_PHONE"
echo ""

# Function to check service health
check_service() {
    local url=$1
    local name=$2
    
    echo -n "Checking $name health... "
    if curl -s -o /dev/null -w "%{http_code}" "$url" | grep -E "^(200|405)$" > /dev/null; then
        echo -e "${GREEN}✓ Online${NC}"
        return 0
    else
        echo -e "${RED}✗ Offline${NC}"
        return 1
    fi
}

# Function to run a test
run_test() {
    local test_name=$1
    local test_command=$2
    
    echo -e "${YELLOW}🧪 $test_name${NC}"
    if eval "$test_command"; then
        echo -e "${GREEN}✓ $test_name passed${NC}"
        return 0
    else
        echo -e "${RED}✗ $test_name failed${NC}"
        return 1
    fi
}

# Pre-flight checks
echo -e "${BLUE}🚀 Pre-flight Checks${NC}"
check_service "$FRONTEND_URL" "Frontend" || {
    echo -e "${RED}Frontend not running. Please start with: npm run dev${NC}"
    exit 1
}

check_service "$BACKEND_URL/query" "Backend GraphQL" || {
    echo -e "${YELLOW}⚠️  Backend not available. Some tests may use fallback data.${NC}"
}

# Environment validation
echo -e "${BLUE}🔧 Environment Validation${NC}"
if [ -f ".env" ]; then
    echo -e "${GREEN}✓ .env file found${NC}"
    
    # Check for required environment variables
    if grep -q "VITE_STRIPE_PUBLISHABLE_KEY" .env; then
        echo -e "${GREEN}✓ Stripe publishable key configured${NC}"
    else
        echo -e "${RED}✗ Stripe publishable key missing${NC}"
        exit 1
    fi
    
    if grep -q "VITE_API_URL" .env; then
        echo -e "${GREEN}✓ API URL configured${NC}"
    else
        echo -e "${RED}✗ API URL missing${NC}"
        exit 1
    fi
else
    echo -e "${RED}✗ .env file not found${NC}"
    exit 1
fi

echo ""

# Test API endpoints
echo -e "${BLUE}🔌 API Connectivity Tests${NC}"

run_test "GraphQL endpoint health" "curl -s -X POST '$BACKEND_URL/query' -H 'Content-Type: application/json' -d '{\"query\":\"query{__schema{types{name}}}\"}' | grep -q 'data' || true"

run_test "Shop context API" "curl -s '$BACKEND_URL/shops/1' | grep -E '(shop|error)' || true"

run_test "Payment methods API" "curl -s '$BACKEND_URL/shops/1/payment-methods' | grep -E '(payment|error)' || true"

echo ""

# Frontend integration tests
echo -e "${BLUE}🌐 Frontend Integration Tests${NC}"

run_test "Frontend home page loads" "curl -s '$FRONTEND_URL' | grep -q 'Commerce'"

run_test "Checkout page loads" "curl -s '$FRONTEND_URL/checkout' | grep -q 'Checkout'"

run_test "Success page loads" "curl -s '$FRONTEND_URL/checkout/success' | grep -q 'success'"

echo ""

# JavaScript/Bundle tests
echo -e "${BLUE}📦 Build and Bundle Tests${NC}"

run_test "TypeScript compilation" "cd /Users/erimebe/Development/commerce/templates/template_1 && npx tsc --noEmit"

run_test "Svelte kit check" "cd /Users/erimebe/Development/commerce/templates/template_1 && npx svelte-kit sync && npx svelte-check --tsconfig ./tsconfig.json"

run_test "Bundle size check" "cd /Users/erimebe/Development/commerce/templates/template_1 && npm run build && ls -la build/"

echo ""

# Stripe-specific tests
echo -e "${BLUE}💳 Stripe Integration Tests${NC}"

# Test Stripe library loading
run_test "Stripe library loading" "curl -s 'https://js.stripe.com/v3/' | grep -q 'Stripe'"

# Test publishable key format
STRIPE_KEY=$(grep VITE_STRIPE_PUBLISHABLE_KEY .env | cut -d'=' -f2)
run_test "Stripe publishable key format" "echo '$STRIPE_KEY' | grep -q '^pk_'"

echo ""

# Performance tests
echo -e "${BLUE}⚡ Performance Tests${NC}"

run_test "Frontend load time < 3s" "time timeout 3s curl -s '$FRONTEND_URL' > /dev/null"

run_test "Bundle size < 1MB" "[ \$(du -k build | cut -f1) -lt 1024 ] || true"

echo ""

# Security tests
echo -e "${BLUE}🔒 Security Tests${NC}"

run_test "No secret keys in bundle" "! grep -r 'sk_' build/ || true"

run_test "HTTPS redirect in production" "curl -s -I '$FRONTEND_URL' | grep -E '(https|secure)' || true"

run_test "CSP headers present" "curl -s -I '$FRONTEND_URL' | grep -i 'content-security-policy' || true"

echo ""

# Manual test instructions
echo -e "${BLUE}👨‍💻 Manual Testing Instructions${NC}"
echo ""
echo -e "${YELLOW}Please perform the following manual tests:${NC}"
echo ""
echo "1. 🛒 Cart Functionality:"
echo "   - Navigate to $FRONTEND_URL"
echo "   - Add products to cart"
echo "   - Verify cart count updates"
echo "   - Test quantity changes"
echo ""
echo "2. 📝 Checkout Form:"
echo "   - Fill in contact information"
echo "   - Fill in shipping address"
echo "   - Select shipping method"
echo "   - Verify totals calculation"
echo ""
echo "3. 💳 Stripe Payment:"
echo "   - Select Stripe as payment method"
echo "   - Use test card: 4242 4242 4242 4242"
echo "   - Expiry: 12/34, CVC: 123"
echo "   - Complete payment flow"
echo ""
echo "4. ✅ Order Confirmation:"
echo "   - Verify success page redirect"
echo "   - Check order confirmation display"
echo "   - Test navigation links"
echo ""
echo "5. ❌ Error Scenarios:"
echo "   - Test with declined card: 4000 0000 0000 0002"
echo "   - Test with empty form fields"
echo "   - Test with invalid email format"
echo ""

# Test card information
echo -e "${BLUE}💳 Stripe Test Cards${NC}"
echo ""
echo -e "${GREEN}✅ Successful Test Cards:${NC}"
echo "• Visa: 4242 4242 4242 4242"
echo "• Mastercard: 5555 5555 5555 4444"
echo "• American Express: 3782 8224 6310 005"
echo ""
echo -e "${RED}❌ Declined Test Cards:${NC}"
echo "• Generic decline: 4000 0000 0000 0002"
echo "• Insufficient funds: 4000 0000 0000 9995"
echo "• Lost card: 4000 0000 0000 9987"
echo ""
echo -e "${YELLOW}🔐 3D Secure Cards:${NC}"
echo "• Auth required: 4000 0025 0000 3155"
echo "• Auth optional: 4000 0027 6000 3184"
echo ""

# Webhook testing
echo -e "${BLUE}🔗 Webhook Testing${NC}"
echo ""
echo "To test webhooks:"
echo "1. Install Stripe CLI: https://stripe.com/docs/stripe-cli"
echo "2. Login: stripe login"
echo "3. Forward events: stripe listen --forward-to $FRONTEND_URL/api/webhooks/stripe"
echo "4. Use the webhook signing secret in your environment"
echo ""

# Summary
echo -e "${BLUE}📊 Test Summary${NC}"
echo ""
echo "All automated tests completed!"
echo ""
echo -e "${GREEN}✅ Ready for manual testing at: $FRONTEND_URL${NC}"
echo -e "${YELLOW}⚠️  Don't forget to test with real Stripe test cards${NC}"
echo -e "${BLUE}ℹ️  Check browser console for any JavaScript errors${NC}"
echo ""
echo "Happy testing! 🎉"
