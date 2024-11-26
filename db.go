package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "Ubuntu"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "tag_db"
)

func connectToDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("connectoToDB1", err)
		return nil, errors.New(err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Println("connectoToDB2", err)
		return nil, errors.New(err.Error())
	}

	return db, nil
}
