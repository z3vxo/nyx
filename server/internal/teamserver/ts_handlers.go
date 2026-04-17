package teamserver

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/z3vxo/kronos/internal/database"
)

func SendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := ErrorResponse{Error: message}
	json.NewEncoder(w).Encode(resp)
}

func ts_AgentListHandler(w http.ResponseWriter, r *http.Request) {

	data, err := database.Db_ListAgents()
	if err != nil {
		SendJSONError(w, "Failed retreiving agents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func ts_AgentResolveHandler(w http.ResponseWriter, r *http.Request) {
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

func ts_CommandNewHandler(w http.ResponseWriter, r *http.Request) {
	var cmd TaskEntry
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		SendJSONError(w, "Error decoding json", http.StatusInternalServerError)
		return
	}

	err := database.InsertCommand(cmd.Cmd_type, cmd.Guid, cmd.TaskID, cmd.Param1, cmd.Param2)
	if err != nil {
		SendJSONError(w, "failed inserting command", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
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
