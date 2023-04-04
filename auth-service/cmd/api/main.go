package main

import (
	"auth/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "8081"

var counts int64

type Config struct{
	DB *sql.DB
	Models data.Models
}

func main() {
	log.Println("starting auth service ")

	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect with postgres")
	}

	app := Config{
		DB: conn,
		Models: data.New(conn),
	}

	srv := http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func openDB(dsc string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsc)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgest not yet ready ... ")
			counts ++
		} else {
			log.Println("connected to postgres")
			return connection
		}

		if counts >10 {
			log.Println("err")
			return nil
		}

		log.Println("Backing off for 2 sec...")
		time.Sleep(2*time.Second)
		continue
	}
}