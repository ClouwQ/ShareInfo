package config

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Environment string          `env:"ENVIRONMENT" envDefault:"dev"`
	Server      ServerConfig    `envPrefix:"SERVER_"`
	Database    DatabaseConfig  `envPrefix:"DB_"`
	Redis       RedisConfig     `envPrefix:"REDIS_"`
	CORS        CORSConfig      `envPrefix:"CORS_"`
	RateLimit   RateLimitConfig `envPrefix:"RATE_LIMIT_"`
	Logger      LoggerConfig    `envPrefix:"LOG_"`
}

type ServerConfig struct {
	Host         string        `env:"SERVER_HOST" envDefault:"0.0.0.0"`
	Port         int           `env:"SERVER_PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" envDefault:"120s"`
}

type DatabaseConfig struct {
	Host            string        `env:"DB_HOST" envDefault:"0.0.0.0"`
	Port            int           `env:"DB_PORT" envDefault:"5432"`
	User            string        `env:"DB_USER" envDefault:"postgres"`
	Password        string        `env:"DB_PASSWORD"`
	DBName          string        `env:"DBNAME"`
	SSLMode         string        `env:"DB_SSLMODE" envDefault:"disable"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`
}

// DSN — строка на подключение к Postgres
func (dbc *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbc.Host, dbc.Port, dbc.User, dbc.Password, dbc.DBName, dbc.SSLMode,
	)
}

type RedisConfig struct {
	Host         string        `env:"REDIS_HOST" envDefault:"localhost"`
	Port         int           `env:"REDIS_PORT" envDefault:"6379"`
	Password     string        `env:"REDIS_PASSWORD"`
	DB           int           `env:"REDIS_DB" envDefault:"0"`
	PoolSize     int           `env:"REDIS_POOL_SIZE" envDefault:"10"`
	MinIdleConns int           `env:"REDIS_MIN_IDLE_CONNS" envDefault:"5"`
	DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT" envDefault:"5s"`
	ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT" envDefault:"3s"`
	WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT" envDefault:"3s"`
	MaxRetries   int           `env:"REDIS_MAX_RETRIES" envDefault:"3"`
}

// Address — строка для подключения к Redis
func (rdc *RedisConfig) Address() string { return fmt.Sprintf("%s:%d", rdc.Host, rdc.Port) }

type CORSConfig struct {
	AllowedOrigins   string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
	AllowedMethods   string `env:"CORS_ALLOWED_METHODS" envDefault:"GET,POST,PUT,DELETE,PATCH,OPTIONS"`
	AllowedHeaders   string `env:"CORS_ALLOWED_HEADERS" envDefault:"Origin,Content-Type,Accept,Authorization,X-User-Id,X-User-Role"`
	AllowCredentials bool   `env:"CORS_ALLOW_CREDENTIALS" envDefault:"true"`
	MaxAge           int    `env:"CORS_MAX_AGE" envDefault:"86400"`
}

// GetAllowedOrigins возвращает список разрешенных origins
func (c *CORSConfig) GetAllowedOrigins() []string {
	if c.AllowedOrigins == "*" {
		return []string{"*"}
	}
	origins := strings.Split(c.AllowedOrigins, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}
	return origins
}

type RateLimitConfig struct {
	Enabled           bool `env:"RATE_LIMIT_ENABLED" envDefault:"true"`
	RequestsPerMinute int  `env:"RATE_LIMIT_REQUESTS_PER_MINUTE" envDefault:"60"`
	BurstSize         int  `env:"RATE_LIMIT_BURST_SIZE" envDefault:"10"`
}

type LoggerConfig struct {
	Level      string `env:"LOG_LEVEL" envDefault:"info"`
	Format     string `env:"LOG_FORMAT" envDefault:"json"`
	OutputPath string `env:"LOG_OUTPUT_PATH" envDefault:"stdout"`
}
