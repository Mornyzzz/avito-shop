package buy

import (
	"avito-shop/internal/entity"
	"avito-shop/internal/repository"
	"context"
	"fmt"
)

type UseCase struct {
	repoBalance   BalanceRepo
	repoInventory InventoryRepo
}

func New(rB *repository.BalanceRepo, rI *repository.InventoryRepo) *UseCase {
	return &UseCase{
		repoBalance:   rB,
		repoInventory: rI,
	}
}

//go:generate mockgen -source=buy.go -destination=./mocks_test.go -package=usecase_test

type (
	Buy interface {
		BuyItem(ctx context.Context, username, item string) error
	}

	BalanceRepo interface {
		GetUserBalance(ctx context.Context, username string) (int, error)
		DecreaseBalance(ctx context.Context, username string, amount int) error
	}

	InventoryRepo interface {
		GetItemPrice(ctx context.Context, name string) (int, error)
		AddInventory(ctx context.Context, inventory entity.Inventory) error
		ExistsInventoryItem(ctx context.Context, username string, item string) (bool, error)
		IncrementInventoryItemQuantity(ctx context.Context, username string, item string) error
	}
)

func (uc *UseCase) BuyItem(ctx context.Context, username, item string) error {
	const op = "usecase.BuyItem"

	price, err := uc.repoInventory.GetItemPrice(ctx, item)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	balance, err := uc.repoBalance.GetUserBalance(ctx, username)
	if err != nil {
		return err
	}

	if balance < price {
		return fmt.Errorf("%s: %w", op, "NO BALANCE")
	}

	if err = uc.repoBalance.DecreaseBalance(ctx, username, price); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	exists, err := uc.repoInventory.ExistsInventoryItem(ctx, username, item)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !exists {
		err = uc.repoInventory.AddInventory(ctx, entity.Inventory{
			Username: username,
			Item:     item,
			Quantity: 0})
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if err = uc.repoInventory.IncrementInventoryItemQuantity(ctx, username, item); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
