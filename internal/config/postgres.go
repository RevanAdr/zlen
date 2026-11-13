package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Env struct {
	User     string
	Password string
	Port     int
	Host     string
	Db       string
}

type Config struct {
}

func NewConfig() *pgxpool.Pool {
	if loadErr := godotenv.Load(".env"); loadErr != nil {
		panic(loadErr)
	}

	env := Env{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Port:     5432,
		Host:     "localhost",
		Db:       os.Getenv("POSTGRES_DB"),
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", env.User, env.Password, env.Host, env.Port, env.Db)

	pool, poolErr := pgxpool.New(context.Background(), dbUrl)

	if poolErr != nil {
		panic(poolErr)
	}

	return pool
}
