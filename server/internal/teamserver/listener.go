package teamserver

import (
	"context"
	"crypto/tls"
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

var listenerAdjectives = []string{"golden", "cursed", "sacred", "eternal", "fallen", "divine", "forsaken", "ancient", "wrathful", "fated", "hollow", "sunken", "boundless", "immortal", "exiled"}
var listenerNouns = []string{"olympus", "tartarus", "elysium", "styx", "erebus", "acheron", "lethe", "phlegethon", "asphodel", "cocytus", "ithaca", "troy", "delphi", "sparta", "thebes"}

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
	Host       string
	Status     bool
	certType   bool
}

type Listeners struct {
	Mu           sync.RWMutex
	ListenerMap  map[string]Listener
	GetEndpoint  string
	PostEndpoint string
}

func (ts *TeamServer) UpdateListenerMapStatus(id string, status bool) {
	ts.Listeners.Mu.Lock()
	entry := ts.Listeners.ListenerMap[id]
	entry.Status = status
	ts.Listeners.ListenerMap[id] = entry
	ts.Listeners.Mu.Unlock()
}

func BuildListenerHttp(port int, protocol string, h *server.AgentHandler, host string, letsEncrypt bool) (*http.Server, error) {
	r := chi.NewRouter()
	r.Get(config.Cfg.Server.GetEndpoint, h.AgentCheckInHandler)
	r.Post(config.Cfg.Server.PostEndpoint, h.AgentUploadHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if protocol == "https" {
		if !letsEncrypt {
			cert, err := GenSelSigned(host)
			if err != nil {
				return nil, errors.New("failed generating self signed cert")
			}
			srv.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
				MinVersion:   tls.VersionTLS12,
			}
		}
	}
	return srv, nil

}

func (ts *TeamServer) StartListenersFromDB() error {
	ToStart, err := ts.db.GetListeners()
	if err != nil {
		return err
	}
	h := &server.AgentHandler{DB: ts.db, Broker: ts.SSE}

	for _, l := range ToStart {

		srv, err := BuildListenerHttp(l.Port, l.Protocol, h, l.Host, l.CertType)
		if err != nil {
			return err
		}
		ts.Listeners.Mu.Lock()
		ts.Listeners.ListenerMap[l.Guid] = Listener{
			httpServer: srv,
			Port:       l.Port,
			Name:       l.Name,
			Host:       l.Host,
			Protocol:   l.Protocol,
			Status:     false,
			certType:   l.CertType,
		}
		ts.Listeners.Mu.Unlock()

		if l.Status == true {
			if err := ts.StartListener(l.Guid); err != nil {
				ts.Listeners.Mu.Lock()
				delete(ts.Listeners.ListenerMap, l.Guid)
				ts.Listeners.Mu.Unlock()
				fmt.Printf("failed starting listener %s: %v\n", l.Guid, err)
				continue
			}
		}

		ts.Logger.Info("listener started from db", "event", "listener-restore", "id", l.Guid, "protocol", l.Protocol, "port", l.Port)
	}
	return nil
}

func (ts *TeamServer) NewListener(port int, Protocol string, user, host string, letsEncrypt bool) (string, string, error) {
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
		Port:     port,
		Name:     name,
		Host:     host,
		Protocol: Protocol,
		Status:   false,
	}
	ts.Listeners.Mu.Unlock()

	if err := ts.db.InsertListener(port, id, name, Protocol, host, letsEncrypt, false); err != nil {
		ts.Listeners.Mu.Lock()
		delete(ts.Listeners.ListenerMap, id)
		ts.Listeners.Mu.Unlock()
		return "", "", err
	}
	ts.Logger.Info("New Listener Created", "event", "listener-create", "id", id, "operator", user, "proto", Protocol, "port", port)

	return id, name, nil
}

func (ts *TeamServer) StartListener(id string) error {
	ts.Listeners.Mu.RLock()
	l, ok := ts.Listeners.ListenerMap[id]
	ts.Listeners.Mu.RUnlock()
	if !ok {
		return errors.New("listener not found")
	}

	if l.Status == true {
		return errors.New("listener already running!")
	}

	srv, err := BuildListenerHttp(l.Port, l.Protocol, &server.AgentHandler{DB: ts.db, Broker: ts.SSE}, l.Host, l.certType)
	if err != nil {
		return errors.New("Failed Creating server object")
	}
	ts.Listeners.Mu.Lock()
	l.httpServer = srv
	ts.Listeners.ListenerMap[id] = l
	ts.Listeners.Mu.Unlock()

	errCh := make(chan error, 1)
	go func() {
		var err error
		if l.Protocol == "https" {
			err = l.httpServer.ListenAndServeTLS("", "")
		} else {
			err = l.httpServer.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			ts.db.DeleteListener(id)
			errCh <- err
		}
	}()
	ts.UpdateListenerMapStatus(id, true)
	ts.db.UpdateListenerStatus(id, true)

	select {
	case err := <-errCh:
		return err
	case <-time.After(100 * time.Millisecond):
		return nil
	}
}

func (ts *TeamServer) DeleteListner(id string) error {
	ts.Listeners.Mu.Lock()
	l, ok := ts.Listeners.ListenerMap[id]
	if !ok {
		ts.Listeners.Mu.Unlock()
		return errors.New("listener not found")
	}
	delete(ts.Listeners.ListenerMap, id)
	ts.Listeners.Mu.Unlock()

	if l.Status == true {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		l.httpServer.Shutdown(ctx)
	}

	if err := ts.db.DeleteListener(id); err != nil {
		fmt.Println(err)
		return err
	}
	delete(ts.Listeners.ListenerMap, id)

	return nil
}

func (ts *TeamServer) StopListener(id string, user string) error {
	ts.Listeners.Mu.Lock()
	l, ok := ts.Listeners.ListenerMap[id]
	if !ok {
		ts.Listeners.Mu.Unlock()
		return errors.New("listener not found")
	}
	//delete(ts.Listeners.ListenerMap, id)
	ts.Listeners.Mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	l.httpServer.Shutdown(ctx)

	if err := ts.db.UpdateListenerStatus(id, false); err != nil {
		fmt.Println(err)
		return err
	}

	ts.UpdateListenerMapStatus(id, false)

	ts.Logger.Info("listener stopped", "event", "listener-stop", "id", id, "operator", user)

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
			Host:     i.Host,
			Status:   i.Status,
		})
	}

	return listener, nil
}
