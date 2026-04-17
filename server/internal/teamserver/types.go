package teamserver

import (
	"net"
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
	TaskID   string `json:"task_id"`
	Param1   string `json:"param_1"`
	Param2   string `json:"param_2"`
}

type UserLogin struct {
	Username string `json:"user"`
	Password string `json:"passwd"`
}

type Broker struct {
	Channels map[string]chan string
	mu       sync.RWMutex
}

type TeamServer struct {
	Listener   net.Listener
	httpServer *http.Server
	SSE        *Broker
	Auth       *auth.Auth
	db         *database.DB
}
