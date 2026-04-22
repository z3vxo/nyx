package teamserver

import (
	"log/slog"
	"net/http"

	"github.com/z3vxo/kronos/internal/auth"
	"github.com/z3vxo/kronos/internal/broker"
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

type NewListener struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	CertType bool   `json:"letsencrypt"` // 0 = self signed, 1 = lets encrypt
}

type TeamServer struct {
	httpServer *http.Server
	SSE        *broker.Broker
	Auth       *auth.Auth
	db         *database.DB
	Listeners  *Listeners
	Logger     *slog.Logger
}

type ListenerEntry struct {
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Status   bool   `json:"status"`
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
	PPid         int32  `json:"ppid"`
	WinVer       string `json:"win_version`
	InternalIP   string `json:"internal_ip"`
	ExternalIP   string `json:"external_ip"`
	IsElevated   bool   `json:"is_elev`
	Arch         byte   `json:"arch"`
	LastCheckin  int64  `json:"last_checkin"`
	RegisterTime int64  `json:"reg_date`
}
