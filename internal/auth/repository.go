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
	results, err := r.Db.Query(ctx, "SELECT id, name, role FROM users;")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	var users []User

	for results.Next() {
		user := User{}
		if err := results.Scan(&user.Id, &user.Name, &user.Role); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil

}

func (r *UserRepositoryImpl) Get(ctx context.Context, id int) (User, error) {
	var result User
	err := r.Db.QueryRow(ctx, "SELECT id, name, role FROM users WHERE id = $1", id).Scan(&result.Id, &result.Name, &result.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return result, err
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user UserRequest) (UserCreate, error) {
	result := UserCreate{}
	err := r.Db.QueryRow(ctx, "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, name, role, email", user.Name, user.Email, user.Password).Scan(&result.Id, &result.Name, &result.Role, &result.Email)
	if err != nil {
		return result, err
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
	err := r.Db.QueryRow(ctx, "SELECT id, name, email, password FROM users WHERE email = $1", email).Scan(&result.ID, &result.Name, &result.Email, &result.Password)

	if errors.Is(err, pgx.ErrNoRows) {
		return result, ErrNotFound
	}
	return result, err
}
