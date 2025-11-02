#!/bin/sh
set -e

echo "🔍 Checking database connectivity..."
echo "DATABASE_URL: ${DATABASE_URL}"

# Set maximum retry attempts (3 minutes with 2 second intervals = 90 attempts)
MAX_ATTEMPTS=90
ATTEMPT=0

# Function to test database connection
test_db_connection() {
    local url="$1"
    local label="$2"
    echo "🔍 Testing $label..."
    if atlas migrate status --dir file://migrations --url "$url" --revisions-schema atlas_schema_revisions 2>/tmp/atlas_error.log; then
        return 0
    else
        echo "❌ Connection failed. Atlas error:"
        cat /tmp/atlas_error.log
        return 1
    fi
}

# Prepare alternative URLs
DB_URL_WITH_SSLMODE_DISABLE=""
if echo "${DATABASE_URL}" | grep -q "sslmode="; then
    # URL already has sslmode parameter
    DB_URL_WITH_SSLMODE_DISABLE=$(echo "${DATABASE_URL}" | sed 's/sslmode=[^&]*/sslmode=disable/')
else
    # Add sslmode=disable to the URL
    if echo "${DATABASE_URL}" | grep -q "?"; then
        DB_URL_WITH_SSLMODE_DISABLE="${DATABASE_URL}&sslmode=disable"
    else
        DB_URL_WITH_SSLMODE_DISABLE="${DATABASE_URL}?sslmode=disable"
    fi
fi

# Wait for database to be ready with timeout
echo "⏳ Attempting to connect to database (timeout: 3 minutes)..."
DB_CONNECTION_URL=""

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    ATTEMPT=$((ATTEMPT + 1))
    echo "⏳ Connection attempt $ATTEMPT/$MAX_ATTEMPTS..."
    
    # Try original DATABASE_URL first
    if test_db_connection "${DATABASE_URL}" "original URL" >/dev/null 2>&1; then
        echo "✅ Database is ready after $ATTEMPT attempts!"
        DB_CONNECTION_URL="${DATABASE_URL}"
        break
    fi
    
    # Try with sslmode=disable every few attempts
    if [ $((ATTEMPT % 5)) -eq 0 ]; then
        if test_db_connection "${DB_URL_WITH_SSLMODE_DISABLE}" "URL with sslmode=disable" >/dev/null 2>&1; then
            echo "✅ Database is ready after $ATTEMPT attempts (using sslmode=disable)!"
            DB_CONNECTION_URL="${DB_URL_WITH_SSLMODE_DISABLE}"
            break
        fi
    fi
    
    if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
        echo "❌ ERROR: Database connection failed after $MAX_ATTEMPTS attempts (3 minutes timeout)"
        echo ""
        echo "🔍 Final detailed connection tests..."
        echo ""
        echo "Testing original URL:"
        test_db_connection "${DATABASE_URL}" "original URL" || true
        echo ""
        echo "Testing with sslmode=disable:"
        test_db_connection "${DB_URL_WITH_SSLMODE_DISABLE}" "sslmode=disable URL" || true
        echo ""
        echo "❌ Deployment failed due to database connectivity issues."
        echo "💡 Common fixes:"
        echo "   1. Add '?sslmode=disable' to your DATABASE_URL in Dokploy"
        echo "   2. Check if database container is running: 'docker ps'"
        echo "   3. Verify network connectivity between containers"
        echo "   4. Check database logs for connection errors"
        exit 1
    fi
    
    sleep 2
done

echo "🚀 Running database migrations with URL: ${DB_CONNECTION_URL}"
if ! atlas migrate apply \
    --dir file://migrations \
    --url "${DB_CONNECTION_URL}" \
    --revisions-schema atlas_schema_revisions \
    --tx-mode none \
    --allow-dirty; then
    echo "❌ ERROR: Database migration failed!"
    echo "❌ Please check migration files and database state"
    exit 1
fi

echo "✅ Migrations completed successfully!"

echo "🚀 Starting backend service..."
exec /app/bin/api
