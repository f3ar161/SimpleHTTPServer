package internal

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type app struct {
	DB *sql.DB
}

// ConnectionDB: postgres DB connection
func ConnectionDB() *app {

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Failed to get DB url")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB conn %v", err)
	}

	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // close db conn
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping DB %v", err)
	}

	application := &app{DB: db}

	log.Println("Connected to the DB....")
	return application
}
