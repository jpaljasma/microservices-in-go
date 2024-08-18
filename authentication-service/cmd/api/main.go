package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jpaljasma/authentication/data"
)

const webPort = "80"

const maxCounts int64 = 10

var counts int64 = 0

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Printf("Starting authentication service on port %s ...\n", webPort)

	// connect to database
	conn := connectToDB()

	if conn == nil {
		log.Panic("Cannot connect to database")
	}

	// configuring
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// app.DB = connectToDB()

	// define http server with routes
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start authentication server
	log.Panic(srv.ListenAndServe())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// test DB connection
	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err == nil {
			log.Println("Successfully connected to PostrgeSQL database!")
			return connection
		} else {
			log.Println("PostgreSQL server not ready yet ...")
			counts++
		}
		if counts > maxCounts {
			log.Println(err)
			return nil
		}

		// apply jitter
		jitter, _ := rand.Int(rand.Reader, big.NewInt(333))
		sleepDelay := 300*counts + jitter.Int64()

		log.Printf("Sleeping %d milliseconds between retries ...", sleepDelay)
		// Sleep 100 ms, increasing timeout linearly
		time.Sleep(time.Duration(sleepDelay) * time.Millisecond)
		continue
	}
}
