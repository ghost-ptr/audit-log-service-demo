package main

import (
	"net/http"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK!"))
}

func main() {
	auditServeMux := http.NewServeMux()
	auditServeMux.HandleFunc("/", httpHandler)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: auditServeMux,
	}
	httpServer.ListenAndServe()
}
