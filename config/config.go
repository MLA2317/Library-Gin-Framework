package config

import (
    "fmt"
    "os"
    "strconv"
    "github.com/joho/godotenv"
)

type Config struct {
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
    Server   ServerConfig
    Upload   UploadConfig
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type RedisConfig struct {
    Host     string
    Port     string
    Password string
    DB       int
}

type JWTConfig struct {
    Secret          string
    ExpirationHours int
}

type ServerConfig struct {
    Port string
}

type UploadConfig struct {
    Path        string
    MaxFileSize int64
}

func Load() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        return nil, fmt.Errorf("error loading .env file: %w", err)
    }

    redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
    jwtExp, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
    maxFileSize, _ := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "10485760"), 10, 64)

    return &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "laziz2317"),
            DBName:   getEnv("DB_NAME", "bookgolang"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnv("REDIS_PORT", "6379"),
            Password: getEnv("REDIS_PASSWORD", ""),
            DB:       redisDB,
        },
        JWT: JWTConfig{
            Secret:          getEnv("JWT_SECRET", "secret"),
            ExpirationHours: jwtExp,
        },
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
        },
        Upload: UploadConfig{
            Path:        getEnv("UPLOAD_PATH", "./uploads"),
            MaxFileSize: maxFileSize,
        },
    }, nil
}

func getEnv(key, defaultVal string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultVal
}

func (c *DatabaseConfig) DSN() string {
    return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}