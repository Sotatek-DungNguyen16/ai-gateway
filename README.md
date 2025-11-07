# AI Gateway for Smart Code Review

A production-ready Go service that provides AI-powered code review capabilities. This gateway acts as a unified interface to multiple AI providers (Google Gemini, OpenAI, Anthropic Claude) and formats their responses for integration with the [Smart Code Review GitHub Action](https://github.com/Sotatek-DungNguyen16/ai-review-gateway).

## 🌟 Features

- **Multiple AI Providers**: Support for Google Gemini, OpenAI GPT, and Anthropic Claude
- **Unified API**: Single endpoint that works with multiple AI backends
- **Structured Output**: Returns reviewdog-compatible diagnostic format
- **Secure**: API key authentication for request validation
- **Production Ready**: Docker support, health checks, comprehensive logging
- **Configurable**: Flexible model and provider selection per request
- **Efficient**: Context-aware prompts optimized for code review

## 🏗️ Architecture

```
┌─────────────────────┐
│  GitHub Action      │
│  (client)           │
└──────────┬──────────┘
           │ HTTP POST
           │ /review
           ▼
┌─────────────────────┐
│  AI Gateway         │
│  (this service)     │
│  - Authentication   │
│  - Request parsing  │
│  - Provider routing │
└──────────┬──────────┘
           │
    ┌──────┴──────┬──────────┐
    ▼             ▼          ▼
┌────────┐  ┌─────────┐  ┌─────────┐
│ Gemini │  │ OpenAI  │  │ Claude  │
└────────┘  └─────────┘  └─────────┘
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- At least one AI provider API key:
  - [Google AI Studio](https://makersuite.google.com/app/apikey) for Gemini
  - [OpenAI Platform](https://platform.openai.com/api-keys) for GPT
  - [Anthropic Console](https://console.anthropic.com/) for Claude
- Docker (optional, for containerized deployment)

### Local Development

1. **Clone and navigate to the directory:**

```bash
cd ai-gateway
```

2. **Copy environment configuration:**

```bash
cp env.example .env
```

3. **Edit `.env` with your API keys:**

```bash
# Required: Add at least one API provider key
API_KEYS=your-gateway-secret-key
GOOGLE_API_KEY=your-google-api-key
# OPENAI_API_KEY=your-openai-key  # Optional
# ANTHROPIC_API_KEY=your-claude-key  # Optional
```

4. **Install dependencies:**

```bash
go mod download
```

5. **Run the service:**

```bash
go run .
# Or use make
make run
```

The service will start on `http://localhost:8080`

### Docker Deployment

1. **Copy environment file:**

```bash
cp env.example .env
```

2. **Edit `.env` with your configuration**

3. **Build and run with Docker Compose:**

```bash
docker-compose up -d
```

4. **Check logs:**

```bash
docker-compose logs -f
```

5. **Stop the service:**

```bash
docker-compose down
```

## 📡 API Reference

### Health Check

```bash
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "ai-gateway"
}
```

### Code Review

```bash
POST /review
Content-Type: multipart/form-data
X-API-Key: your-api-key

Fields:
- metadata: JSON string with review configuration
- git_diff: File containing the git diff
```

**Metadata JSON Structure:**

```json
{
  "ai_model": "gemini-2.0-flash",
  "ai_provider": "google",
  "language": "typescript",
  "review_mode": "file",
  "git_info": {
    "commit_hash": "abc123",
    "branch_name": "feature/new-feature",
    "pr_number": "42",
    "repo_url": "https://github.com/user/repo",
    "author": {
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

**Response Format:**

```json
{
  "source": {
    "name": "ai-review",
    "url": ""
  },
  "overview": "Brief summary of the review",
  "diagnostics": [
    {
      "message": "Issue description",
      "location": {
        "path": "src/file.ts",
        "range": {
          "start": {"line": 10, "column": 5},
          "end": {"line": 10, "column": 15}
        }
      },
      "severity": "ERROR",
      "code": {
        "value": "security",
        "url": ""
      },
      "suggestion": "How to fix the issue"
    }
  ]
}
```

### Example Request

```bash
# Create test metadata
cat > metadata.json << 'EOF'
{
  "ai_model": "gemini-2.0-flash",
  "ai_provider": "google",
  "language": "javascript",
  "review_mode": "file"
}
EOF

# Create test diff
cat > test.diff << 'EOF'
diff --git a/test.js b/test.js
+function unsafeQuery(input) {
+  return db.query("SELECT * FROM users WHERE id = " + input);
+}
EOF

# Send request
curl -X POST http://localhost:8080/review \
  -H "X-API-Key: your-api-key" \
  -F "metadata=$(cat metadata.json)" \
  -F "git_diff=@test.diff"
```

## ⚙️ Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | No | `8080` | Server port |
| `API_KEYS` | **Yes** | - | Comma-separated list of valid API keys |
| `GOOGLE_API_KEY` | No* | - | Google Gemini API key |
| `OPENAI_API_KEY` | No* | - | OpenAI API key |
| `ANTHROPIC_API_KEY` | No* | - | Anthropic Claude API key |
| `DEFAULT_AI_PROVIDER` | No | `google` | Default AI provider |
| `DEFAULT_AI_MODEL` | No | `gemini-2.0-flash` | Default AI model |

\* At least one AI provider API key is required

### Supported Models

#### Google Gemini
- `gemini-2.0-flash` (recommended, fastest)
- `gemini-1.5-pro`
- `gemini-1.5-flash`
- `gemini-pro`

#### OpenAI
- `gpt-4o` (recommended)
- `gpt-4-turbo`
- `gpt-4`
- `gpt-3.5-turbo`

#### Anthropic Claude
- `claude-3-5-sonnet-20241022` (recommended)
- `claude-3-opus-20240229`
- `claude-3-sonnet-20240229`
- `claude-3-haiku-20240307`

## 🔒 Security

### API Key Management

1. **Generate Strong Keys:**
```bash
# Generate a secure random key
openssl rand -base64 32
```

2. **Use Multiple Keys:**
```bash
API_KEYS=key-for-ci,key-for-dev,key-for-prod
```

3. **Rotate Keys Regularly:**
   - Add new key to `API_KEYS`
   - Update clients
   - Remove old key after migration

### Best Practices

- ✅ Use HTTPS in production (reverse proxy with SSL)
- ✅ Store API keys in secrets management (GitHub Secrets, AWS Secrets Manager, etc.)
- ✅ Use different keys for different environments
- ✅ Monitor API usage and set up alerts
- ✅ Run behind a firewall or VPN for internal use

## 🚢 Production Deployment

### Using Docker

```bash
# Build image
docker build -t ai-gateway:latest .

# Run container
docker run -d \
  -p 8080:8080 \
  -e API_KEYS="your-secret-key" \
  -e GOOGLE_API_KEY="your-google-key" \
  --name ai-gateway \
  --restart unless-stopped \
  ai-gateway:latest
```

### Using systemd

1. **Build binary:**
```bash
make build
```

2. **Create systemd service** (`/etc/systemd/system/ai-gateway.service`):

```ini
[Unit]
Description=AI Gateway Service
After=network.target

[Service]
Type=simple
User=ai-gateway
WorkingDirectory=/opt/ai-gateway
EnvironmentFile=/opt/ai-gateway/.env
ExecStart=/opt/ai-gateway/ai-gateway
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

3. **Enable and start:**
```bash
sudo systemctl enable ai-gateway
sudo systemctl start ai-gateway
sudo systemctl status ai-gateway
```

### Reverse Proxy (Nginx)

```nginx
server {
    listen 443 ssl http2;
    server_name ai-gateway.yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Increase timeouts for AI processing
        proxy_read_timeout 180s;
        proxy_connect_timeout 180s;
    }
}
```

## 🔧 Development

### Project Structure

```
ai-gateway/
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   │   └── config.go
│   ├── models/            # Data structures
│   │   └── models.go
│   ├── handlers/          # HTTP handlers
│   │   └── review.go
│   ├── middleware/        # HTTP middleware
│   │   └── middleware.go
│   ├── providers/         # AI provider implementations
│   │   ├── provider.go    # Interface and registry
│   │   ├── gemini.go      # Google Gemini
│   │   ├── openai.go      # OpenAI GPT
│   │   └── claude.go      # Anthropic Claude
│   └── prompt/            # Prompt engineering
│       └── prompt.go
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

### Building

```bash
# Build binary
make build

# Run tests
make test

# Clean build artifacts
make clean

# Install dependencies
make deps
```

### Adding a New AI Provider

1. Create a new file in `internal/providers/`
2. Implement the `AIProvider` interface:
```go
type AIProvider interface {
    Review(ctx context.Context, request *models.ReviewRequest) (*models.AIProviderResponse, error)
    Name() string
    SupportedModels() []string
}
```
3. Register the provider in `main.go`

## 🧪 Testing

### Manual Testing

```bash
# Test health endpoint
curl http://localhost:8080/health

# Test review endpoint
bash test-review.sh  # See example above
```

### Integration with GitHub Action

Update your GitHub workflow to use your deployed gateway:

```yaml
- name: Run Smart Code Review
  uses: hiiamtrong/smart-code-review@v1
  with:
    github_token: ${{ secrets.GITHUB_TOKEN }}
    ai_gateway_url: https://your-gateway.yourdomain.com
    ai_gateway_api_key: ${{ secrets.AI_GATEWAY_API_KEY }}
    ai_model: gemini-2.0-flash
    ai_provider: google
```

## 📊 Monitoring

### Health Checks

```bash
# Simple health check
curl http://localhost:8080/health

# With monitoring (Prometheus format can be added)
# See /metrics endpoint (TODO: implement)
```

### Logs

The service logs all requests with:
- HTTP method and path
- Status code
- Response time
- Client IP
- Error details

```bash
# View logs
docker-compose logs -f

# Or with journalctl
sudo journalctl -u ai-gateway -f
```

## 🤝 Integration with Smart Code Review Action

This gateway is designed to work seamlessly with the Smart Code Review GitHub Action. The action:

1. Collects git diff from PR or commits
2. Sends request to this gateway via `/review` endpoint
3. Receives structured diagnostics
4. Posts inline comments on GitHub PR using reviewdog

## 🐛 Troubleshooting

### Common Issues

**1. "No AI providers configured"**
- Ensure at least one `*_API_KEY` environment variable is set
- Check API key validity with the provider

**2. "Invalid API key"**
- Verify `X-API-Key` header matches a key in `API_KEYS`
- Check for whitespace or formatting issues in env file

**3. "Provider not available"**
- Check that the requested provider's API key is configured
- Verify the provider name (google, openai, anthropic)

**4. "Context deadline exceeded"**
- Large diffs may timeout (default 120s)
- Consider splitting large reviews
- Check AI provider API status

**5. Empty or invalid responses**
- Check AI provider API quotas
- Verify network connectivity to AI provider
- Review logs for detailed error messages

## 📝 License

MIT License - see the main repository for details.

## 🙏 Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📮 Support

- **Issues**: [GitHub Issues](https://github.com/Sotatek-DungNguyen16/ai-review-gateway/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Sotatek-DungNguyen16/ai-review-gateway/discussions)

---

Built with ❤️ for better code reviews

