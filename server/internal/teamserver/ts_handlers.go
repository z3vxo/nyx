package teamserver

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/z3vxo/kronos/internal/database"
	"github.com/z3vxo/kronos/internal/httputil"
)

func (ts *TeamServer) AgentListHandler(w http.ResponseWriter, r *http.Request) {

	agents, err := ts.db.ListAgents()
	if err != nil {
		httputil.SendJSONError(w, "Failed retreiving agents", http.StatusInternalServerError)
		return
	}

	payload := database.Agents{Total: len(agents), Agent: agents}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&payload)

}

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

func (ts *TeamServer) CommandNewHandler(w http.ResponseWriter, r *http.Request) {
	var cmd TaskEntry
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		httputil.SendJSONError(w, "Error decoding json", http.StatusInternalServerError)
		return
	}

	err := ts.db.InsertCommand(cmd.Cmd_type, cmd.TaskID, cmd.Guid, cmd.Param1, cmd.Param2)
	if err != nil {
		httputil.SendJSONError(w, "failed inserting command", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

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

func (ts *TeamServer) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	guid := chi.URLParam(r, "guid")
	if guid == "" {
		httputil.SendJSONError(w, "missing guid", http.StatusBadRequest)
		return
	}

	tasks, err := ts.db.ListTasks(guid)
	if err != nil {
		httputil.SendJSONError(w, "database error, failed loading tasks", http.StatusInternalServerError)
		return
	}

	payload := database.TaskEntrys{Total: len(tasks), Tasks: tasks}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&payload)

}

func (ts *TeamServer) StartListenerHandler(w http.ResponseWriter, r *http.Request) {
	var Info NewListener
	if err := json.NewDecoder(r.Body).Decode(&Info); err != nil {
		httputil.SendJSONError(w, "Failed decoding json", http.StatusInternalServerError)
		return
	}

	id, err := ts.NewListener(Info.Port)
	if err != nil {
		fmt.Println(err)
		httputil.SendJSONError(w, "failed Creating listener", http.StatusInternalServerError)
		return
	}

	if err := ts.StartListener(id); err != nil {
		httputil.SendJSONError(w, "failed starting listener", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"listener_id": id})
}

func (ts *TeamServer) StopListenerHandler(w http.ResponseWriter, r *http.Request) {
	guid := chi.URLParam(r, "guid")
	if guid == "" {
		httputil.SendJSONError(w, "missing guid", http.StatusBadRequest)
		return
	}

	err := ts.StopListener(guid)
	if err != nil {
		httputil.SendJSONError(w, "failed deleting listener from db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

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

	refresh, err := ts.Auth.CraftRefreshJWT(log.Username)
	if err != nil {
		httputil.SendJSONError(w, "Failed Crafting refresh jwt", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token, "refresh": refresh})
}

func (ts *TeamServer) ListListenerHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ts.ListListeners()
	if err != nil {
		httputil.SendJSONError(w, "Failed Listing Listeners", http.StatusInternalServerError)
		return
	}
	res := ListListenersResp{
		Total:     len(data),
		Listeners: data,
	}

	json.NewEncoder(w).Encode(res)

}
