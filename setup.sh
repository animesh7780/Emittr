#!/bin/bash

# Setup script for 4 in a Row game
# This script helps with local development setup

set -e

echo "🎮 Setting up 4 in a Row Game..."

# Check prerequisites
echo "📋 Checking prerequisites..."

if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

echo "✅ Docker and Docker Compose found"

# Create environment files if they don't exist
echo "📄 Setting up environment files..."

if [ ! -f "backend/.env" ]; then
    cp backend/.env.example backend/.env
    echo "✅ Created backend/.env"
else
    echo "⚠️  backend/.env already exists"
fi

if [ ! -f "analytics/.env" ]; then
    cp analytics/.env.example analytics/.env
    echo "✅ Created analytics/.env"
else
    echo "⚠️  analytics/.env already exists"
fi

# Build and start services
echo ""
echo "🐳 Building and starting Docker services..."
echo "This may take a few minutes on first run..."

docker-compose up --build

echo ""
echo "✅ Setup complete!"
echo ""
echo "🎯 Access the application:"
echo "   Frontend: http://localhost:3000"
echo "   Backend API: http://localhost:8080"
echo "   API Docs: http://localhost:8080/health"
echo ""
echo "To stop the services, run: docker-compose down"
