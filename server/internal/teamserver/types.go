package teamserver

import (
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
	Port int `json:"port"`
}

type TeamServer struct {
	httpServer *http.Server
	SSE        *Broker
	Auth       *auth.Auth
	db         *database.DB
	Listeners  *Listeners
}

type ListenerEntry struct {
	Port   int
	Name   string
	Status string
}

type ListListenersResp struct {
	Total     int             `json:"total"`
	Listeners []ListenerEntry `json:"listeners"`
}
