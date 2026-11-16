package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	List(ctx context.Context) ([]User, error)
	Get(ctx context.Context, id int) (User, error)
	Create(ctx context.Context, user User) (User, error)
}

type UserRepositoryImpl struct {
	Db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &UserRepositoryImpl{Db: db}
}

var ErrNotFound error = errors.New("User not found")

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

func (r *UserRepositoryImpl) Get(ctx context.Context, id int) (User, error) {
	tx, err := r.Db.Begin(ctx)
	result := User{}
	if err != nil {
		return result, err
	}

	qErr := tx.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", id).Scan(&result.Id, &result.Name, &result.Role)
	if errors.Is(pgx.ErrNoRows, qErr) {
		return result, ErrNotFound
	}

	if qErr != nil {
		return result, qErr
	}

	if comerr := tx.Commit(ctx); comerr != nil {
		return result, comerr
	}

	return result, nil
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user User) (User, error) {
	tx, err := r.Db.Begin(ctx)
	result := User{}
	if err != nil {
		return result, err
	}

	qErr := tx.QueryRow(ctx, "INSERT INTO users (name, role) VALUES ($1, $2) RETURNING id", user.Name, user.Role).Scan(&result.Id, result.Name, result.Role)
	if qErr != nil {
		return result, qErr
	}

	if comerr := tx.Commit(ctx); comerr != nil {
		return result, comerr
	}

	return result, nil
}
