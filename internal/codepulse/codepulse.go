package codepulse

import (
	"log"
	"net/http"
)

// filled by ldflags during build
var Version string

func StartApplication(Version string) {
	codepulse := NewCodepulse()
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", codepulse))
}
