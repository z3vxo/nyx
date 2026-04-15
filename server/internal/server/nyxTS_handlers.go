package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/z3vxo/nyx/internal/database"
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

func nyx_AgentHandler(w http.ResponseWriter, r *http.Request) {

	data, err := database.Db_ListAgents()
	if err != nil {
		SendJSONError(w, "Failed retreiving agents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func nyx_AgentResolveHandler(w http.ResponseWriter, r *http.Request) {
	codeName := chi.URLParam(r, "codename")
	if codeName == "" {
		SendJSONError(w, "missing codename", http.StatusBadRequest)
		return
	}

	AgentGuid, err := database.ResolveCodename(codeName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			SendJSONError(w, "Agent Codename not found", http.StatusNotFound)
			return
		}
		SendJSONError(w, "database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"guid": AgentGuid})
}

type TaskEntry struct {
	Cmd_type int    `json:"type"`
	Guid     string `json:"guid"`
	Param1   string `json:"param1"`
	Param2   string `json:"param2"`
}

func nyx_CommandNewHandler(w http.ResponseWriter, r *http.Request) {
	var cmd TaskEntry
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		SendJSONError(w, "Error decoding json", http.StatusInternalServerError)
		return
	}

	err := database.InsertCommand(cmd.Cmd_type, cmd.Guid, cmd.Param1, cmd.Param2)
	if err != nil {
		SendJSONError(w, "failed inserting command", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	return
}

type UserLogin struct {
	Username string `json:"user"`
	Password string `json:"passwd"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var log UserLogin
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		SendJSONError(w, "failed decoding json", http.StatusInternalServerError)
		return
	}

	if !CheckLogin(log.Username, log.Password) {
		SendJSONError(w, "invalid login", http.StatusUnauthorized)
		return
	}

	token, err := CraftJWT(log.Username)
	if err != nil {
		SendJSONError(w, "Failed Crafting jwt", http.StatusInternalServerError)

		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
