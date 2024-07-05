package codepulse

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// filled by ldflags during build
var Version string

func connectDB() *sql.DB {
	db, err := sql.Open("sqlite3", "data/codepulseDB.db")
	if err != nil {
		// Handle error
		log.Printf("failed to open DB conn")
		os.Exit(1)
	}
	log.Printf("DB connected!\n")
	return db
}

func StartApplication(Version string) {
	log.Println("starting application codepulse version", Version)
	codepulse := NewCodepulse()
	codepulse.db = connectDB()
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", codepulse))
}
