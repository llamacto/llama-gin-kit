# Llama Gin Kit

Llama Gin Kit is a modern Go scaffold for AI-powered development, built with Cursor and Windsurf conventions. It supports automated testing, LLM APIs, and agent-based architecture out of the box.

## Features

- 📦 Modular architecture designed for AI coding
- 🤖 Built-in LLM API integrations (OpenAI, Claude, etc.)
- 🔐 JWT authentication with AI-enhanced security
- 📝 Auto-generated Swagger API documentation
- 🚦 Intelligent rate limiting
- 📨 Asynchronous task queue for AI workflows
- 🔄 WebSocket support for real-time AI interactions
- 📊 GORM database operations with AI query optimization
- 💾 Redis cache optimized for LLM responses
- 📧 Smart email service with AI templates
- 🔍 Unified error handling with AI diagnostics
- 📝 Structured logging (Zap) for AI debugging
- ⚙️ Configuration management optimized for AI services
- 🎯 Agent-based architecture support
- 🧠 Ready-to-use AI prompt templates
- 🔄 Streaming response support for LLM interactions

## Quick Start

### Requirements

- Go 1.21+
- PostgreSQL 12+
- Redis 6.0+
- OpenAI API Key (or other LLM provider)

### Installation

```bash
git clone https://github.com/llamacto/llama-gin-kit.git
cd llama-gin-kit
go mod download
```

### Configuration

1. Copy the environment variable template:
```bash
cp .env.example .env
```

2. Configure your LLM API keys and other services in `.env`:
```bash
# LLM Configuration
OPENAI_API_KEY=your_openai_api_key
# Add other LLM providers as needed

# Database
DB_USERNAME=your_db_username
DB_PASSWORD=your_db_password
DB_NAME=llama_gin_kit

# JWT for AI-enhanced auth
JWT_SECRET=your_jwt_secret

# Other services...
```

3. Copy the config file template (if available):
```bash
cp config/config.example.yaml config/config.yaml
```

### Run

```bash
# Run database migration
make migrate

# Start the AI-powered service
make run
```

## Project Structure

```
llama-gin-kit/
├── cmd/                   # Entry files (server, migrate, tools, etc.)
├── app/                   # Business modules (user, ai-agents, etc.)
│   ├── user/             # User management
│   └── agents/           # AI agent implementations
├── config/               # Configuration management
├── middleware/           # Gin middleware (including AI middleware)
├── pkg/                  # Utility packages
│   ├── ai/              # AI service integrations
│   ├── llm/             # LLM client implementations
│   └── agents/          # Agent framework
├── routes/               # Route management
├── storage/              # Static/persistent resources
├── docs/                 # API documentation
└── templates/            # AI prompt templates
```

## AI Features

### LLM Integration

The kit comes with built-in support for multiple LLM providers:

- **OpenAI GPT models** (GPT-4, GPT-3.5-turbo)
- **Streaming responses** for real-time AI interactions
- **Prompt template management** for consistent AI outputs
- **Token usage tracking** and cost optimization

### Agent-Based Architecture

Build sophisticated AI agents with:

- **Multi-step reasoning** workflows
- **Tool integration** for external API calls
- **Memory management** for context retention
- **Parallel processing** for complex tasks

### AI-Enhanced APIs

- **Intelligent text processing** endpoints
- **Automated content generation** services
- **Real-time language translation** with context awareness
- **Smart data analysis** and insights generation

## Development Guide

### Add a New AI Module

1. Create a new module directory under `app/`
2. Implement model, repository, service, and handler with AI integration
3. Add LLM-specific functionality in `pkg/ai/`
4. Register routes in `routes/` with appropriate middleware

### Add Custom LLM Provider

1. Implement the LLM interface in `pkg/llm/`
2. Add configuration in `config/`
3. Register the provider in your service initialization

### Run Tests

```bash
make test
```

### Generate API Documentation

```bash
make swagger
```

## Environment Variables

All sensitive information, secrets, and API keys are configured via the `.env` file. Do not commit real secrets to the repository; only commit `.env.example`.

Critical environment variables for AI features:
```bash
# LLM APIs
OPENAI_API_KEY=<your-openai-api-key>
ANTHROPIC_API_KEY=<your-anthropic-api-key>

# Database
DB_USERNAME=<your-db-username>
DB_PASSWORD=<your-db-password>
DB_NAME=llama_gin_kit

# Security
JWT_SECRET=<your-jwt-secret>

# Redis for caching LLM responses
REDIS_HOST=localhost
REDIS_PASSWORD=<your-redis-password>
```

## Deployment

### Docker

```bash
# Build the AI-powered image
docker build -t llama-gin-kit .

# Run the container with AI services
docker run -p 8080:8080 -e OPENAI_API_KEY=your_key llama-gin-kit
```

### Production Considerations

- Use environment variables for all LLM API keys
- Configure proper rate limiting for AI endpoints
- Set up monitoring for LLM usage and costs
- Implement proper error handling for AI service failures

## AI Coding Optimizations

This scaffold is specifically optimized for AI-assisted development:

- **Cursor IDE integration** with proper .cursorrules
- **Windsurf conventions** for seamless AI coding
- **Automated test generation** templates
- **AI-friendly code structure** for better LLM understanding
- **Built-in prompt engineering** utilities

## Contributing

Pull requests and issues are welcome! This project is designed to evolve with the AI coding ecosystem.

## License

MIT License

---

Built with ❤️ for the AI coding community. Optimized for Cursor, Windsurf, and modern LLM workflows. 
