#!/bin/sh
set -e

echo "🔍 Checking database connectivity..."

# Wait for database to be ready
until atlas migrate status --dir file://migrations --url "${DATABASE_URL}" --revisions-schema atlas_schema_revisions >/dev/null 2>&1; do
    echo "⏳ Waiting for database to be ready..."
    sleep 2
done

echo "✅ Database is ready!"

echo "🚀 Running database migrations..."
atlas migrate apply \
    --dir file://migrations \
    --url "${DATABASE_URL}" \
    --revisions-schema atlas_schema_revisions \
    --tx-mode none \
    --allow-dirty

echo "✅ Migrations completed successfully!"

echo "🚀 Starting backend service..."
exec /app/bin/api
