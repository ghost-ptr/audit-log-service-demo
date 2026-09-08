package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type EventRequest struct {
	User       string    `json:"user"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	OccurredAt time.Time `json:"occurred_at"`
	SourceIP   string    `json:"source_ip"`
}

type Event struct {
	EventRequest
	ReceivedAt         time.Time `json:"received_at"`
	OccurredAtInferred bool      `json:"occurred_at_inferred"`
}

type EventStore struct {
	mu     sync.Mutex
	events []Event
}

func (store *EventStore) Add(event Event) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.events = append(store.events, event)
}

func (store *EventStore) Read() []Event {
	store.mu.Lock()
	defer store.mu.Unlock()
	snapshots := make([]Event, len(store.events))
	copy(snapshots, store.events)
	return snapshots
}

func NewEvent(request EventRequest) Event {
	var event Event
	currentTime := time.Now()

	event.EventRequest = request
	event.ReceivedAt = currentTime

	if event.OccurredAt.IsZero() {
		event.OccurredAt = currentTime
		event.OccurredAtInferred = true
	}

	return event
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	var request EventRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event := NewEvent(request)
	if event.User == "" || event.Action == "" || event.Resource == "" {
		http.Error(w, "Invalid event", http.StatusBadRequest)
		return
	}
	//auditStore.Add(event)

	log.Printf("Received event: %+v", event)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	auditServeMux := http.NewServeMux()
	auditServeMux.HandleFunc("POST /events", httpHandler)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: auditServeMux,
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
