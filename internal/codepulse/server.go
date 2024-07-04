package codepulse

import (
	"codepulse/internal/middleware"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server {
	mux := http.NewServeMux()
	server := &Server{mux: mux}
	server.routes()
	return server
}

// URL Mappings
func (s *Server) routes() {
	s.mux.Handle("/welcome", middleware.RateLimiting(middleware.Logging(http.HandlerFunc(s.welcomeHandler))))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
