package teamserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/z3vxo/kronos/internal/config"
	"github.com/z3vxo/kronos/internal/server"
)

type Listener struct {
	httpServer *http.Server
	Port       int
}

type Listeners struct {
	Mu           sync.RWMutex
	ListenerMap  map[string]Listener
	GetEndpoint  string
	PostEndpoint string
}

func (ts *TeamServer) NewListener(port int) (string, error) {
	id := uuid.NewString()

	ts.Listeners.Mu.Lock()
	for _, l := range ts.Listeners.ListenerMap {
		if l.Port == port {
			ts.Listeners.Mu.Unlock()
			return "", errors.New("already Listening on port")
		}
	}

	r := chi.NewRouter()
	r.Get(config.Cfg.Server.GetEndpoint, server.AgentCheckInHandler)
	r.Post(config.Cfg.Server.PostEndpoint, server.AgentUploadHandler)

	ts.Listeners.ListenerMap[id] = Listener{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      r,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Port: port,
	}
	ts.Listeners.Mu.Unlock()

	if err := ts.db.InsertListener(port, id); err != nil {
		ts.Listeners.Mu.Lock()
		delete(ts.Listeners.ListenerMap, id)
		ts.Listeners.Mu.Unlock()
		return "", err
	}

	return id, nil

}

func (ts *TeamServer) StartListener(id string) error {
	ts.Listeners.Mu.RLock()
	l, ok := ts.Listeners.ListenerMap[id]
	ts.Listeners.Mu.RUnlock()
	if !ok {
		return errors.New("listener not found")
	}

	go func() {
		if err := l.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listener %s error: %v\n", id, err)
		}
	}()

	return nil
}

func (ts *TeamServer) StopListener(id string) error {
	ts.Listeners.Mu.Lock()
	l, ok := ts.Listeners.ListenerMap[id]
	if !ok {
		ts.Listeners.Mu.Unlock()
		return errors.New("listener not found")
	}
	delete(ts.Listeners.ListenerMap, id)
	ts.Listeners.Mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	l.httpServer.Shutdown(ctx)

	if err := ts.db.DeleteListener(id); err != nil {
		return err
	}
	return nil

}

func (ts *TeamServer) StopAllListeners() {
	ts.Listeners.Mu.Lock()
	listeners := make([]Listener, 0, len(ts.Listeners.ListenerMap))
	for _, l := range ts.Listeners.ListenerMap {
		listeners = append(listeners, l)
	}
	ts.Listeners.Mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, l := range listeners {
		l.httpServer.Shutdown(ctx)
	}
}
