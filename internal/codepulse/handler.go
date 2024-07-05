package codepulse

import (
	"encoding/json"
	"log"
	"net/http"
)

func (app *Codepulse) welcomeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := app.db.Exec(`CREATE TABLE IF NOT EXISTS reqs (name TEXT);`)
	if err != nil {
		log.Printf("failed to create table: %v", err)
	} else {
		log.Print("table create success")
	}

	_, _ = app.db.Exec("INSERT INTO reqs (name) VALUES (?)", r.RemoteAddr)

	response := map[string]interface{}{
		"app":              "CodePulse",
		"remoteAddress":    r.RemoteAddr,
		"computedResponse": computedResponse(),
	}
	json.NewEncoder(w).Encode(response)
}
