#!/bin/bash

# Script to generate environment-specific Oathkeeper access rules
# Usage: ./generate-access-rules.sh <environment> <base_url> <backend_url> <hydra_url> <auth_handler_url> <store_deployer_url>

set -e

ENVIRONMENT=$1
BASE_URL=$2
BACKEND_URL=$3
HYDRA_URL=$4
AUTH_HANDLER_URL=$5
STORE_DEPLOYER_URL=$6

if [[ -z "$ENVIRONMENT" || -z "$BASE_URL" || -z "$BACKEND_URL" || -z "$HYDRA_URL" || -z "$AUTH_HANDLER_URL" || -z "$STORE_DEPLOYER_URL" ]]; then
    echo "Usage: $0 <environment> <base_url> <backend_url> <hydra_url> <auth_handler_url> <store_deployer_url>"
    echo ""
    echo "Example:"
    echo "  $0 local http://127.0.0.1:8080 http://local-backend.naytife.svc.cluster.local:8000 http://local-hydra-public.naytife-auth.svc.cluster.local:4444 http://local-auth-handler.naytife-auth.svc.cluster.local:3000 http://local-store-deployer.naytife.svc.cluster.local:9003"
    echo "  $0 staging https://api-staging.naytife.com http://backend.naytife.svc.cluster.local:8000 http://hydra-public.naytife-auth.svc.cluster.local:4444 http://auth-handler.naytife-auth.svc.cluster.local:3000 http://store-deployer.naytife.svc.cluster.local:9003"
    echo "  $0 production https://api.naytife.com http://backend.naytife.svc.cluster.local:8000 http://hydra-public.naytife-auth.svc.cluster.local:4444 http://auth-handler.naytife-auth.svc.cluster.local:3000 http://store-deployer.naytife.svc.cluster.local:9003"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OVERLAY_DIR="$SCRIPT_DIR/../overlays/$ENVIRONMENT"

# Create overlay directory if it doesn't exist
mkdir -p "$OVERLAY_DIR"

# Generate access rules using template substitution
cat > "$OVERLAY_DIR/oathkeeper-access-rules-$ENVIRONMENT-patch.yaml" << EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: oathkeeper-access-rules
  namespace: naytife-auth
