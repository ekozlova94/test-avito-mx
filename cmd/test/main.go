package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"go.uber.org/zap"
	"test-avito-merchant-experience/internal/app/test/getter/prodgetter"
	"test-avito-merchant-experience/internal/app/test/server"
	"test-avito-merchant-experience/internal/app/test/store/sqlstore"
	"test-avito-merchant-experience/internal/app/test/tasks"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var port, dbPath string

func init() {
	flag.StringVar(&port, "port", "server.port", "server port")
	flag.StringVar(&dbPath, "db-path", "database.sql", "path to db")
}

func main() {
	flag.Parse()

	if err := Start(); err != nil {
		log.Fatalf("Error while running application: %s", err)
	}
}

func Start() error {
	logger, _ := zap.NewDevelopment()
	//noinspection GoUnhandledErrorResult
	defer logger.Sync()

	db, err := newDB(dbPath)
	if err != nil {
		return err
	}
	//noinspection GoUnhandledErrorResult
	defer db.Close()

	store := sqlstore.New(db)
	g := prodgetter.NewGetter(logger)
	t := tasks.NewTasks(store, logger, g)
	t.StartBackgroundTask()

	srv := server.NewServer(store, logger, t, g)
	logger.Info("Listening and serving", zap.String("port", port))
	return srv.Router.Run(":" + port)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := runMigrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

func runMigrate(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
