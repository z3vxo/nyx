package teamserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/z3vxo/kronos/internal/auth"
	"github.com/z3vxo/kronos/internal/httputil"
)

func (ts *TeamServer) NameToGuid(name string) string {
	ts.Listeners.Mu.Lock()
	var id string
	for i, l := range ts.Listeners.ListenerMap {
		if l.Name == name {
			id = i
		}
	}
	ts.Listeners.Mu.Unlock()
	return id
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

func (ts *TeamServer) StartListenerHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		httputil.SendJSONError(w, "missing listener name", http.StatusBadRequest)
		return
	}

	id := ts.NameToGuid(name)
	fmt.Printf("NAME: %s\nID %s\n", name, id)
	if id == "" {
		httputil.SendJSONError(w, "listener not found", http.StatusNotFound)
		return
	}

	if err := ts.StartListener(id); err != nil {
		httputil.SendJSONError(w, "failed starting Listener", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})

}

func (ts *TeamServer) NewListenerHandler(w http.ResponseWriter, r *http.Request) {
	var Info NewListener
	if err := json.NewDecoder(r.Body).Decode(&Info); err != nil {
		httputil.SendJSONError(w, "Failed decoding json", http.StatusInternalServerError)
		return
	}

	user, _ := r.Context().Value(auth.UsernameKey).(string)
	fmt.Println(Info.Protocol)
	id, name, err := ts.NewListener(Info.Port, Info.Protocol, user, Info.Host, Info.CertType)
	if err != nil {
		fmt.Println(err)
		httputil.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := ts.StartListener(id); err != nil {
		httputil.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"listener_name": name})
}

func (ts *TeamServer) DeleteListnerHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		httputil.SendJSONError(w, "missing listener name", http.StatusBadRequest)
		return
	}

	id := ts.NameToGuid(name)
	if id == "" {
		httputil.SendJSONError(w, "listener not found", http.StatusNotFound)
		return
	}

	if err := ts.DeleteListner(id); err != nil {
		httputil.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

func (ts *TeamServer) StopListenerHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		httputil.SendJSONError(w, "missing listener name", http.StatusBadRequest)
		return
	}

	guid := ts.NameToGuid(name)
	if guid == "" {
		httputil.SendJSONError(w, "listener not found", http.StatusNotFound)
		return
	}
	user, _ := r.Context().Value(auth.UsernameKey).(string)
	err := ts.StopListener(guid, user)
	if err != nil {
		httputil.SendJSONError(w, "failed deleting listener from db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}
