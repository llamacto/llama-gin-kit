# API Examples

## 项目信息接口

### GET /

获取项目完整信息，包括版本、功能特性、可用接口等。

**请求:**
```bash
curl http://localhost:6066/
```

**响应:**
```json
{
  "success": true,
  "message": "Welcome to Llama Gin Kit API",
  "data": {
    "name": "Llama Gin Kit",
    "description": "A comprehensive REST API service built with Go and Gin framework, featuring user management, organization/team management, API key authentication, and modular architecture.",
    "version": "v1.0.0",
    "go_version": "1.23.0+",
    "build_time": "2025-06-19 19:01:55",
    "environment": "debug",
    "api": {
      "version": "v1",
      "base_url": "/v1",
      "endpoints": [
        "POST /v1/register - User registration",
        "POST /v1/login - User login",
        "GET /v1/users/profile - Get user profile",
        "POST /v1/organizations - Create organization",
        "GET /v1/organizations - List organizations",
        "POST /v1/teams - Create team",
        "GET /v1/teams/:id - Get team details",
        "POST /v1/apikeys - Create API key",
        "GET /v1/apikeys - List API keys"
      ],
      "features": [
        "JWT Authentication",
        "API Key Authentication",
        "User Management",
        "Organization Management",
        "Team Management",
        "Role-based Access Control",
        "Email Notifications",
        "PostgreSQL Database",
        "Docker Support",
        "Swagger Documentation"
      ]
    },
    "links": {
      "documentation": "/swagger/index.html",
      "health": "/v1/health/status",
      "swagger": "/swagger/*any"
    }
  }
}
```

## 健康检查接口

### GET /ping

简单的健康检查接口，用于快速验证服务状态。

**请求:**
```bash
curl http://localhost:6066/ping
```

**响应:**
```json
{
  "message": "pong"
}
```

### GET /v1/health/status

详细的健康状态检查，包含版本信息。

**请求:**
```bash
curl http://localhost:6066/v1/health/status
```

**响应:**
```json
{
  "status": "ok",
  "version": "v1"
}
```

## 用户管理接口

### POST /v1/register

用户注册接口。

**请求:**
```bash
curl -X POST http://localhost:6066/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**响应:**
```json
{
  "id": 1,
  "created_at": "2025-06-19T19:02:16.668183508+08:00",
  "updated_at": "2025-06-19T19:02:16.668183508+08:00",
  "deleted_at": null,
  "username": "testuser",
  "email": "test@example.com",
  "nickname": "",
  "avatar": "",
  "phone": "",
  "bio": "",
  "status": 1,
  "last_login": null
}
```

### POST /v1/login

用户登录接口。

**请求:**
```bash
curl -X POST http://localhost:6066/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

**响应:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "nickname": "",
    "avatar": "",
    "phone": "",
    "bio": "",
    "status": 1
  }
}
```

## Docker 部署示例

### 本地开发环境

```bash
# 构建镜像
docker build -t ginext:latest .

# 运行容器
docker run -d --name ginext \
  -p 6066:6066 \
  -e DB_HOST="your-db-host" \
  -e DB_PORT="5432" \
  -e DB_USERNAME="your-username" \
  -e DB_PASSWORD="your-password" \
  -e DB_NAME="gin-kit" \
  -e JWT_SECRET="your-jwt-secret" \
  ginext:latest

# 测试根路径
curl http://localhost:6066/
```

### 生产环境

```bash
# 运行生产容器
docker run -d --name ginext-prod \
  -p 8080:6066 \
  -e DB_HOST="production-db-host" \
  -e DB_PORT="5432" \
  -e DB_USERNAME="prod-username" \
  -e DB_PASSWORD="prod-password" \
  -e DB_NAME="gin-kit" \
  -e GIN_MODE="release" \
  -e JWT_SECRET="production-jwt-secret" \
  ginext:latest

# 测试项目信息
curl http://localhost:8080/
```

## 错误响应示例

### 400 Bad Request

```json
{
  "error": "Invalid request format",
  "code": 400,
  "message": "The request body is malformed or missing required fields"
}
```

### 401 Unauthorized

```json
{
  "error": "Unauthorized",
  "code": 401,
  "message": "Invalid or missing authentication token"
}
```

### 500 Internal Server Error

```json
{
  "error": "Internal server error",
  "code": 500,
  "message": "An unexpected error occurred"
}
``` 
