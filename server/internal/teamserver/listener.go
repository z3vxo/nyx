package teamserver

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/z3vxo/kronos/internal/config"
	"github.com/z3vxo/kronos/internal/server"
)

var listenerAdjectives = []string{"silent", "ghost", "iron", "shadow", "crimson", "steel", "hollow", "obsidian", "frozen", "ashen"}
var listenerNouns = []string{"falcon", "specter", "anvil", "wraith", "ember", "vault", "cipher", "nexus", "bastion", "phantom"}

func generateListenerName() string {
	adj := listenerAdjectives[rand.Intn(len(listenerAdjectives))]
	noun := listenerNouns[rand.Intn(len(listenerNouns))]
	return fmt.Sprintf("%s-%s", adj, noun)
}

type Listener struct {
	httpServer *http.Server
	Port       int
	Name       string
	Protocol   string
}

type Listeners struct {
	Mu           sync.RWMutex
	ListenerMap  map[string]Listener
	GetEndpoint  string
	PostEndpoint string
}

func BuildListenerHttp(port int, protocol string) *http.Server {
	r := chi.NewRouter()
	r.Get(config.Cfg.Server.GetEndpoint, server.AgentCheckInHandler)
	r.Post(config.Cfg.Server.PostEndpoint, server.AgentUploadHandler)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

}

func (ts *TeamServer) StartListenersFromDB() error {
	ToStart, err := ts.db.GetListeners()
	if err != nil {
		return err
	}
	for _, l := range ToStart {
		ts.Listeners.Mu.Lock()
		ts.Listeners.ListenerMap[l.Guid] = Listener{
			httpServer: BuildListenerHttp(l.Port, l.Protocol),
			Port:       l.Port,
			Name:       l.Name,
			Protocol:   l.Protocol,
		}
		ts.Listeners.Mu.Unlock()

		if err := ts.StartListener(l.Guid); err != nil {
			ts.Listeners.Mu.Lock()
			delete(ts.Listeners.ListenerMap, l.Guid)
			ts.Listeners.Mu.Unlock()
			fmt.Printf("failed starting listener %s: %v\n", l.Guid, err)
			continue
		}
	}
	return nil
}

func (ts *TeamServer) NewListener(port int, Protocol string) (string, string, error) {
	id := uuid.NewString()
	name := generateListenerName()

	ts.Listeners.Mu.Lock()
	for _, l := range ts.Listeners.ListenerMap {
		if l.Port == port {
			ts.Listeners.Mu.Unlock()
			return "", "", errors.New("already Listening on port")
		}
	}

	ts.Listeners.ListenerMap[id] = Listener{
		httpServer: BuildListenerHttp(port, Protocol),
		Port:       port,
		Name:       name,
		Protocol:   Protocol,
	}
	ts.Listeners.Mu.Unlock()

	if err := ts.db.InsertListener(port, id, name, Protocol); err != nil {
		ts.Listeners.Mu.Lock()
		delete(ts.Listeners.ListenerMap, id)
		ts.Listeners.Mu.Unlock()
		return "", "", err
	}

	return id, name, nil
}

func (ts *TeamServer) StartListener(id string) error {
	ts.Listeners.Mu.RLock()
	l, ok := ts.Listeners.ListenerMap[id]
	ts.Listeners.Mu.RUnlock()
	if !ok {
		return errors.New("listener not found")
	}

	go func() {
		var err error
		if l.Protocol == "https" {
			err = l.httpServer.ListenAndServeTLS(config.Cfg.Server.Cert, config.Cfg.Server.Key)
		} else {
			err = l.httpServer.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
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

func (ts *TeamServer) ListListeners() ([]ListenerEntry, error) {
	ts.Listeners.Mu.RLock()
	defer ts.Listeners.Mu.RUnlock()

	var listener []ListenerEntry
	for _, i := range ts.Listeners.ListenerMap {
		listener = append(listener, ListenerEntry{
			Port:     i.Port,
			Name:     i.Name,
			Protocol: i.Protocol,
			Status:   "running",
		})
	}

	return listener, nil
}
