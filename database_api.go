package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func configureDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Couldn't load .env variables: %v", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	ssl := os.Getenv("DB_SSLMODE")
	dbName := os.Getenv("DB_NAME")
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, dbName, ssl)

	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Couldn't open db: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Couldn't connect to the db: %v", err)
	}

	fmt.Println("DB connection successful")

	createUNSTable()
}

func createUNSTable() {
	query := `CREATE TABLE IF NOT EXISTS uns_data (
	id BIGSERIAL PRIMARY KEY,
	event_timestamp TIMESTAMPTZ NOT NULL,
	topic TEXT NOT NULL,
	ot_data JSONB,
	it_data JSONB
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}
}

func insertUNSData(data UNSData) error {
	otJSON, err := json.Marshal(data.OT)
	if err != nil {
		return fmt.Errorf("Failed to marshal OT data: %w", err)
	}

	itJSON, err := json.Marshal(data.IT)
	if err != nil {
		return fmt.Errorf("Failed to marshal IT data: %w", err)
	}

	query := `
	INSERT INTO uns_data (event_timestamp, topic, ot_data, it_data)
	VALUES ($1, $2, $3, $4)
	`

	_, err = db.Exec(query, time.Now(), data.Topic, otJSON, itJSON)
	if err != nil {
		return fmt.Errorf("Failed to insert UNS data: %w", err)
	}

	fmt.Println("Inserted UNS data in the table successfully")
	return nil
}
