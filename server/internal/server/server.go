package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupTeamServer() error {

	r := chi.NewRouter()

	r.Route("/rest", func(r chi.Router) {
		r.Post("/login", loginHandler)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleWare)
			r.Get("/agents/list", nyx_AgentHandler)
			r.Get("/agents/resolve/{codename}", nyx_AgentResolveHandler)

			r.Post("/commands/new", nyx_CommandNewHandler)
			// r.Post("/commands/delete", nyx_CommandDeleteHandler)
			//
			// r.Post("/listeners/start, nyx_StartListenerHandler)
			// r.Post("/listeners/stop, nyx_StopListenerHandler)

		})
	})
	fmt.Println("Server Started!")

	http.ListenAndServe(":3000", r)

	return nil
}
