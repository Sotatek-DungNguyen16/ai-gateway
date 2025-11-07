# Quick Start Guide

## 1. Prerequisites

- **Go 1.21+**: [Install Go](https://golang.org/doc/install)
- **AI Provider API Key**: Get at least one:
  - [Google AI Studio](https://makersuite.google.com/app/apikey) for Gemini (Recommended - Free tier available)
  - [OpenAI Platform](https://platform.openai.com/api-keys) for GPT
  - [Anthropic Console](https://console.anthropic.com/) for Claude

## 2. Setup (Local Development)

```bash
# 1. Navigate to ai-gateway directory
cd ai-gateway

# 2. Run setup script
chmod +x setup.sh
./setup.sh

# 3. Edit .env file with your keys
nano .env  # or use your favorite editor
```

Required in `.env`:
```bash
API_KEYS=your-secret-gateway-key
GOOGLE_API_KEY=your-google-api-key  # Get from Google AI Studio
```

## 3. Run the Service

```bash
# Option 1: Direct run
go run .

# Option 2: Build and run
make build
./ai-gateway

# Option 3: Using make
make run
```

## 4. Test the Service

```bash
# Test health endpoint
curl http://localhost:8080/health

# Run full test suite (requires service to be running)
export AI_GATEWAY_API_KEY=your-secret-gateway-key
./test-review.sh
```

## 5. Use with GitHub Action

Update your GitHub workflow:

```yaml
# .github/workflows/code-review.yml
name: Code Review

on:
  pull_request:
    types: [opened, synchronize]

jobs:
  review:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write
      checks: write
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Smart Code Review
        uses: hiiamtrong/smart-code-review@v1
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          # Point to your deployed gateway
          ai_gateway_url: http://your-server:8080
          ai_gateway_api_key: ${{ secrets.AI_GATEWAY_API_KEY }}
          ai_model: gemini-2.0-flash
          ai_provider: google
```

## 6. Deploy (Docker)

```bash
# 1. Edit .env with your configuration
cp env.example .env
nano .env

# 2. Build and run with Docker Compose
docker-compose up -d

# 3. Check logs
docker-compose logs -f

# 4. Stop
docker-compose down
```

## 7. Deploy (Production Server)

### Using systemd

```bash
# 1. Build binary
make build

# 2. Copy to server
scp ai-gateway user@your-server:/opt/ai-gateway/
scp .env user@your-server:/opt/ai-gateway/

# 3. Create systemd service
sudo nano /etc/systemd/system/ai-gateway.service
```

```ini
[Unit]
Description=AI Gateway Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/ai-gateway
EnvironmentFile=/opt/ai-gateway/.env
ExecStart=/opt/ai-gateway/ai-gateway
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

```bash
# 4. Start service
sudo systemctl enable ai-gateway
sudo systemctl start ai-gateway
sudo systemctl status ai-gateway
```

## Configuration

### Supported AI Models

#### Google Gemini (Recommended)
- `gemini-2.0-flash` ⚡ Fastest, great quality
- `gemini-1.5-pro` 🎯 Best quality
- `gemini-1.5-flash` ⚡ Fast, good quality

#### OpenAI
- `gpt-4o` 🎯 Latest, best
- `gpt-4-turbo` ⚡ Fast GPT-4
- `gpt-3.5-turbo` 💰 Cheapest

#### Anthropic Claude
- `claude-3-5-sonnet-20241022` 🎯 Best quality
- `claude-3-opus-20240229` 🎯 High quality
- `claude-3-haiku-20240307` ⚡ Fast, cheap

## Troubleshooting

### "No AI providers configured"
```bash
# Check your .env file has at least one API key
cat .env | grep API_KEY
```

### "Invalid API key"
```bash
# The X-API-Key header must match one of the keys in API_KEYS
# Test with:
curl -H "X-API-Key: your-key" http://localhost:8080/health
```

### "Build failed"
```bash
# Clean and retry
go clean
go mod download
go mod tidy
go build .
```

### Service not responding
```bash
# Check if running
ps aux | grep ai-gateway

# Check logs
journalctl -u ai-gateway -f  # systemd
docker-compose logs -f       # docker
```

## Next Steps

1. ✅ Setup AI Gateway (you're here!)
2. 📖 Read [full README](README.md) for advanced configuration
3. 🔒 Setup HTTPS with nginx/caddy for production
4. 📊 Add monitoring (logs, metrics)
5. 🚀 Deploy and integrate with GitHub Actions

## Getting Help

- 📚 [Full Documentation](README.md)
- 🐛 [Report Issues](https://github.com/Sotatek-DungNguyen16/ai-review-gateway/issues)
- 💬 [Discussions](https://github.com/Sotatek-DungNguyen16/ai-review-gateway/discussions)

---

**Pro Tips:**
- 🎁 Google Gemini has a generous free tier - great for getting started
- ⚡ Use `gemini-2.0-flash` for fast reviews
- 💰 Monitor your API usage to avoid unexpected costs
- 🔒 Always use HTTPS in production
- 📝 Keep your API keys in environment variables, never in code

