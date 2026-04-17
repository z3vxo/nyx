package teamserver

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/z3vxo/kronos/internal/httputil"
)

// @Summary      List all agents
// @Tags         agents
// @Produce      json
// @Success      200  {object}  database.Agents
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /ts/rest/agents/list [get]
func (ts *TeamServer) AgentListHandler(w http.ResponseWriter, r *http.Request) {

	data, err := ts.db.ListAgents()
	if err != nil {
		httputil.SendJSONError(w, "Failed retreiving agents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// @Summary      Resolve agent codename to GUID
// @Tags         agents
// @Produce      json
// @Param        codename  path      string  true  "Agent codename"
// @Success      200       {object}  map[string]string
// @Failure      400       {object}  ErrorResponse
// @Failure      404       {object}  ErrorResponse
// @Failure      500       {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /ts/rest/agents/resolve/{codename} [get]
func (ts *TeamServer) AgentResolveHandler(w http.ResponseWriter, r *http.Request) {
	codeName := chi.URLParam(r, "codename")
	if codeName == "" {
		httputil.SendJSONError(w, "missing codename", http.StatusBadRequest)
		return
	}

	AgentGuid, err := ts.db.ResolveCodename(codeName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.SendJSONError(w, "Agent Codename not found", http.StatusNotFound)
			return
		}
		httputil.SendJSONError(w, "database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"guid": AgentGuid})
}

// @Summary      Queue a new command for an agent
// @Tags         commands
// @Accept       json
// @Produce      json
// @Param        body  body      TaskEntry  true  "Command payload"
// @Success      200   {object}  map[string]string
// @Failure      500   {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /ts/rest/commands/new [post]
func (ts *TeamServer) CommandNewHandler(w http.ResponseWriter, r *http.Request) {
	var cmd TaskEntry
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		httputil.SendJSONError(w, "Error decoding json", http.StatusInternalServerError)
		return
	}

	err := ts.db.InsertCommand(cmd.Cmd_type, cmd.Guid, cmd.TaskID, cmd.Param1, cmd.Param2)
	if err != nil {
		httputil.SendJSONError(w, "failed inserting command", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

// @Summary      Delete a queued command by task ID
// @Tags         commands
// @Accept       json
// @Produce      json
// @Param        body  body      TaskDelete  true  "Task to delete"
// @Success      200   {object}  map[string]string
// @Failure      500   {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /ts/rest/commands/delete [post]
func (ts *TeamServer) CommandDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var task TaskDelete

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		httputil.SendJSONError(w, "Error Decoding json", http.StatusInternalServerError)
		return
	}

	if err := ts.db.DeleteTask(task.TaskID); err != nil {
		httputil.SendJSONError(w, "Failed Deleting task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})

}

// @Summary      Authenticate and receive a JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      UserLogin  true  "Login credentials"
// @Success      200   {object}  map[string]string
// @Failure      401   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /ts/rest/login [post]
func (ts *TeamServer) loginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var log UserLogin
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		httputil.SendJSONError(w, "failed decoding json", http.StatusInternalServerError)
		return
	}

	if !ts.Auth.CheckLogin(log.Username, log.Password) {
		httputil.SendJSONError(w, "invalid login", http.StatusUnauthorized)
		return
	}

	token, err := ts.Auth.CraftJWT(log.Username)
	if err != nil {
		httputil.SendJSONError(w, "Failed Crafting jwt", http.StatusInternalServerError)

		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
