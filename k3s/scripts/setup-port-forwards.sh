#!/bin/bash

echo "🚀 Setting up port forwards for Naytife services..."

# Kill any existing port forwards
pkill -f "kubectl.*port-forward" || true

# Wait a moment for processes to die
sleep 2

# Start port forwards in background
echo "  📡 Setting up API Gateway (Oathkeeper) on port 8080..."
kubectl port-forward svc/oathkeeper-proxy -n naytife-auth 8080:8080 &

echo "  🔙 Setting up Backend API on port 8000..."
kubectl port-forward svc/backend -n naytife 8000:8000 &

echo "  🔑 Setting up Auth Handler on port 3000..."
kubectl port-forward svc/auth-handler -n naytife-auth 3000:3000 &

echo "  🏗️  Setting up Template Registry on port 9001..."
kubectl port-forward svc/template-registry -n naytife 9001:9001 &

echo "  🚀 Setting up Store Deployer on port 9003..."
kubectl port-forward svc/store-deployer -n naytife 9003:9003 &

echo "  🐘 Setting up PostgreSQL on port 5432..."
kubectl port-forward svc/postgres -n naytife 5432:5432 &

echo "  📊 Setting up Redis on port 6379..."
kubectl port-forward svc/redis -n naytife 6379:6379 &

echo "  🆔 Setting up Hydra Public on port 4444..."
kubectl port-forward svc/hydra-public -n naytife-auth 4444:4444 &

echo "  🆔 Setting up Hydra Admin on port 4445..."
kubectl port-forward svc/hydra-admin -n naytife-auth 4445:4445 &

# Wait a moment for forwards to establish
sleep 3

echo ""
echo "✅ Port forwards established!"
echo ""
echo "🔗 Service URLs:"
echo "  🔐 API Gateway:      http://localhost:8080"
echo "  🔙 Backend API:      http://localhost:8000"
echo "  🔑 Auth Handler:     http://localhost:3000"
echo "  🏗️  Template Registry: http://localhost:9001"
echo "  🚀 Store Deployer:   http://localhost:9003"
echo "  🐘 PostgreSQL:       localhost:5432"
echo "  📊 Redis:            localhost:6379"
echo "  🆔 Hydra Public:     http://localhost:4444"
echo "  🆔 Hydra Admin:      http://localhost:4445"
echo ""
echo "📝 To stop all port forwards: pkill -f 'kubectl.*port-forward'"
echo ""
echo "🌐 Press Ctrl+C to stop port forwarding..."

# Keep the script running
wait
