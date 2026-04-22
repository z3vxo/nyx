package teamserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/z3vxo/kronos/internal/database"
	"github.com/z3vxo/kronos/internal/httputil"
)

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
