package repository

import (
	"avito-shop/pkg/postgres"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

type BalanceRepo struct {
	*postgres.Postgres
}

func NewBalanceRepo(pg *postgres.Postgres) *BalanceRepo {
	return &BalanceRepo{pg}
}

type Balance interface {
	InitBalance(ctx context.Context, username string) error
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

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	return nil
}

func (r *BalanceRepo) GetUserBalance(ctx context.Context, username string) (int, error) {
	const op = "repository.balance.GetUserBalance"

	var balance int
	query, _, err := sq.Select("coins").
		From("balances").
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()

	err = rows.Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return balance, nil
}

func (r *BalanceRepo) DecreaseBalance(ctx context.Context, username string, amount int) error {
	const op = "repository.balance.DecreaseBalance"

	query, args, err := sq.Update("balances").
		Set("coins", sq.Expr("coins - ?", amount)).
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *BalanceRepo) IncreaseBalance(ctx context.Context, username string, amount int) error {
	const op = "repository.balance.IncreaseBalance"

	query, args, err := sq.Update("balances").
		Set("coins", sq.Expr("coins + ?", amount)).
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
