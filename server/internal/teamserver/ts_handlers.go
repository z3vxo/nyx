package teamserver

import (
	"encoding/json"
	"net/http"

	"github.com/z3vxo/kronos/internal/httputil"
)

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

	ts.Logger.Debug("Operator logged in", "user", log.Username)

	json.NewEncoder(w).Encode(map[string]string{"token": token, "refresh": refresh})
}
