package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// GlobalConfig stores the global configuration
var GlobalConfig *Config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Log      LogConfig
	OpenAI   OpenAIConfig
	R2       R2Config
	Email    EmailConfig
	App      AppConfig
}

type ServerConfig struct {
	Port           int    `json:"port"`
	Mode           string `json:"mode"`
	ReadTimeout    int    `json:"read_timeout"`
	WriteTimeout   int    `json:"write_timeout"`
	MaxHeaderBytes int    `json:"max_header_bytes"`
}

type DatabaseConfig struct {
	Driver          string `json:"driver"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"-"` // 敏感信息不序列化
	DBName          string `json:"dbname"`
	SSLMode         string `json:"sslmode"`
	Timezone        string `json:"timezone"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	MaxOpenConns    int    `json:"max_open_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Password     string `json:"-"` // 敏感信息不序列化
	DB           int    `json:"db"`
	PoolSize     int    `json:"pool_size"`
	MinIdleConns int    `json:"min_idle_conns"`
}

type JWTConfig struct {
	Secret         string        `json:"-"` // 敏感信息不序列化
	ExpireDays     int           `json:"expire_days"`
	ExpireDuration time.Duration `json:"-"`
}

type LogConfig struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
	Compress   bool   `json:"compress"`
}

type OpenAIConfig struct {
	APIKey string `json:"-"` // 敏感信息不序列化
}

type R2Config struct {
	AccessKeyID     string `json:"-"` // 敏感信息不序列化
	SecretAccessKey string `json:"-"` // 敏感信息不序列化
	Bucket          string `json:"bucket"`
	Region          string `json:"region"`
	Endpoint        string `json:"endpoint"`
	PublicURL       string `json:"public_url"`
	PublicDomain    string `json:"public_domain"`
}

type EmailConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"-"` // 敏感信息不序列化
	From         string `json:"from"`
	ResendAPIKey string `json:"-"` // 敏感信息不序列化
}

type AppConfig struct {
	Name      string        `json:"name"`
	Version   string        `json:"version"`
	Secret    string        `json:"-"` // 敏感信息不序列化
	JWTSecret string        `json:"-"` // 敏感信息不序列化
	JWTExpire time.Duration `json:"jwt_expire"`
}

// Load loads configuration from .env file
func Load() (*Config, error) {
	// 仅在开发环境加载 .env 文件
	if os.Getenv("SERVER_MODE") == "" || os.Getenv("SERVER_MODE") == "debug" {
		_ = godotenv.Load()
	}

	config := &Config{}

	// Load server config
	if err := loadServerConfig(config); err != nil {
		return nil, err
	}

	// Load database config
	if err := loadDatabaseConfig(config); err != nil {
		return nil, err
	}

	// Load redis config
	if err := loadRedisConfig(config); err != nil {
		return nil, err
	}

	// Load JWT config
	if err := loadJWTConfig(config); err != nil {
		return nil, err
	}

	// Load log config
	if err := loadLogConfig(config); err != nil {
		return nil, err
	}

	// Load OpenAI config
	if err := loadOpenAIConfig(config); err != nil {
		return nil, err
	}

	// Load R2 config
	if err := loadR2Config(config); err != nil {
		return nil, err
	}

	// Load email config
	if err := loadEmailConfig(config); err != nil {
		return nil, err
	}

	// Load app config
	if err := loadAppConfig(config); err != nil {
		return nil, err
	}

	// Validate config
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	GlobalConfig = config
	return config, nil
}

func loadServerConfig(config *Config) error {
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "6066"))
	if err != nil {
		return fmt.Errorf("invalid SERVER_PORT: %v", err)
	}

	readTimeout, err := strconv.Atoi(getEnv("SERVER_READ_TIMEOUT", "60"))
	if err != nil {
		return fmt.Errorf("invalid SERVER_READ_TIMEOUT: %v", err)
	}

	writeTimeout, err := strconv.Atoi(getEnv("SERVER_WRITE_TIMEOUT", "60"))
	if err != nil {
		return fmt.Errorf("invalid SERVER_WRITE_TIMEOUT: %v", err)
	}

	maxHeaderBytes, err := strconv.Atoi(getEnv("SERVER_MAX_HEADER_BYTES", "1048576"))
	if err != nil {
		return fmt.Errorf("invalid SERVER_MAX_HEADER_BYTES: %v", err)
	}

	config.Server = ServerConfig{
		Port:           port,
		Mode:           getEnv("SERVER_MODE", "debug"),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	return nil
}

func loadDatabaseConfig(config *Config) error {
	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return fmt.Errorf("invalid DB_PORT: %v", err)
	}

	maxIdleConns, err := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
	if err != nil {
		return fmt.Errorf("invalid DB_MAX_IDLE_CONNS: %v", err)
	}

	maxOpenConns, err := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "100"))
	if err != nil {
		return fmt.Errorf("invalid DB_MAX_OPEN_CONNS: %v", err)
	}

	connMaxLifetime, err := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME", "3600"))
	if err != nil {
		return fmt.Errorf("invalid DB_CONN_MAX_LIFETIME: %v", err)
	}

	config.Database = DatabaseConfig{
		Driver:          getEnv("DB_DRIVER", "postgres"),
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            port,
		Username:        getEnv("DB_USERNAME", "postgres"),
		Password:        getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "llama_gin_kit"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		Timezone:        getEnv("DB_TIMEZONE", "Asia/Shanghai"),
		MaxIdleConns:    maxIdleConns,
		MaxOpenConns:    maxOpenConns,
		ConnMaxLifetime: connMaxLifetime,
	}

	return nil
}

func loadRedisConfig(config *Config) error {
	port, err := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	if err != nil {
		return fmt.Errorf("invalid REDIS_PORT: %v", err)
	}

	db, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return fmt.Errorf("invalid REDIS_DB: %v", err)
	}

	poolSize, err := strconv.Atoi(getEnv("REDIS_POOL_SIZE", "10"))
	if err != nil {
		return fmt.Errorf("invalid REDIS_POOL_SIZE: %v", err)
	}

	minIdleConns, err := strconv.Atoi(getEnv("REDIS_MIN_IDLE_CONNS", "5"))
	if err != nil {
		return fmt.Errorf("invalid REDIS_MIN_IDLE_CONNS: %v", err)
	}

	config.Redis = RedisConfig{
		Host:         getEnv("REDIS_HOST", "localhost"),
		Port:         port,
		Password:     getEnv("REDIS_PASSWORD", ""),
		DB:           db,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	}

	return nil
}

func loadJWTConfig(config *Config) error {
	expireDays, err := strconv.Atoi(getEnv("JWT_EXPIRE_DAYS", "7"))
	if err != nil {
		return fmt.Errorf("invalid JWT_EXPIRE_DAYS: %v", err)
	}

	config.JWT = JWTConfig{
		Secret:         getEnv("JWT_SECRET", ""),
		ExpireDays:     expireDays,
		ExpireDuration: time.Duration(expireDays) * 24 * time.Hour,
	}

	return nil
}

func loadLogConfig(config *Config) error {
	maxSize, err := strconv.Atoi(getEnv("LOG_MAX_SIZE", "100"))
	if err != nil {
		return fmt.Errorf("invalid LOG_MAX_SIZE: %v", err)
	}

	maxAge, err := strconv.Atoi(getEnv("LOG_MAX_AGE", "30"))
	if err != nil {
		return fmt.Errorf("invalid LOG_MAX_AGE: %v", err)
	}

	maxBackups, err := strconv.Atoi(getEnv("LOG_MAX_BACKUPS", "7"))
	if err != nil {
		return fmt.Errorf("invalid LOG_MAX_BACKUPS: %v", err)
	}

	compress, err := strconv.ParseBool(getEnv("LOG_COMPRESS", "true"))
	if err != nil {
		return fmt.Errorf("invalid LOG_COMPRESS: %v", err)
	}

	config.Log = LogConfig{
		Level:      getEnv("LOG_LEVEL", "debug"),
		Filename:   getEnv("LOG_FILENAME", "logs/app.log"),
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		Compress:   compress,
	}

	return nil
}

func loadOpenAIConfig(config *Config) error {
	config.OpenAI = OpenAIConfig{
		APIKey: getEnv("OPENAI_API_KEY", ""),
	}
	return nil
}

func loadR2Config(config *Config) error {
	config.R2 = R2Config{
		AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
		SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
		Bucket:          getEnv("R2_BUCKET", ""),
		Region:          getEnv("R2_REGION", "auto"),
		Endpoint:        getEnv("R2_ENDPOINT", ""),
		PublicURL:       getEnv("R2_PUBLIC_URL", ""),
		PublicDomain:    getEnv("R2_PUBLIC_DOMAIN", ""),
	}
	return nil
}

func loadEmailConfig(config *Config) error {
	port, err := strconv.Atoi(getEnv("EMAIL_PORT", "587"))
	if err != nil {
		return fmt.Errorf("invalid EMAIL_PORT: %v", err)
	}

	config.Email = EmailConfig{
		Host:         getEnv("EMAIL_HOST", "smtp.gmail.com"),
		Port:         port,
		Username:     getEnv("EMAIL_USERNAME", ""),
		Password:     getEnv("EMAIL_PASSWORD", ""),
		From:         getEnv("EMAIL_FROM", ""),
		ResendAPIKey: getEnv("EMAIL_RESEND_API_KEY", ""),
	}
	return nil
}

func loadAppConfig(config *Config) error {
	expireDays, err := strconv.Atoi(getEnv("APP_JWT_EXPIRE_DAYS", "7"))
	if err != nil {
		return fmt.Errorf("invalid APP_JWT_EXPIRE_DAYS: %v", err)
	}

	config.App = AppConfig{
		Name:      getEnv("APP_NAME", "Llama-Gin-Kit"),
		Version:   getEnv("APP_VERSION", "1.0.0"),
		Secret:    getEnv("APP_SECRET", ""),
		JWTSecret: getEnv("APP_JWT_SECRET", ""),
		JWTExpire: time.Duration(expireDays) * 24 * time.Hour,
	}
	return nil
}

func validateConfig(config *Config) error {
	// Validate required fields
	if config.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}

	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
