package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/z3vxo/kronos/internal/config"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := ErrorResponse{Error: message}
	json.NewEncoder(w).Encode(resp)
}

func Send404(w http.ResponseWriter) {
	w.Header().Set("Server", "kronos")
	w.Header().Set("Content-Type", "text/html")

	path := config.Cfg.Server.NotFoundFile
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = home + path[1:]
	}

	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "<h1>404 not found</h1>", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(content)
}
