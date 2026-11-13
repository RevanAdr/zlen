package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	List(ctx context.Context) ([]User, error)
}

type AuthRepositoryImpl struct {
	Db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &AuthRepositoryImpl{Db: db}
}

func (r *AuthRepositoryImpl) List(ctx context.Context) ([]User, error) {
	users := []User{}
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return users, err
	}

	results, qErr := tx.Query(ctx, "SELECT * FROM users;")

	if qErr != nil {
		return users, err
	}

	for results.Next() {
		result := User{}
		if err := results.Scan(result.Id, result.Name, result.Role); err != nil {
			return users, err
		}
		users = append(users, result)
	}

	return users, nil

}
