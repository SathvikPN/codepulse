package codepulse

import (
	"codepulse/internal/middleware"
	"database/sql"
	"net/http"
)

type Codepulse struct {
	httpMultiplexer *http.ServeMux
	db              *sql.DB
}

func NewCodepulse() *Codepulse {
	mux := http.NewServeMux()
	server := &Codepulse{httpMultiplexer: mux}
	server.routes()
	return server
}

// URL Mappings
func (app *Codepulse) routes() {
	app.httpMultiplexer.Handle("/welcome", middleware.RateLimiting(middleware.Logging(http.HandlerFunc(app.welcomeHandler))))
}

func (app *Codepulse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.httpMultiplexer.ServeHTTP(w, r)
}
