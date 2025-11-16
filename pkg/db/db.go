package db

import (
	"database/sql"
	"os"
	"strconv"
	"strings"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/knadh/koanf/v2"

	"pr-reviewer-assign-service/pkg/app"
)

const (
	appName             = "pr-reviewer-assign-service"
	defaultPostgresPort = 5432
)

func Init(_ *app.App, cfg *koanf.Koanf) (*sql.DB, error) {
	connConfig, err := pgx.ParseConfig(buildDSN(cfg))
	if err != nil {
		return nil, err
	}

	connConfig.Tracer = otelpgx.NewTracer()

	db := stdlib.OpenDB(*connConfig)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	return db, nil
}

func buildDSN(cfg *koanf.Koanf) string {
	dsn := &strings.Builder{}

	dsn.WriteString("host=")
	dsn.WriteString(cfg.MustString("host"))

	port := defaultPostgresPort

	if cfg.Int("port") > 0 {
		port = cfg.Int("port")
	}

	dsn.WriteString(" port=")
	dsn.WriteString(strconv.Itoa(port))

	dsn.WriteString(" dbname=")
	dsn.WriteString(cfg.MustString("database"))

	searchPath := cfg.String("search_path")
	if cfg.String("search_path") != "" {
		dsn.WriteString(" search_path=")
		dsn.WriteString(searchPath)
	}

	dsn.WriteString(" application_name=")
	dsn.WriteString(appName)

	user := cfg.String("user")
	if userEnv, found := os.LookupEnv("DATABASE_USER"); found {
		user = userEnv
	}

	password := cfg.String("password")
	if passwordEnv, found := os.LookupEnv("DATABASE_PASSWORD"); found {
		password = passwordEnv
	}

	dsn.WriteString(" user=")
	dsn.WriteString(user)
	dsn.WriteString(" password=")
	dsn.WriteString(password)

	sslMode := "disable"

	if cfg.Exists("sslmode") {
		sslMode = cfg.String("sslmode")
	}

	dsn.WriteString(" sslmode=")
	dsn.WriteString(sslMode)

	dsnStr := dsn.String()

	return dsnStr
}
