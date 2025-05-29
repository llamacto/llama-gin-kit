# ZGI-GinKit

ZGI-GinKit is an enterprise-level web development scaffold based on the Gin framework, integrating common enterprise features and best practices.

## Features

- 📦 Modular architecture
- 🔐 JWT authentication
- 📝 Swagger API documentation
- 🚦 Rate limiting
- 📨 Asynchronous task queue
- 🔄 WebSocket support
- 📊 GORM database operations
- 💾 Redis cache
- 📧 Email service
- 🔍 Unified error handling
- 📝 Structured logging (Zap)
- ⚙️ Configuration management (Viper)

## Quick Start

### Requirements

- Go 1.21+
- PostgreSQL 12+
- Redis 6.0+

### Installation

```bash
git clone https://github.com/zgiai/zgi-ginkit.git
cd zgi-ginkit
go mod download
```

### Configuration

1. Copy the environment variable template:
```bash
cp .env.example .env
```
2. Fill in the required variables in `.env` (all sensitive information should use placeholders for open source security).

3. Copy the config file template (if available):
```bash
cp config/config.example.yaml config/config.yaml
```

### Run

```bash
# Run database migration
make migrate

# Start the service
make run
```

## Project Structure

```
your-gin-project/
├── cmd/                   # Entry files (server, migrate, tools, etc.)
├── app/                   # Business modules (e.g., user)
├── config/                # Configuration management
├── middleware/            # Gin middleware
├── pkg/                   # Utility packages
├── routes/                # Route management
├── storage/               # Static/persistent resources
├── docs/                  # API documentation
└── ...
```

## Development Guide

### Add a New Module

1. Create a new module directory under `app`
2. Implement model, repository, service, and handler
3. Register the new module's routes in `routes`

### Run Tests

```bash
make test
```

### Generate API Documentation

```bash
make swagger
```

## Environment Variables

All sensitive information, secrets, and credentials are configured via the `.env` file. Do not commit real secrets to the repository; only commit `.env.example`.

Common environment variable examples:
```
DB_USERNAME=<your-db-username>
DB_PASSWORD=<your-db-password>
JWT_SECRET=<your-jwt-secret>
OPENAI_API_KEY=<your-openai-api-key>
...
```

## Deployment

### Docker

```bash
# Build the image
docker build -t zgi-ginkit .

# Run the container
docker run -p 8080:8080 zgi-ginkit
```

## Contribution Guide

Pull requests and issues are welcome!

## License

MIT License 
