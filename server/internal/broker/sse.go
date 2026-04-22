package broker

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func NewBroker() *Broker {
	return &Broker{
		Channels: make(map[string]chan string),
	}
}

type Broker struct {
	Channels map[string]chan string
	mu       sync.RWMutex
}

func (b *Broker) AddSubscriber() (string, chan string) {
	id := fmt.Sprintf("%016x", rand.Uint64())
	fmt.Printf("New Sub: %s\n", id)
	ch := make(chan string, 8)
	b.mu.Lock()
	b.Channels[id] = ch
	b.mu.Unlock()

	return id, ch

}

func (b *Broker) RemoveSubscriber(id string) {
	b.mu.Lock()
	if ch, ok := b.Channels[id]; ok {
		close(ch)
		delete(b.Channels, id)
	}
	b.mu.Unlock()

}

func (b *Broker) Broadcast(msg string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.Channels {
		select {
		case ch <- msg:
		default:
		}
	}
}

func (b *Broker) EventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE Not Supported", http.StatusInternalServerError)
		return
	}

	id, ch := b.AddSubscriber()
	defer b.RemoveSubscriber(id)

	heart := time.NewTicker(15 * time.Second)
	defer heart.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-heart.C:
			data, _ := json.Marshal(map[string]int{"type": 3})
			fmt.Fprint(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}

}
