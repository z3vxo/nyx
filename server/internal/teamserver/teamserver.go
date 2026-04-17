package teamserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/z3vxo/kronos/internal/auth"
	"github.com/z3vxo/kronos/internal/config"
	"github.com/z3vxo/kronos/internal/database"
)

func NewTeamServer() (*TeamServer, error) {
	a := auth.NewAuth(config.Cfg.TS.Auth.Username, config.Cfg.TS.Auth.Password,
		config.Cfg.TS.Auth.JwtSecret, config.Cfg.TS.Auth.TokenHours)
	d, err := database.NewDB()
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
		SSE:  NewBroker(),
		Auth: a,
		db:   d,
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
		r.Get("/events", ts.SSE.EventHandler)

		r.Group(func(r chi.Router) {
			r.Use(ts.Auth.AuthMiddleWare)
			r.Get("/rest/agents/list", ts.AgentListHandler)
			r.Get("/rest/agents/resolve/{codename}", ts.AgentResolveHandler)

			r.Post("/rest/commands/new", ts.CommandNewHandler)
			r.Post("/rest/commands/delete", ts.CommandDeleteHandler)

			//r.Get("/rest/listeners/list", ts_ListListener)
			//r.Post("/rest/listeners/start, ts_StartListenerHandler)
			//r.Post("/rest/listeners/stop, ts_StopListenerHandler)
			//
			// r.Post("/listeners/start, nyx_StartListenerHandler)
			// r.Post("/listeners/stop, nyx_StopListenerHandler)

		})
	})
	fmt.Println("Server Started!")

	return ts.httpServer.ListenAndServe()
}

func (ts *TeamServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ts.httpServer.Shutdown(ctx)
}
