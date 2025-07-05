package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "merchshop/cmd/docs"

	_ "github.com/lib/pq"

	"merchshop/internal/api/http/auth"
	"merchshop/internal/api/http/handlers"
	"merchshop/internal/api/http/router"
	"merchshop/internal/config"
	"merchshop/internal/repository"
	"merchshop/internal/usecase"
)

// @title MerchShop API
// @version 1.0
// @description API для мерчшопа

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Загрузка конфигурации
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	//Инициализация бд
	db, err := initializeDatabase(cfg.DB.DSN())
	if err != nil {
		log.Fatalf("failed to initialize  database: %v", err)
	}

	// // Чтение файла migrations.sql
	// migrationSQL, err := os.ReadFile("../migrations/migrations.sql")
	// if err != nil {
	// 	panic(fmt.Sprintf("err read: %v", err))
	// }

	// ctx := context.Background()
	// _, err = db.ExecContext(ctx, string(migrationSQL))
	// if err != nil {
	// 	panic(fmt.Sprintf("err migr: %v", err))
	// }

	// fmt.Println("migr use!")

	// Инициализация репозиториев
	repo := repository.NewRepositories(db)

	// Инициализация JWT manager
	tokenManager, err := auth.NewJWTManager(cfg.Auth.SigningKey, cfg.Auth.TokenTTL)
	if err != nil {
		if cerr := db.Close(); cerr != nil {
			log.Printf("Ошибка при закрытии БД: %v", cerr)
		}

		db.Close()
		log.Fatalf("failed to initialize token manager: %v", err)
	}

	defer db.Close()

	// Инициализация use cases
	useCases := usecase.NewUseCases(repo)

	// Инициализация хендлеров
	handler := handlers.NewHandler(
		useCases.User,
		useCases.Transaction,
		useCases.Purchase,
		useCases.Merch,
		tokenManager,
	)

	// Инициализация роутера
	httpRouter := router.NewRouter(handler, tokenManager)

	// Запуск HTTP сервера
	startServer(httpRouter, cfg.Server.Port, cfg.Server.ReadTimeout, cfg.Server.WriteTimeout)

}

func loadConfig() (*config.Config, error) {
	cfg, err := config.LoadConfig("../configs")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}

func initializeDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func startServer(r http.Handler, port int, readTimeout, writeTimeout time.Duration) {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	go func() {
		log.Println("Server is starting")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		cancel()
		log.Fatalf("Server Shutdown: %v", err)
	}

	cancel()

	log.Println("Server exiting")
}
