package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"

	"github.com/spf13/viper"
)

const usageText = `This program runs command on the db. Supported commands are:
  - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
  - create_schema - creates initial set of tables from models (structs).

Usage:
  go run *.go <command> [args]
`

type Config struct {
	PostgresUrl string `mapstructure:"POSTGRES_URL"`
}

func main() {

	var config Config
	viper.SetConfigFile("../.env")
	viper.ReadInConfig()

	err := viper.Unmarshal(&config)

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Set up db connection
	opt, err := pg.ParseURL(config.PostgresUrl)
	if err != nil {
		panic(err)
	}
	db := pg.Connect(opt)

	defer db.Close()

	// 2. Read flags (up, down, version)
	flag.Parse()

	if flag.Arg(0) == "usage" {
		usage()
		os.Exit(2)
	}

	// if flag.Arg(0) == "create_schema" {
	// 	createSchema(db)
	// 	os.Exit(2)
	// }

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Migration failed: %s\n", err)
		os.Exit(1)
	}

	if newVersion != oldVersion {
		fmt.Printf("Migrated from %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("Version is %d\n", oldVersion)
	}

}

func usage() {
	fmt.Print(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}

// func createSchema(db *pg.DB) {
// 	models := []interface{}{
// 		(*repository.Track)(nil),
// 	}

// 	opt := &orm.CreateTableOptions{
// 		IfNotExists: true,
// 	}

// 	for _, model := range models {
// 		err := db.Model(model).CreateTable(opt)

// 		if err != nil {
// 			panic(err)
// 		}

// 		fmt.Printf("Created table successfully: %s", model)
// 	}

// 	fmt.Printf("Created schema successfully")
// }
