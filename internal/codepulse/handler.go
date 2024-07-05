package codepulse

import (
	"encoding/json"
	"net/http"
)

func (s *Codepulse) welcomeHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"app":              "CodePulse",
		"remoteAddress":    r.RemoteAddr,
		"computedResponse": computedResponse(),
	}
	json.NewEncoder(w).Encode(response)
}
