package turso

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// Config controls the libSQL connection. Zero pool settings use database/sql defaults.
type Config struct {
	URL             string
	Token           string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// ConnectFromEnv reads TURSO_DATABASE_URL and TURSO_AUTH_TOKEN.
func ConnectFromEnv(ctx context.Context) (*sqlx.DB, error) {
	return Connect(ctx, Config{URL: os.Getenv("TURSO_DATABASE_URL"), Token: os.Getenv("TURSO_AUTH_TOKEN")})
}

// Connect opens a Turso connection and verifies it with PingContext.
func Connect(ctx context.Context, cfg Config) (*sqlx.DB, error) {
	if strings.TrimSpace(cfg.URL) == "" || strings.TrimSpace(cfg.Token) == "" {
		return nil, errors.New("ymirr/turso: URL and token are required")
	}
	if cfg.MaxOpenConns < 0 || cfg.MaxIdleConns < 0 || cfg.ConnMaxLifetime < 0 || cfg.ConnMaxIdleTime < 0 {
		return nil, errors.New("ymirr/turso: pool settings cannot be negative")
	}
	if cfg.MaxOpenConns > 0 && cfg.MaxIdleConns > cfg.MaxOpenConns {
		return nil, errors.New("ymirr/turso: MaxIdleConns cannot exceed MaxOpenConns")
	}
	u, err := url.Parse(cfg.URL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, errors.New("ymirr/turso: invalid database URL")
	}
	if u.Scheme != "libsql" && u.Scheme != "https" && u.Scheme != "wss" {
		return nil, errors.New("ymirr/turso: unsupported database URL scheme")
	}
	q := u.Query()
	q.Set("authToken", cfg.Token)
	u.RawQuery = q.Encode()
	db, err := sqlx.Open("libsql", u.String())
	if err != nil {
		return nil, fmt.Errorf("ymirr/turso: open: %w", err)
	}
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ymirr/turso: ping failed: %w", err)
	}
	return db, nil
}
