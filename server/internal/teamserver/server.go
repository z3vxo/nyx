package server

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/z3vxo/kronos/internal/config"
)

type TeamServer struct {
	Listener   net.Listener
	httpServer *http.Server
}

func NewTeamServer() *TeamServer {
	return &TeamServer{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", config.Cfg.TS.ListenInterface, config.Cfg.TS.Port),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 0,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (ts *TeamServer) Start() error {

	r := chi.NewRouter()

	r.Route("/rest", func(r chi.Router) {
		r.Post("/login", loginHandler)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleWare)
			r.Get("/agents/list", nyx_AgentListHandler)
			r.Get("/agents/resolve/{codename}", nyx_AgentResolveHandler)

			r.Post("/commands/new", nyx_CommandNewHandler)
			// r.Post("/commands/delete", nyx_CommandDeleteHandler)
			//
			// r.Post("/listeners/start, nyx_StartListenerHandler)
			// r.Post("/listeners/stop, nyx_StopListenerHandler)

		})
	})
	fmt.Println("Server Started!")

	return ts.httpServer.ListenAndServe()
}
