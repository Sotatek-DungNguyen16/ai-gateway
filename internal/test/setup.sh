#!/bin/bash
# Setup script for AI Gateway

set -e

echo "🚀 Setting up AI Gateway..."

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    echo "Visit: https://golang.org/doc/install"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✓ Found Go version: $GO_VERSION"

# Check for env file
if [ ! -f .env ]; then
    if [ -f env.example ]; then
        echo "📝 Creating .env file from env.example..."
        cp env.example .env
        echo "⚠️  Please edit .env file and add your API keys"
    else
        echo "⚠️  No .env file found. Please create one with your configuration."
    fi
else
    echo "✓ .env file exists"
fi

# Download dependencies
echo "📦 Downloading Go dependencies..."
if go mod download; then
    echo "✓ Dependencies downloaded"
else
    echo "⚠️  Some dependencies may need manual intervention"
fi

# Tidy up
echo "🧹 Tidying up dependencies..."
go mod tidy || echo "⚠️  go mod tidy had warnings (may be ok)"

# Try to build
echo "🔨 Building application..."
if go build -o ai-gateway .; then
    echo "✓ Build successful!"
    echo ""
    echo "🎉 Setup complete!"
    echo ""
    echo "Next steps:"
    echo "1. Edit .env file and add your API keys"
    echo "2. Run: ./ai-gateway"
    echo "   Or: make run"
    echo "   Or: docker-compose up"
else
    echo "❌ Build failed. Please check the errors above."
    exit 1
fi

