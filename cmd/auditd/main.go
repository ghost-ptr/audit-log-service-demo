package main

import (
	"encoding/json"
	"fmt"
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

func httpHandler(w http.ResponseWriter, r *http.Request) {
	var request EventRequest
	var event Event
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	currentTime := time.Now()

	event.EventRequest = request
	event.ReceivedAt = currentTime

	if event.OccurredAt.IsZero() {
		event.OccurredAt = currentTime
		event.OccurredAtInferred = true
	}

	if event.User == "" || event.Action == "" || event.Resource == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid event"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK!"))

	fmt.Printf("Received event: %+v\n", event)
}

func main() {
	auditServeMux := http.NewServeMux()
	auditServeMux.HandleFunc("POST /events", httpHandler)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: auditServeMux,
	}
	httpServer.ListenAndServe()
}
