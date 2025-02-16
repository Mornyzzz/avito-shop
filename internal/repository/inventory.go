package repository

import (
	"avito-shop/internal/entity"
	e "avito-shop/pkg/errors"
	"avito-shop/pkg/postgres"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv4/v2"
)

type InventoryRepo struct {
	*postgres.Postgres
}

func NewInventoryRepo(pg *postgres.Postgres) *InventoryRepo {
	return &InventoryRepo{pg}
}

//go:generate mockery --name=Inventory

type Inventory interface {
	GetItemPrice(ctx context.Context, name string) (int, error)
	ExistsInventoryItem(ctx context.Context, username string, item string) (bool, error)
	IncrementInventoryItemQuantity(ctx context.Context, username, item string) error
	AddInventory(ctx context.Context, inventory entity.Inventory) error
	GetInventory(ctx context.Context, username string) ([]entity.InventoryItem, error)
}

func (r *InventoryRepo) GetItemPrice(ctx context.Context, name string) (int, error) {
	const op = "repository.inventory.GetItemPrice"

	var price int
	query, args, err := sq.Select("price").
		From("item").
		Where(sq.Eq{"name": name}).
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

	err = rows.Scan(&price)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if rows.Next() {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return price, nil
}

func (r *InventoryRepo) ExistsInventoryItem(ctx context.Context, username, item string) (bool, error) {
	const op = "repository.inventory.ExistsInventoryItem"

	query, _, err := sq.Select("EXISTS(SELECT 1 FROM inventory WHERE username = $1 AND item = $2)").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var exists bool
	conn := trmpgx.DefaultCtxGetter.DefaultTrOrDB(ctx, r.Pool)
	err = conn.QueryRow(ctx, query, username, item).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	return exists, nil
}

func (r *InventoryRepo) IncrementInventoryItemQuantity(ctx context.Context, username, item string) error {
	const op = "repository.inventory.IncrementInventoryItemQuantity"

	query, args, err := sq.Update("inventory").
		Set("quantity", sq.Expr("quantity + 1")).
		Where(sq.Eq{"username": username, "item": item}).
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

func (r *InventoryRepo) AddInventory(ctx context.Context, inventory entity.Inventory) error {
	const op = "repository.inventory.AddInventory"

	query, args, err := sq.Insert("inventory").
		Columns("username", "item", "quantity").
		Values(inventory.Username, inventory.Item, inventory.Quantity).
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

func (r *InventoryRepo) GetInventory(ctx context.Context, username string) ([]entity.InventoryItem, error) {
	const op = "repository.inventory.getInventory"

	query, args, err := sq.Select("item, quantity").
		From("inventory").
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	conn := trmpgx.DefaultCtxGetter.DefaultTrOrDB(ctx, r.Pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var inventory []entity.InventoryItem

	if !rows.Next() {
		return nil, fmt.Errorf("%s: %w", op, e.ErrNotFound)
	}

	for rows.Next() {
		var item entity.InventoryItem
		if err = rows.Scan(&item.Name, &item.Quantity); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		inventory = append(inventory, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return inventory, nil
}
