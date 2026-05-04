package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/BCoaracy/pacientezero/internal/db"
	"github.com/BCoaracy/pacientezero/internal/handler"
	"github.com/BCoaracy/pacientezero/internal/middleware"
	"github.com/BCoaracy/pacientezero/internal/repository"
	"github.com/BCoaracy/pacientezero/internal/service"
)

func main() {
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	pool, err := db.NewPool(context.Background())
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	jwtSecret := []byte(mustEnv("JWT_SECRET"))

	usuarioRepo := repository.NewUsuarioRepository(pool)
	pacienteRepo := repository.NewPacienteRepository(pool)

	usuarioSvc := service.NewUsuarioService(usuarioRepo, jwtSecret)
	pacienteSvc := service.NewPacienteService(pacienteRepo)

	usuarioH := handler.NewUsuarioHandler(usuarioSvc)
	pacienteH := handler.NewPacienteHandler(pacienteSvc)

	auth := middleware.Auth(jwtSecret)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Rotas públicas
	mux.HandleFunc("POST /login", usuarioH.Login)
	mux.HandleFunc("POST /usuarios", usuarioH.Create)

	// Rotas protegidas — usuários
	mux.Handle("GET /usuarios", auth(http.HandlerFunc(usuarioH.List)))
	mux.Handle("GET /usuarios/{id}", auth(http.HandlerFunc(usuarioH.GetByID)))
	mux.Handle("PUT /usuarios/{id}", auth(http.HandlerFunc(usuarioH.Update)))
	mux.Handle("DELETE /usuarios/{id}", auth(http.HandlerFunc(usuarioH.Delete)))

	// Rotas protegidas — pacientes
	mux.Handle("POST /pacientes", auth(http.HandlerFunc(pacienteH.Create)))
	mux.Handle("GET /pacientes", auth(http.HandlerFunc(pacienteH.List)))
	mux.Handle("GET /pacientes/{id}", auth(http.HandlerFunc(pacienteH.GetByID)))
	mux.Handle("PUT /pacientes/{id}", auth(http.HandlerFunc(pacienteH.Update)))
	mux.Handle("DELETE /pacientes/{id}", auth(http.HandlerFunc(pacienteH.Delete)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	root := middleware.RateLimit(60)(middleware.CORS(securityHeaders(mux)))

	log.Printf("API iniciando na porta %s", port)
	if err := http.ListenAndServe(":"+port, root); err != nil {
		log.Fatalf("falha ao iniciar servidor: %v", err)
	}
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("variável de ambiente obrigatória não definida: %s", key)
	}
	return v
}
