package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/RevanAdr/zlen/internal/auth"
	"github.com/RevanAdr/zlen/internal/config"
	"github.com/RevanAdr/zlen/internal/httpx"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db := config.NewConfig()
	repo := auth.NewUserRepository(db)
	service := auth.NewUserService(repo)
	handler := auth.NewUserHandler(service)
	route := auth.NewUseRoute(handler)

	if env := godotenv.Load(".env"); env != nil {
		panic(env)
	}

	key := os.Getenv("JWT_KEY")
	if len(key) < 32 {
		log.Fatal("Jwt key too short or missing!")
	}
	jwt := []byte(key)
	router := route.NewRouter()

	server := &http.Server{
		Addr:              ":8000",
		Handler:           httpx.Logging(router),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logger.Info("listening", "addr", server.Addr)
	log.Fatal(server.ListenAndServe())
}
