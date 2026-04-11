package bootstrap

import (
	"log/slog"

	"github.com/go-pg/pg/v10"
)

// GetConnection returns our pg database connection
// usage:
// db := config.GetConnection()
// defer db.Close()
func GetDbConnection(c *Config, l *slog.Logger) *pg.DB {
	opt, err := pg.ParseURL(c.PostgresUrl)
	if err != nil {
		panic(err)
	}
	db := pg.Connect(opt)

	l.Info("Connected to postgres successfully.")

	return db
}
