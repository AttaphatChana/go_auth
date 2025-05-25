package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port         string        `json:"port"` // this is called a struct tag
	DBHost       string        `json:"host"`
	RedisHost    string        `json:"redis_host"`
	JWTSecret    string        `json:"jwt_secret"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

func getEnv(key string) (string, error) {
	v, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("environment variable %s is required", key)
	}
	return v, nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseDurationOrDefault(s, def string) (time.Duration, error) {
	if s == "" {
		return time.ParseDuration(def)
	}
	return time.ParseDuration(s)
}

func Load() (*Config, error) {
	var err error
	cfg := &Config{}

	// Required values for DB connection parts:
	dbHost, err := getEnv("DB_HOST")
	if err != nil {
		return nil, err
	}
	dbPort, err := getEnv("DB_PORT")
	if err != nil {
		return nil, err
	}
	dbUser, err := getEnv("DB_USER")
	if err != nil {
		return nil, err
	}
	dbPass, err := getEnv("DB_PASS")
	if err != nil {
		return nil, err
	}
	dbName, err := getEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	// Assemble full Postgres connection string
	cfg.DBHost = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// Other required values:
	if cfg.RedisHost, err = getEnv("REDIS_ADDR"); err != nil {
		return nil, err
	}
	if cfg.JWTSecret, err = getEnv("JWT_SECRET"); err != nil {
		return nil, err
	}

	// Optional with defaults:
	cfg.Port = getenvDefault("PORT", "8000")

	if cfg.ReadTimeout, err = parseDurationOrDefault(os.Getenv("READ_TIMEOUT"), "5s"); err != nil {
		return nil, fmt.Errorf("invalid READ_TIMEOUT: %w", err)
	}
	if cfg.WriteTimeout, err = parseDurationOrDefault(os.Getenv("WRITE_TIMEOUT"), "10s"); err != nil {
		return nil, fmt.Errorf("invalid WRITE_TIMEOUT: %w", err)
	}

	// Sanity checks:
	if _, err := strconv.Atoi(cfg.Port); err != nil {
		return nil, errors.New("PORT must be a valid integer")
	}

	return cfg, nil
}

func connectDB(cfg *Config) (*pgxpool.Pool, error) {
	// Use cfg.DBHost (connection string) to open the pool
	pool, err := pgxpool.New(context.Background(), cfg.DBHost)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func ConnectRedis(cfg *Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisHost,
		// Add other fields like Password or DB index if needed
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}

func RetryRedis(cfg *Config) (*redis.Client, error) {
	rdb, err := ConnectRedis(cfg)
	attempts := 0
	for attempts <= 5 {
		rdb, err = ConnectRedis(cfg)
		if err == nil {
			return rdb, nil
		}
		time.Sleep(10 * time.Second)
		attempts++
	}
	return rdb, err
}

func RetryDB(cfg *Config) (*pgxpool.Pool, error) {
	pool, err := connectDB(cfg)
	attempts := 0
	for attempts <= 5 {
		pool, err = connectDB(cfg)
		if err == nil {
			return pool, nil
		}
		time.Sleep(10 * time.Second)
		attempts++
	}
	return pool, err

}
