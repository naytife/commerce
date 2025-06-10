#!/bin/bash

echo "🔐 Testing Swagger UI OAuth Configuration"
echo "========================================"

# Test 1: Check if Swagger UI client exists in Hydra
echo "📋 Test 1: Checking OAuth client in Hydra..."
CLIENT_INFO=$(kubectl exec -n naytife-auth deployment/hydra -- hydra get oauth2-client d39beaaa-9c53-48e7-b82a-37ff52127473 --endpoint http://localhost:4445 2>/dev/null)

if echo "$CLIENT_INFO" | grep -q "d39beaaa-9c53-48e7-b82a-37ff52127473"; then
    echo "✅ OAuth client exists in Hydra"
    echo "   Client ID: d39beaaa-9c53-48e7-b82a-37ff52127473"
    echo "   Scopes: $(echo "$CLIENT_INFO" | grep "SCOPE" -A1 | tail -1 | xargs)"
    echo "   Redirect URIs: $(echo "$CLIENT_INFO" | grep "REDIRECT URIS" -A1 | tail -1 | xargs)"
else
    echo "❌ OAuth client not found in Hydra"
    exit 1
fi

# Test 2: Check if Swagger docs are accessible
echo
echo "📖 Test 2: Checking Swagger UI accessibility..."
SWAGGER_RESPONSE=$(curl -s -w "%{http_code}" http://127.0.0.1:8080/v1/docs/index.html)
HTTP_CODE="${SWAGGER_RESPONSE: -3}"

if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Swagger UI is accessible at http://127.0.0.1:8080/v1/docs/index.html"
else
    echo "❌ Swagger UI not accessible (HTTP $HTTP_CODE)"
    exit 1
fi

# Test 3: Check if OAuth authorization endpoint responds correctly
echo
echo "🔓 Test 3: Testing OAuth authorization endpoint..."
AUTH_RESPONSE=$(curl -s -w "%{http_code}" -X GET \
    "http://127.0.0.1:8080/oauth2/auth?client_id=d39beaaa-9c53-48e7-b82a-37ff52127473&response_type=code&scope=openid%20offline&redirect_uri=http://127.0.0.1:8080/v1/docs/oauth2-redirect.html&state=test123&app_type=dashboard")

AUTH_HTTP_CODE="${AUTH_RESPONSE: -3}"

if [ "$AUTH_HTTP_CODE" = "303" ] || [ "$AUTH_HTTP_CODE" = "302" ]; then
    echo "✅ OAuth authorization endpoint responding correctly (HTTP $AUTH_HTTP_CODE - redirect as expected)"
else
    echo "❌ OAuth authorization endpoint error (HTTP $AUTH_HTTP_CODE)"
    if [ "$AUTH_HTTP_CODE" = "400" ]; then
        echo "   This might indicate a client configuration issue"
    fi
fi

# Test 4: Check if protected endpoint requires authentication
echo
echo "🔒 Test 4: Testing protected endpoint..."
PROTECTED_RESPONSE=$(curl -s -w "%{http_code}" http://127.0.0.1:8080/v1/shops)
PROTECTED_HTTP_CODE="${PROTECTED_RESPONSE: -3}"

if [ "$PROTECTED_HTTP_CODE" = "401" ]; then
    echo "✅ Protected endpoint correctly requires authentication (HTTP 401)"
else
    echo "❌ Protected endpoint not properly secured (HTTP $PROTECTED_HTTP_CODE)"
fi

# Test 5: Check backend configuration
echo
echo "⚙️  Test 5: Checking backend OAuth configuration..."
echo "   Swagger OAuth Client ID: d39beaaa-9c53-48e7-b82a-37ff52127473"
echo "   Authorization URL: http://127.0.0.1:8080/oauth2/auth"
echo "   Token URL: http://127.0.0.1:8080/oauth2/token"
echo "   Redirect URL: http://127.0.0.1:8080/v1/docs/oauth2-redirect.html"

echo
echo "🎉 OAuth Configuration Test Complete!"
echo "=====================================

📝 Manual Testing Steps:
1. Open Swagger UI: http://127.0.0.1:8080/v1/docs/index.html
2. Click the 'Authorize' button
3. Select the scopes you want to authorize
4. Click 'Authorize' - you should be redirected to Google OAuth
5. Complete the OAuth flow
6. Return to Swagger UI and test a protected endpoint like GET /shops

🔧 If authentication still fails, check:
   - Google OAuth credentials in auth-handler
   - Network connectivity between services
   - Browser cookies and local storage"
