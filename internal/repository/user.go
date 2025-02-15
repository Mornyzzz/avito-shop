package repository

import (
	"avito-shop/internal/entity"
	e "avito-shop/internal/errors"
	"avito-shop/pkg/postgres"
	"context"
	_ "database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

type User interface {
	Get(ctx context.Context, username string) (*entity.User, error)
	Add(ctx context.Context, user entity.User) error
}

func (r *UserRepo) Get(ctx context.Context, username string) (*entity.User, error) {
	const op = "repository.user.Get"

	query, args, err := sq.Select("name", "password").
		From("users").
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()

	var user entity.User

	err = rows.Scan(&user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s:%w", op, e.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	return &entity.User{}, nil
}

func (r *UserRepo) Add(ctx context.Context, user entity.User) error {
	const op = "repository.user.Add"

	query, args, err := sq.Insert("users").
		Columns("username", "password").
		Values(user.Username, user.Password).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}
	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to execute query: %w", op, err)
	}
	return nil
}
