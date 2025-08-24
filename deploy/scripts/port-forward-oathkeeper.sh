#!/bin/bash

# Port forward Oathkeeper proxy to localhost:8080 for local development
echo "Setting up port forwarding for Oathkeeper..."
echo "Access Oathkeeper at: http://localhost:8080"
echo "Press Ctrl+C to stop port forwarding"

kubectl port-forward -n naytife-auth service/oathkeeper-proxy 8080:8080
