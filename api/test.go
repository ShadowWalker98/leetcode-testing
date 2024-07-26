package main

import (
	"Testing/internal/data"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type Application struct {
	cfg     config
	models  data.Models
	logger  *log.Logger
	version string
}

func main() {
	var cfg config

	cfg.db.dsn = "postgres://problemsuser:password@localhost/problems?sslmode=disable"

	conn, err := openDB(cfg)

	if err != nil {
		fmt.Println(err)
	}

	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := Application{
		cfg:     cfg,
		models:  data.NewModels(conn),
		logger:  logger,
		version: "1.0",
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 3000),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(time.Minute * 15)
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type database struct {
	dsn string
}

type config struct {
	db database
}
