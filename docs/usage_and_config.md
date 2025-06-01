# Project Usage & Configuration Guide

This document provides a step-by-step guide for configuring and using the Llama Gin Kit project. It is intended for international open-source users and is written in English for global accessibility.

## 1. Overview

Llama Gin Kit is a modern Go scaffold for AI-powered development, built with Cursor and Windsurf conventions. It supports JWT authentication, database integration, Redis, email, OpenAI, and agent-based architecture out of the box.

## 2. Configuration Philosophy

- **Environment variables** are the primary way to configure sensitive and environment-specific values.
- **YAML config files** are used for non-sensitive, structured configuration.
- **Never commit real secrets or credentials** to the repository. Use placeholders in `.env.example` and `config.example.yaml`.

## 3. Configuration Workflow

### Step 1: Clone the Repository

```bash
git clone https://github.com/llamacto/llama-gin-kit.git
cd llama-gin-kit
```

### Step 2: Install Dependencies

```bash
go mod download
```

### Step 3: Prepare Environment Variables

- Copy the example file:
  ```bash
  cp .env.example .env
  ```
- Fill in your actual values for all required variables in `.env` (database, JWT secret, LLM API keys, email credentials, etc).

### Step 4: Prepare YAML Config (Optional)

- If your project uses a YAML config (e.g., `config/config.yaml`), copy the example:
  ```bash
  cp config/config.example.yaml config/config.yaml
  ```
- Edit `config.yaml` as needed for your environment (non-sensitive settings only).

### Step 5: Database Migration

```bash
make migrate
```

### Step 6: Start the Application

```bash
make run
```

The server will start on the port specified in your environment variables or config file (default: 6066).

## 4. Configuration Precedence

- **Production**: Only system environment variables are loaded. `.env` is ignored for security.
- **Development**: `.env` is loaded automatically if present, allowing easy local overrides.
- **YAML config**: Used for non-sensitive, structured settings. Never store secrets here.

## 5. Best Practices

- Always keep `.env.example` and `config.example.yaml` up to date with all required variables and settings.
- Never commit `.env` or real config files with secrets to version control.
- Use strong, unique secrets for JWT, database, and third-party services.
- Document any new configuration options in both the code and the example files.
- Configure proper rate limiting for LLM API endpoints to manage costs.

## 6. Example Environment Variables

```
# Database
DB_USERNAME=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=llama_gin_kit

# Security
JWT_SECRET=your_jwt_secret

# LLM APIs
OPENAI_API_KEY=your_openai_key
ANTHROPIC_API_KEY=your_anthropic_key

# Email
EMAIL_FROM="Llama Gin Kit <noreply@llamacto.com>"
EMAIL_RESEND_API_KEY=your_resend_key
```

## 7. Further Reading

- [Gin Documentation](https://gin-gonic.com/docs/)
- [Twelve-Factor App Methodology](https://12factor.net/config)
- [OpenAI API Documentation](https://platform.openai.com/docs)

## 8. Troubleshooting

- If the application fails to start, check for missing or incorrect environment variables.
- Ensure your database and Redis services are running and accessible.
- Verify your LLM API keys are valid and have sufficient quota.
- Review logs for detailed error messages.

---

For more details, see the main [README.md](../README.md) or open an issue on GitHub. 
