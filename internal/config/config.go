package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig      `json:"app"`
	Server   ServerConfig   `json:"server"`
	Print    PrintConfig    `json:"print"`
	Laravel  LaravelConfig  `json:"laravel"`
	Security SecurityConfig `json:"security"`
	Logger   LoggerConfig   `json:"logger"`
}

type AppConfig struct {
	Name        string `json:"name" env:"APP_NAME" default:"PixelPrintService"`
	Environment string `json:"environment" env:"APP_ENV" default:"local"`
	Debug       bool   `json:"debug" env:"APP_DEBUG" default:"true"`
	Version     string `json:"version" env:"APP_VERSION" default:"1.0.0"`
}

type ServerConfig struct {
	Host         string        `json:"host" env:"SERVER_HOST" default:"0.0.0.0"`
	Port         int           `json:"port" env:"SERVER_PORT" default:"8080"`
	ReadTimeout  time.Duration `json:"read_timeout" env:"SERVER_READ_TIMEOUT" default:"30s"`
	WriteTimeout time.Duration `json:"write_timeout" env:"SERVER_WRITE_TIMEOUT" default:"30s"`
	IdleTimeout  time.Duration `json:"idle_timeout" env:"SERVER_IDLE_TIMEOUT" default:"60s"`
}

type PrintConfig struct {
	Command              string        `json:"command" env:"PRINT_COMMAND" default:"Start-Process -FilePath \"%s\" -Verb Print"`
	Timeout              time.Duration `json:"timeout" env:"PRINT_TIMEOUT" default:"60s"`
	AllowedFormats       []string      `json:"allowed_formats" env:"PRINT_ALLOWED_FORMATS" default:"txt,pdf,png,jpg"`
	MaxFileSize          int64         `json:"max_file_size" env:"PRINT_MAX_FILE_SIZE" default:"10485760"` // 10MB
	TempDir              string        `json:"temp_dir" env:"PRINT_TEMP_DIR" default:"./temp"`
	QueueSize            int           `json:"queue_size" env:"PRINT_QUEUE_SIZE" default:"100"`
	KeepArchive          bool          `json:"keep_archive" env:"PRINT_KEEP_ARCHIVE" default:"false"`
	ArchiveRetentionDays int           `json:"archive_retention_days" env:"PRINT_ARCHIVE_RETENTION_DAYS" default:"90"`
	ArchiveDir           string        `json:"archive_dir" env:"PRINT_ARCHIVE_DIR" default:"./archive"`
	CleanupInterval      int           `json:"cleanup_interval" env:"PRINT_CLEANUP_INTERVAL" default:"60"`
}

type LaravelConfig struct {
	APIKey     string        `json:"api_key" env:"LARAVEL_API_KEY" required:"true"`
	BaseURL    string        `json:"base_url" env:"LARAVEL_BASE_URL" default:"http://127.0.0.1"`
	Timeout    time.Duration `json:"timeout" env:"LARAVEL_TIMEOUT" default:"30s"`
	RetryCount int           `json:"retry_count" env:"LARAVEL_RETRY_COUNT" default:"3"`
}

type SecurityConfig struct {
	EnableCORS     bool     `json:"enable_cors" env:"SECURITY_ENABLE_CORS" default:"true"`
	AllowedIPs     []string `json:"allowed_ips" env:"SECURITY_ALLOWED_IPS"`
	RateLimit      int      `json:"rate_limit" env:"SECURITY_RATE_LIMIT" default:"100"` // requetes par minute
	TrustedProxies []string `json:"trusted_proxies" env:"SECURITY_TRUSTED_PROXIES"`
}

type LoggerConfig struct {
	Level      string `json:"level" env:"LOGGING_LEVEL" default:"info"`
	FilePath   string `json:"file_path" env:"LOGGING_FILE_PATH" default:"logs/app.log"`
	MaxSize    int    `json:"max_size" env:"LOGGING_MAX_SIZE" default:"10"`
	MaxBackups int    `json:"max_backups" env:"LOGGING_MAX_BACKUPS" default:"5"`
	MaxAge     int    `json:"max_age" env:"LOGGING_MAX_AGE" default:"30"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("erreur lors du chargement du fichier .env : %v", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("erreur lors du parsing de la configuration : %v", err)
	}

	return cfg, nil
}
