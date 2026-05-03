package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("API iniciando na porta %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("falha ao iniciar servidor: %v", err)
	}
}
