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
	Host         string        `env:"HOST" envDefault:"0.0.0.0"`
	Port         int           `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" envDefault:"120s"`
	MaxFileSize  int           `env:"MAX_FILE_SIZE" envDefault:"10000000"`
}

type DatabaseConfig struct {
	Host            string        `env:"HOST" envDefault:"localhost"`
	Port            int           `env:"PORT" envDefault:"5432"`
	User            string        `env:"USER" envDefault:"postgres"`
	Password        string        `env:"PASSWORD"`
	DBName          string        `env:"NAME"`
	SSLMode         string        `env:"SSL_MODE" envDefault:"disable"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" envDefault:"5m"`
}

func (dbc *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbc.Host, dbc.Port, dbc.User, dbc.Password, dbc.DBName, dbc.SSLMode,
	)
}

type RedisConfig struct {
	Host         string        `env:"HOST" envDefault:"localhost"`
	Port         int           `env:"PORT" envDefault:"6379"`
	Password     string        `env:"PASSWORD"`
	DB           int           `env:"DB" envDefault:"0"`
	PoolSize     int           `env:"POOL_SIZE" envDefault:"10"`
	MinIdleConns int           `env:"MIN_IDLE_CONNS" envDefault:"5"`
	DialTimeout  time.Duration `env:"DIAL_TIMEOUT" envDefault:"5s"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"3s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"3s"`
	MaxRetries   int           `env:"MAX_RETRIES" envDefault:"3"`
}

func (rdc *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", rdc.Host, rdc.Port)
}

type CORSConfig struct {
	AllowedOrigins   string `env:"ALLOWED_ORIGINS" envDefault:"*"`
	AllowedMethods   string `env:"ALLOWED_METHODS" envDefault:"GET,POST,PUT,DELETE,PATCH,OPTIONS"`
	AllowedHeaders   string `env:"ALLOWED_HEADERS" envDefault:"Origin,Content-Type,Accept,Authorization,X-User-Id,X-User-Role"`
	AllowCredentials bool   `env:"ALLOW_CREDENTIALS" envDefault:"true"`
	MaxAge           int    `env:"MAX_AGE" envDefault:"86400"`
}

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
	Enabled           bool `env:"ENABLED" envDefault:"true"`
	RequestsPerMinute int  `env:"REQUESTS_PER_MINUTE" envDefault:"60"`
	BurstSize         int  `env:"BURST_SIZE" envDefault:"10"`
}

type LoggerConfig struct {
	Level      string `env:"LEVEL" envDefault:"info"`
	Format     string `env:"FORMAT" envDefault:"json"`
	OutputPath string `env:"OUTPUT_PATH" envDefault:"stdout"`
}
