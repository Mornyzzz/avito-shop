package repository

import (
	"avito-shop/internal/entity"
	"avito-shop/pkg/postgres"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

type TransactionRepo struct {
	*postgres.Postgres
}

type Transaction interface {
	AddTransaction(ctx context.Context, txn entity.CoinTransaction) error
	GetReceivedTransactions(ctx context.Context, username string) ([]entity.ReceivedTransaction, error)
	GetSentTransactions(ctx context.Context, username string) ([]entity.SentTransaction, error)
}

func NewTransactionRepo(pg *postgres.Postgres) *TransactionRepo {
	return &TransactionRepo{pg}
}

func (r *TransactionRepo) AddTransaction(ctx context.Context, txn entity.CoinTransaction) error {
	const op = "repository.transaction.AddTransaction"

	query, args, err := sq.Insert("coinTransaction").
		Columns("fromUser", "toUser", "amount").
		Values(txn.FromUser, txn.ToUser, txn.Amount).
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

func (r *TransactionRepo) GetReceivedTransactions(ctx context.Context, username string) ([]entity.ReceivedTransaction, error) {
	const op = "repository.transaction.GetReceivedTransactions"

	query, args, err := sq.Select("fromUser", "amount").
		From("coinTransaction").
		Where(sq.Eq{"toUser": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}
	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}
	defer rows.Close()

	var receivedTxns []entity.ReceivedTransaction

	for rows.Next() {
		var item entity.ReceivedTransaction
		if err = rows.Scan(&item.FromUser, &item.Amount); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		receivedTxns = append(receivedTxns, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return receivedTxns, nil
}

func (r *TransactionRepo) GetSentTransactions(ctx context.Context, username string) ([]entity.SentTransaction, error) {
	const op = "repository.transaction.GetSentTransactions"

	query, args, err := sq.Select("toUser", "amount").
		From("Cointransaction").
		Where(sq.Eq{"fromUser": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}
	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}
	defer rows.Close()

	var sentTxns []entity.SentTransaction

	for rows.Next() {
		var item entity.SentTransaction
		if err = rows.Scan(&item.ToUser, &item.Amount); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		sentTxns = append(sentTxns, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sentTxns, nil
}
