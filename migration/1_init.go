package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table tracks...")

		_, err := db.Exec(`CREATE TABLE tracks (
			id serial PRIMARY KEY,
			name text NOT NULL,
			author text NOT NULL,
			genres text[],
			release_date timestamptz DEFAULT now(),
			created_at timestamptz DEFAULT now()
		)`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table tracks...")
		_, err := db.Exec("DROP TABLE tracks")
		return err
	})
}
