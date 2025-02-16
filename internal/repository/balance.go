package repository

import (
	e "avito-shop/pkg/errors"
	"avito-shop/pkg/postgres"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
)

type BalanceRepo struct {
	*postgres.Postgres
}

func NewBalanceRepo(pg *postgres.Postgres) *BalanceRepo {
	return &BalanceRepo{pg}
}

//go:generate mockery --name=Balance

type Balance interface {
	InitBalance(ctx context.Context, username string, amount int) error
	GetUserBalance(ctx context.Context, username string) (int, error)
	DecreaseBalance(ctx context.Context, username string, amount int) error
	IncreaseBalance(ctx context.Context, username string, amount int) error
}

func (r *BalanceRepo) InitBalance(ctx context.Context, username string, amount int) error {
	const op = "repository.balance.InitBalance"

	query, args, err := sq.Insert("balance").
		Columns("username", "coins").
		Values(username, amount).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}
	conn := trmpgx.DefaultCtxGetter.DefaultTrOrDB(ctx, r.Pool)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	return nil
}

func (r *BalanceRepo) GetUserBalance(ctx context.Context, username string) (int, error) {
	const op = "repository.balance.GetUserBalance"

	var balance int
	query, args, err := sq.Select("coins").
		From("balance").
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	conn := trmpgx.DefaultCtxGetter.DefaultTrOrDB(ctx, r.Pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, fmt.Errorf("%s: %w", op, e.ErrNotFound)
	}

	err = rows.Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if rows.Next() {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return balance, nil
}

func (r *BalanceRepo) DecreaseBalance(ctx context.Context, username string, amount int) error {
	const op = "repository.balance.DecreaseBalance"

	query, args, err := sq.Update("balance").
		Set("coins", sq.Expr("coins - ?", amount)).
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	conn := trmpgx.DefaultCtxGetter.DefaultTrOrDB(ctx, r.Pool)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *BalanceRepo) IncreaseBalance(ctx context.Context, username string, amount int) error {
	const op = "repository.balance.IncreaseBalance"

	query, args, err := sq.Update("balance").
		Set("coins", sq.Expr("coins + ?", amount)).
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	conn := trmpgx.DefaultCtxGetter.DefaultTrOrDB(ctx, r.Pool)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
