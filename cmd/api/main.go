package main

import (
	"log"
	"net/http"
	"os"

	"github.com/BCoaracy/pacientezero/internal/db"
)

func main() {
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("migrations: %v", err)
	}

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