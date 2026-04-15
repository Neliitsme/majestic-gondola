package bootstrap

import (
	"context"
	"log/slog"

	"github.com/go-pg/pg/v10"
)

// GetConnection returns our pg database connection
// usage:
// db := config.GetConnection()
// defer db.Close()
func NewDbConnection(c *Config, l *slog.Logger) *pg.DB {
	opt, err := pg.ParseURL(c.PostgresUrl)
	if err != nil {
		panic(err)
	}
	db := pg.Connect(opt)

	ctx := context.Background()

	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	if l != nil {
		l.Info("Connected to postgres successfully.")
	}

	return db
}
