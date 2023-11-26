package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5" // pgx PostgreSQL driver
)

type Seat struct {
	ID      int
	Student string
}

func main() {
	// Database connection string
	postgresUrl := "postgres://faraz:A1s2d3f4.@localhost:5432/Seat"
	db, err := pgx.Connect(context.Background(), os.Getenv(postgresUrl))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())

	// Execute the query
	rows, err := db.Query(context.Background(), "SELECT CASE WHEN MOD(id, 2) = 0 THEN id - 1 WHEN id = (SELECT MAX(id) FROM Seat) AND MOD(id, 2) = 1 THEN id ELSE id + 1 END AS id, student FROM Seat ORDER BY id")
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		var seat Seat
		if err := rows.Scan(&seat.ID, &seat.Student); err != nil {
			log.Fatalf("Query scan failed: %v\n", err)
		}
		fmt.Printf("ID: %d, Student: %s\n", seat.ID, seat.Student)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatalf("Error during rows iteration: %v\n", err)
	}
}
