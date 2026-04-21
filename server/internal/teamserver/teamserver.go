package teamserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/z3vxo/kronos/internal/auth"
	"github.com/z3vxo/kronos/internal/config"
	"github.com/z3vxo/kronos/internal/database"
)

func GetLogFile() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kronos", "logs", "kronos.log")
}

func NewTeamServer() (*TeamServer, error) {
	a := auth.NewAuth(config.Cfg.TS.Auth.Username, config.Cfg.TS.Auth.Password,
		config.Cfg.TS.Auth.JwtSecret, config.Cfg.TS.Auth.TokenHours, config.Cfg.TS.Auth.TokenRefreshHours)
	d, err := database.NewDB()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(GetLogFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &TeamServer{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", config.Cfg.TS.ListenInterface, config.Cfg.TS.Port),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 0,
			IdleTimeout:  60 * time.Second,
		},
		SSE:       NewBroker(),
		Auth:      a,
		db:        d,
		Listeners: &Listeners{ListenerMap: make(map[string]Listener), GetEndpoint: config.Cfg.Server.GetEndpoint, PostEndpoint: config.Cfg.Server.PostEndpoint},
		Logger:    slog.New(slog.NewJSONHandler(file, nil)),
	}, nil
}

func NewBroker() *Broker {
	return &Broker{
		Channels: make(map[string]chan string),
	}
}

func (ts *TeamServer) Start() error {

	r := chi.NewRouter()
	ts.httpServer.Handler = r

	r.Route("/ts", func(r chi.Router) {
		r.Post("/rest/login", ts.loginHandler)

		r.Group(func(r chi.Router) {
			r.Use(ts.Auth.AuthMiddleWare)
			r.Get("/events", ts.SSE.EventHandler)
			r.Get("/rest/agents/list", ts.AgentListHandler)
			r.Get("/rest/agents/resolve/{codename}", ts.AgentResolveHandler)
			r.Get("/rest/agents/info/{codename}", ts.AgentInfoHandler)

			r.Post("/rest/tasks/new", ts.CommandNewHandler)
			r.Post("/rest/tasks/delete", ts.CommandDeleteHandler)
			r.Get("/rest/tasks/list/{guid}", ts.ListTasksHandler)

			r.Get("/rest/listeners/list", ts.ListListenerHandler)
			r.Post("/rest/listeners/start", ts.StartListenerHandler)
			r.Post("/rest/listeners/stop/{name}", ts.StopListenerHandler)

		})
	})
	fmt.Println("Server Started!")
	if err := ts.StartListenersFromDB(); err != nil {
		return err
	}

	return ts.httpServer.ListenAndServe()
}

func (ts *TeamServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ts.StopAllListeners()
	ts.httpServer.Shutdown(ctx)
}
