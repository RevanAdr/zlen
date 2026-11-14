package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	List(ctx context.Context) ([]User, error)
}

type UserRepositoryImpl struct {
	Db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &UserRepositoryImpl{Db: db}
}

func (r *UserRepositoryImpl) List(ctx context.Context) ([]User, error) {
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
		user := User{}
		if err := results.Scan(&user.Id, &user.Name, &user.Role); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	if comerr := tx.Commit(ctx); comerr != nil {
		return users, comerr
	}

	return users, nil

}
