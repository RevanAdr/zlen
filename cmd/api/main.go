package main

import (
	"github.com/RevanAdr/zlen/internal/auth"
	"github.com/RevanAdr/zlen/internal/config"
)

func main() {
	db := config.NewConfig()
	repo := auth.NewAuthRepository(db)
}
