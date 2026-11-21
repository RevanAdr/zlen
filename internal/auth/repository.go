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
	Create(ctx context.Context, user UserRequest) (UserCreate, error)
	Delete(ctx context.Context, id int) error
	FindByEmail(ctx context.Context, email string) (LoginResponse, error)
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

	results, qErr := tx.Query(ctx, "SELECT id, name, role FROM users;")

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

	qErr := tx.QueryRow(ctx, "SELECT id, name, role FROM users WHERE id = $1", id).Scan(&result.Id, &result.Name, &result.Role)
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

func (r *UserRepositoryImpl) Create(ctx context.Context, user UserRequest) (UserCreate, error) {
	tx, err := r.Db.Begin(ctx)
	result := UserCreate{}
	if err != nil {
		return result, err
	}

	qErr := tx.QueryRow(ctx, "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, name, role, email", user.Name, user.Email, user.Password).Scan(&result.Id, &result.Name, &result.Role, &result.Email)
	if qErr != nil {
		tx.Rollback(ctx)
		return result, qErr
	}

	if comerr := tx.Commit(ctx); comerr != nil {
		return result, comerr
	}

	return result, nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id int) error {
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return err
	}

	_, rErr := tx.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if rErr != nil {
		tx.Rollback(ctx)
		return rErr
	}

	if comerr := tx.Commit(ctx); comerr != nil {
		tx.Rollback(ctx)
		return comerr
	}
	return nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (LoginResponse, error) {
	result := LoginResponse{}
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return LoginResponse{}, err
	}

	responseErr := tx.QueryRow(ctx, "SELECT name, email, password FROM users WHERE email = $1", email).Scan(&result.Name, &result.Email, &result.Password)
	switch {
	case errors.Is(responseErr, pgx.ErrNoRows):
		return result, ErrNotFound
	case responseErr != nil:
		return result, responseErr
	}

	if cErr := tx.Commit(ctx); cErr != nil {
		tx.Rollback(ctx)
		return result, cErr
	}

	return result, nil
}
