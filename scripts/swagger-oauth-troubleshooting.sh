#!/bin/bash

echo "🔧 Swagger OAuth Troubleshooting"
echo "================================"

# Check service health
echo "📊 Service Health Check:"
kubectl get pods -n naytife | grep -E "(backend|oathkeeper|hydra)"
echo

# Check OAuth client details
echo "🔐 OAuth Client Configuration:"
kubectl exec -n naytife-auth deployment/hydra -- hydra get oauth2-client d39beaaa-9c53-48e7-b82a-37ff52127473 --endpoint http://localhost:4445
echo

# Check backend logs for OAuth-related errors
echo "🔙 Recent Backend Logs:"
kubectl logs -n naytife deployment/backend --tail=20 | grep -i oauth || echo "No OAuth-related logs found"
echo

# Check auth-handler logs
echo "🔑 Recent Auth Handler Logs:"
kubectl logs -n naytife-auth deployment/auth-handler --tail=10
echo

# Test OAuth authorization endpoint
echo "🔓 Testing OAuth Authorization:"
RESPONSE=$(curl -s -I "http://127.0.0.1:8080/oauth2/auth?client_id=d39beaaa-9c53-48e7-b82a-37ff52127473&response_type=code&scope=openid%20offline&redirect_uri=http://127.0.0.1:8080/v1/docs/oauth2-redirect.html&state=test&app_type=dashboard")
echo "$RESPONSE" | head -1

echo
echo "🌐 Access Points:"
echo "  Swagger UI: http://127.0.0.1:8080/v1/docs/index.html"
echo "  OAuth Auth: http://127.0.0.1:8080/oauth2/auth"
echo "  OAuth Token: http://127.0.0.1:8080/oauth2/token"

echo
echo "✅ If everything above looks good, the OAuth flow should work in Swagger UI!"
