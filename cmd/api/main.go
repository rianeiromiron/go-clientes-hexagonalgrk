package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	httphandler "github.com/grok/crudclienteshex/internal/adapters/http"
	"github.com/grok/crudclienteshex/internal/adapters/repository"
	"github.com/grok/crudclienteshex/internal/application"
	"github.com/grok/crudclienteshex/internal/infrastructure"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if present
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:norimorienair4614@localhost:5432/crudclientesgrok?sslmode=disable")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := infrastructure.NewPostgresPool(ctx, dbURL)
	if err != nil {
		log.Fatalf("Error conectando a PostgreSQL: %v\nAsegúrate de que el contenedor esté corriendo y la base de datos 'crudclientesgrok' exista.", err)
	}
	defer pool.Close()

	// Repository (driven adapter)
	repo := repository.NewPostgresClientRepository(pool)

	// Ensure schema
	if err := repo.EnsureSchema(context.Background()); err != nil {
		log.Fatalf("Error creando esquema: %v", err)
	}

	// Application service (use cases)
	service := application.NewClientService(repo)

	// HTTP handler (driving adapter)
	handler := httphandler.NewClientHandler(service)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Serve static frontend
	webDir := filepath.Join(".", "web")
	fs := http.FileServer(http.Dir(webDir))
	mux.Handle("/", fs)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		fmt.Printf("🚀 Servidor escuchando en http://localhost:%s\n", port)
		fmt.Printf("📁 Frontend: http://localhost:%s/\n", port)
		fmt.Printf("🔌 API:     http://localhost:%s/api/clients\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error del servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nApagando servidor...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Error en shutdown: %v", err)
	}
	fmt.Println("Servidor detenido correctamente.")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