data:
  access-rules.json: |
    [
      {
        "id":"api:rest-api-docs",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/docs",
          "methods": ["GET"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"}, 
        "mutators": [{"handler": "noop"}]
      },
      {
        "id":"api:rest-api-docs-all",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/docs/<**>",
          "methods": ["GET"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"}, 
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "api:user-register",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/auth/register",
          "methods": ["POST"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "api:user-login",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/auth/login",
          "methods": ["POST"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "api:user-refresh",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/auth/refresh",
          "methods": ["POST"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "api:user-logout",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/auth/logout",
          "methods": ["POST"]
        },
        "authenticators": [{"handler": "oauth2_introspection"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "header"}]
      },
      {
        "id": "api:user-profile",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/auth/profile",
          "methods": ["GET", "PUT"]
        },
        "authenticators": [{"handler": "oauth2_introspection"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "header"}]
      },
      {
        "id": "api:cors",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/<**>",
          "methods": ["OPTIONS"]
        },
        "authenticators": [{"handler": "noop"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{
          "handler": "header",
          "config": {
            "headers": {
              "Access-Control-Allow-Origin": "*",
              "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
              "Access-Control-Allow-Headers": "Authorization, Content-Type"
            }
          }
        }]
      },
      {
        "id": "api:customer-register",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/auth/register-customer",
          "methods": ["POST"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "api:customerinfo-public",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/customerinfo<**>",
          "methods": ["GET"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "api:webhooks-public",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/webhooks/<**>",
          "methods": ["POST"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "api:health-endpoints",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/<{health,ready}>",
          "methods": ["GET"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "api:shop-order",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/shops/<*>/orders",
          "methods": ["POST", "GET"]
        },
        "authenticators": [
          { "handler": "anonymous" },
          { "handler": "oauth2_introspection" }
        ],
        "authorizer": { "handler": "allow" },
        "mutators": [{ "handler": "noop" }]
      },
      {
        "id": "api:shop-orders",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/shops/<*>/orders/<**>",
          "methods": ["DELETE", "GET", "PATCH", "PUT"]
        },
        "authenticators": [{"handler": "oauth2_introspection"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{
          "handler": "header",
          "config": {
            "headers": {
              "X-User-Id": "{{ print .Subject }}",
              "Access-Control-Allow-Origin": "*",
              "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
              "Access-Control-Allow-Headers": "Authorization, Content-Type"
            }
          }
        }]
      },
      {
        "id": "oauth:auth-endpoint",
        "upstream": {
          "preserve_host": true,
          "url": "$HYDRA_URL"
        },
        "match": {
          "url": "$BASE_URL/oauth2/auth",
          "methods": ["GET"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "oauth:token-endpoint",
        "upstream": {
          "preserve_host": true,
          "url": "$HYDRA_URL"
        },
        "match": {
          "url": "$BASE_URL/oauth2/token",
          "methods": ["POST", "OPTIONS"]
        },
        "authenticators": [{"handler": "noop"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "oauth:userinfo-endpoint",
        "upstream": {
          "preserve_host": true,
          "url": "$HYDRA_URL"
        },
        "match": {
          "url": "$BASE_URL/userinfo",
          "methods": ["GET", "POST"]
        },
        "authenticators": [{"handler": "oauth2_introspection"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "oidc:well-known",
        "upstream": {
          "preserve_host": true,
          "url": "$HYDRA_URL"
        },
        "match": {
          "url": "$BASE_URL/.well-known/openid-configuration",
          "methods": ["GET"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "oauth:login-endpoint",
        "upstream": {
          "preserve_host": true,
          "url": "$AUTH_HANDLER_URL"
        },
        "match": {
          "url": "$BASE_URL/login",
          "methods": ["GET", "POST"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "oauth:consent-endpoint",
        "upstream": {
          "preserve_host": true,
          "url": "$AUTH_HANDLER_URL"
        },
        "match": {
          "url": "$BASE_URL/consent",
          "methods": ["GET", "POST"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "oauth:callback-endpoint",
        "upstream": {
          "preserve_host": true,
          "url": "$AUTH_HANDLER_URL"
        },
        "match": {
          "url": "$BASE_URL/callback",
          "methods": ["GET"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "oauth:logout-endpoint-auth",
        "upstream": {
          "preserve_host": true,
          "url": "$AUTH_HANDLER_URL"
        },
        "match": {
          "url": "$BASE_URL/logout",
          "methods": ["GET", "POST"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "oauth:error-endpoint",
        "upstream": {
          "preserve_host": true,
          "url": "$AUTH_HANDLER_URL"
        },
        "match": {
          "url": "$BASE_URL/error",
          "methods": ["GET"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "oauth:fallbacks-error-endpoint",
        "upstream": {
          "preserve_host": true,
          "url": "$HYDRA_URL"
        },
        "match": {
          "url": "$BASE_URL/oauth2/fallbacks/error",
          "methods": ["GET"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "oauth:logout-endpoint",
        "upstream": {
          "preserve_host": true,
          "url": "$HYDRA_URL"
        },
        "match": {
          "url": "$BASE_URL/oauth2/sessions/logout",
          "methods": ["GET"]
        },
        "authenticators": [{"handler": "oauth2_introspection"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "api:shops-post-get",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/shops",
          "methods": ["GET", "POST"]
        },
        "authenticators": [{"handler": "oauth2_introspection"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{
          "handler": "header",
          "config": {
            "headers": {
              "X-User-Id": "{{ print .Subject }}",
              "Access-Control-Allow-Origin": "*",
              "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
              "Access-Control-Allow-Headers": "Authorization, Content-Type"
            }
          }
        }]
      },
      {
        "id": "api:shop-crud",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/shops/<*>",
          "methods": ["GET", "PUT", "DELETE"]
        },
        "authenticators": [{"handler": "oauth2_introspection"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{
          "handler": "header",
          "config": {
            "headers": {
              "X-User-Id": "{{ print .Subject }}",
              "Access-Control-Allow-Origin": "*",
              "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
              "Access-Control-Allow-Headers": "Authorization, Content-Type"
            }
          }
        }]
      },
      {
        "id": "api:admin-shop-apis",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/shops/<*>/<{customers,inventory,payment-methods,products,product-types,attributes,deploy,redeploy,deployment-status,update-data,images}><**>",
          "methods": ["GET", "POST", "PUT", "DELETE", "PATCH"]
        },
        "authenticators": [{"handler": "oauth2_introspection"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{
          "handler": "header",
          "config": {
            "headers": {
              "X-User-Id": "{{ print .Subject }}",
              "Access-Control-Allow-Origin": "*",
              "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
              "Access-Control-Allow-Headers": "Authorization, Content-Type"
            }
          }
        }]
      },
      {
        "id": "api:admin-protected-endpoints",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/<{me,userinfo,cart,checkout,images,predefined-product-types,subdomains}><**>",
          "methods": ["GET", "POST", "PUT", "DELETE", "PATCH"]
        },
        "authenticators": [{"handler": "oauth2_introspection"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{
          "handler": "header",
          "config": {
            "headers": {
              "X-User-Id": "{{ print .Subject }}",
              "Access-Control-Allow-Origin": "*",
              "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
              "Access-Control-Allow-Headers": "Authorization, Content-Type"
            }
          }
        }]
      },
      {
        "id": "api:public-graph-playground",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "http://<*>.localhost:8080/graph",
          "methods": ["GET"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "api:public-graph-query",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "http://<*>.localhost:8080/query",
          "methods": ["GET", "POST", "OPTIONS"]
        },
        "authenticators": [
          {"handler": "anonymous"},
          {"handler": "jwt"}
        ],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "build:service",
        "upstream": {
          "preserve_host": true,
          "url": "$STORE_DEPLOYER_URL"
        },
        "match": {
          "url": "$BASE_URL/build<**>",
          "methods": ["GET", "POST"]
        },
        "authenticators": [{"handler": "oauth2_introspection"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "header"}]
      },
      {
        "id": "backend:templates",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/templates<**>",
          "methods": ["GET", "POST", "PUT", "DELETE"]
        },
        "authenticators": [{"handler": "oauth2_introspection"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{
          "handler": "header",
          "config": {
            "headers": {
              "X-User-Id": "{{ print .Subject }}",
              "X-User-Email": "{{ print .Extra.email }}",
              "X-User-Name": "{{ print .Extra.name }}",
              "Access-Control-Allow-Origin": "*",
              "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
              "Access-Control-Allow-Headers": "Authorization, Content-Type"
            }
          }
        }]
      },
      {
        "id": "backend:deployment",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/shops/{shop_id}/deploy<**>",
          "methods": ["GET", "POST", "PUT", "DELETE"]
        },
        "authenticators": [{"handler": "oauth2_introspection"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{
          "handler": "header",
          "config": {
            "headers": {
              "X-User-Id": "{{ print .Subject }}",
              "X-User-Email": "{{ print .Extra.email }}",
              "X-User-Name": "{{ print .Extra.name }}",
              "Access-Control-Allow-Origin": "*",
              "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
              "Access-Control-Allow-Headers": "Authorization, Content-Type"
            }
          }
        }]
      },
      {
        "id": "backend:deployment-status",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/shops/{shop_id}/deployment-status<**>",
          "methods": ["GET"]
        },
        "authenticators": [{"handler": "oauth2_introspection"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{
          "handler": "header",
          "config": {
            "headers": {
              "X-User-Id": "{{ print .Subject }}",
              "X-User-Email": "{{ print .Extra.email }}",
              "X-User-Name": "{{ print .Extra.name }}",
              "Access-Control-Allow-Origin": "*",
              "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
              "Access-Control-Allow-Headers": "Authorization, Content-Type"
            }
          }
        }]
      },
      {
        "id": "backend:services-health",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/health/services",
          "methods": ["GET"]
        },
        "authenticators": [{"handler": "anonymous"}],
        "authorizer": {"handler": "allow"},
        "mutators": [{"handler": "noop"}]
      },
      {
        "id": "api:payments",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/payments<**>",
          "methods": ["GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"]
        },
        "authenticators": [{ "handler": "anonymous" }],
        "authorizer": { "handler": "allow" },
        "mutators": [{ "handler": "noop" }]
      },
      {
        "id": "api:shop-analytics",
        "upstream": {
          "preserve_host": true,
          "url": "$BACKEND_URL"
        },
        "match": {
          "url": "$BASE_URL/v1/shops/<*>/analytics/<**>",
          "methods": ["GET"]
        },
        "authenticators": [{ "handler": "oauth2_introspection" }],
        "authorizer": { "handler": "allow" },
        "mutators": [{
          "handler": "header",
          "config": {
            "headers": {
              "X-User-Id": "{{ print .Subject }}",
              "Access-Control-Allow-Origin": "*",
              "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
              "Access-Control-Allow-Headers": "Authorization, Content-Type"
            }
          }
        }]
      }
    ]
EOF

echo "Generated access rules for $ENVIRONMENT environment: $OVERLAY_DIR/oathkeeper-access-rules-$ENVIRONMENT-patch.yaml"
echo ""
echo "Key URLs configured:"
echo "  Base URL: $BASE_URL"
echo "  Backend URL: $BACKEND_URL"
echo "  Hydra URL: $HYDRA_URL"
echo "  Auth Handler URL: $AUTH_HANDLER_URL"
echo "  Store Deployer URL: $STORE_DEPLOYER_URL"
echo ""
echo "Make sure to update your kustomization.yaml to include this patch file."
