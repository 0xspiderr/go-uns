package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/0xspiderr/go-uns/internal/models"
)

var dbConn *sql.DB

func ConfigureDB() {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	ssl := os.Getenv("DB_SSLMODE")
	dbName := os.Getenv("DB_NAME")
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, dbName, ssl)
	var err error
	dbConn, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Couldn't open db: %v", err)
	}

	// retry database connection every 5 seconds and hang program until the connection is established
	for {
		if err = dbConn.Ping(); err != nil {
			log.Println("Couldn't establish connection to the DB, retrying...")
		}
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
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
	_, err := dbConn.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}
}

func InsertUNSData(data models.UNSData) error {
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

	_, err = dbConn.Exec(query, time.Now(), data.Topic, otJSON, itJSON)
	if err != nil {
		return fmt.Errorf("Failed to insert UNS data: %w", err)
	}

	fmt.Println("Inserted UNS data in the table successfully")
	return nil
}

func Close() {
	if dbConn != nil {
		dbConn.Close()
	}
}
