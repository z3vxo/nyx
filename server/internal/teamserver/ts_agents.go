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
		fmt.Println(err)
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

func (ts *TeamServer) AgentInfoHandler(w http.ResponseWriter, r *http.Request) {
	Name := chi.URLParam(r, "codename")
	if Name == "" {
		httputil.SendJSONError(w, "Missing Codename", http.StatusBadRequest)
		return
	}

	AgentInfo, err := ts.db.ListAgentInfo(Name)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			httputil.SendJSONError(w, "agent not found", http.StatusNotFound)
			return
		}
		httputil.SendJSONError(w, "database error", http.StatusInternalServerError)
		return
	}

	resp := AgentInfoResp{
		User:         AgentInfo.User,
		Host:         AgentInfo.Host,
		ProcPath:     AgentInfo.ProcPath,
		Pid:          AgentInfo.Pid,
		WinVer:       AgentInfo.WinVer,
		InternalIP:   AgentInfo.InternalIP,
		ExternalIP:   AgentInfo.ExternalIP,
		IsElevated:   AgentInfo.IsElevated,
		LastCheckin:  AgentInfo.LastCheckin,
		RegisterTime: AgentInfo.RegisterTime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)

}
