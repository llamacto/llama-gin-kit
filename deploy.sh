#!/bin/bash

# Deployment script for llama-gin-kit
set -e

echo "🚀 Starting deployment of llama-gin-kit..."

# Build the Docker image
echo "📦 Building Docker image..."
docker build -t llama-gin-kit:latest .

# Stop existing containers
echo "🛑 Stopping existing containers..."
docker-compose -f docker-compose.prod.yml down || true

# Start the production environment
echo "🔄 Starting production environment..."
docker-compose -f docker-compose.prod.yml up -d

# Wait for container to be ready
echo "⏳ Waiting for container to be ready..."
sleep 10

# Health check
echo "🏥 Performing health check..."
HEALTH_CHECK=$(curl -s http://localhost:8088/v1/health/status || echo "failed")

if [[ $HEALTH_CHECK == *"ok"* ]]; then
    echo "✅ Deployment successful! Server is running at http://localhost:8088"
    echo "📊 Health status: $HEALTH_CHECK"
    echo ""
    echo "📋 Available endpoints:"
    echo "  - Health: http://localhost:8088/v1/health/status"
    echo "  - Register: POST http://localhost:8088/v1/register"
    echo "  - Login: POST http://localhost:8088/v1/login"
    echo "  - Organizations: http://localhost:8088/v1/organizations"
    echo "  - Teams: http://localhost:8088/v1/teams"
else
    echo "❌ Deployment failed! Health check returned: $HEALTH_CHECK"
    echo "📜 Container logs:"
    docker logs llama-gin-kit-prod
    exit 1
fi
