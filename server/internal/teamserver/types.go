package teamserver

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/z3vxo/kronos/internal/auth"
	"github.com/z3vxo/kronos/internal/database"
)

type TaskDelete struct {
	TaskID int `json:"task_id"`
}

type TaskEntry struct {
	Cmd_type int    `json:"type"`
	Guid     string `json:"guid"`
	TaskID   int    `json:"task_id"`
	Param1   string `json:"param_1"`
	Param2   string `json:"param_2"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Broker struct {
	Channels map[string]chan string
	mu       sync.RWMutex
}

type NewListener struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type TeamServer struct {
	httpServer *http.Server
	SSE        *Broker
	Auth       *auth.Auth
	db         *database.DB
	Listeners  *Listeners
	Logger     *slog.Logger
}

type ListenerEntry struct {
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Status   string `json:"status"`
}

type ListListenersResp struct {
	Total     int             `json:"total"`
	Listeners []ListenerEntry `json:"listeners"`
}

type AgentInfoResp struct {
	User         string `json:"username"`
	Host         string `json:"hostname"`
	ProcPath     string `json:"proc_path`
	Pid          int32  `json:"pid"`
	WinVer       string `json:"win_version`
	InternalIP   string `json:"internal_ip"`
	ExternalIP   string `json:"external_ip"`
	IsElevated   bool   `json:"is_elev`
	LastCheckin  int64  `json:"last_checkin"`
	RegisterTime int64  `json:"reg_date`
}
