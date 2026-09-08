package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	fmt.Printf("Received event: %+v\n", event)
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
		log.Fatalf("Error starting server: %v\n", err)
	}
}
