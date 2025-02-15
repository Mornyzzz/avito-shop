package repository

import (
	"avito-shop/internal/entity"
	"avito-shop/pkg/postgres"
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

type InventoryRepo struct {
	*postgres.Postgres
}

func NewInventoryRepo(pg *postgres.Postgres) *InventoryRepo {
	return &InventoryRepo{pg}
}

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
		From("inventory").
		Where(sq.Eq{"name": name}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()

	err = rows.Scan(&price)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return price, nil
}

func (r *InventoryRepo) ExistsInventoryItem(ctx context.Context, username string, item string) (bool, error) {
	const op = "repository.inventory.GetInventoryItem"

	query, args, err := sq.Select("1").
		From("inventory").
		Where(sq.Eq{"username": username, "item_name": item}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var exists bool
	err = r.Pool.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: failed to execute query: %w", op, err)
	}

	return exists, nil
}

func (r *InventoryRepo) IncrementInventoryItemQuantity(ctx context.Context, username, item string) error {
	const op = "repository.inventory.IncrementInventoryItemQuantity"

	query, args, err := sq.Update("inventory").
		Set("quantity", sq.Expr("quantity + 1")).
		Where(sq.Eq{"username": username, "item_name": item}).
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

func (r *InventoryRepo) AddInventory(ctx context.Context, inventory entity.Inventory) error {
	const op = "repository.inventory.AddInventoryItem"

	query, args, err := sq.Insert("inventory").
		Columns("username", "item_name", "quantity").
		Values(inventory.Username, inventory.Item, inventory.Quantity).
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

func (r *InventoryRepo) GetInventory(ctx context.Context, username string) ([]entity.InventoryItem, error) {
	const op = "repository.inventory.getInventory"

	query, args, err := sq.Select("name, quantity").
		From("inventory").
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var inventory []entity.InventoryItem

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
