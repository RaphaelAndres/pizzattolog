package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/pizzattolog/licencas/internal/auth"
	"github.com/pizzattolog/licencas/internal/handlers"
	"github.com/pizzattolog/licencas/internal/middleware"
	"github.com/pizzattolog/licencas/internal/models"
	"github.com/pizzattolog/licencas/internal/repository"
	"github.com/pizzattolog/licencas/internal/services"
)

func main() {
	// Carrega .env se existir
	_ = godotenv.Load()

	// Conecta ao banco
	db, err := models.Connect()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}

	// Auto-migrate das entidades
	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("Erro na migration: %v", err)
	}

	// Seed do usuário admin padrão
	if err := models.SeedAdmin(db); err != nil {
		log.Printf("Aviso no seed: %v", err)
	}

	// Inicializa MinIO
	minioClient, err := services.NewMinioClient()
	if err != nil {
		log.Fatalf("Erro ao conectar ao MinIO: %v", err)
	}

	// Repositórios
	userRepo := repository.NewUserRepository(db)
	licencaRepo := repository.NewLicencaRepository(db)

	// Serviços
	jwtService := auth.NewJWTService()
	userService := services.NewUserService(userRepo, jwtService)
	licencaService := services.NewLicencaService(licencaRepo, minioClient)
	alertaService := services.NewAlertaService(licencaRepo)

	// Inicia cron de alertas
	alertaService.StartCron()

	// Handlers
	authHandler := handlers.NewAuthHandler(userService)
	licencaHandler := handlers.NewLicencaHandler(licencaService)
	dashboardHandler := handlers.NewDashboardHandler(licencaService)

	// Router
	router := middleware.SetupRouter(authHandler, licencaHandler, dashboardHandler, jwtService)

	// Servidor HTTP
	srv := &http.Server{
		Addr:         ":" + getEnv("APP_PORT", "8080"),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Inicia servidor em goroutine
	go func() {
		log.Printf("🚀 Servidor rodando na porta %s", getEnv("APP_PORT", "8080"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro no servidor: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Encerrando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro no shutdown: %v", err)
	}

	log.Println("Servidor encerrado com sucesso.")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
