package main

import (
	"log"
	"net/http"

	"github.com/RevanAdr/zlen/internal/auth"
	"github.com/RevanAdr/zlen/internal/config"
)

func main() {
	db := config.NewConfig()
	repo := auth.NewUserRepository(db)
	service := auth.NewUserService(repo)
	handler := auth.NewUserHandler(service)
	route := auth.NewUseRoute(handler)

	router := route.NewRouter()
	log.Println("Server listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
