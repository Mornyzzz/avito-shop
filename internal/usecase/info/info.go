package info

import (
	"avito-shop/internal/entity"
	"avito-shop/internal/repository"
	"context"
)

type UseCase struct {
	repoBalance     BalanceRepo
	repoInventory   InventoryRepo
	repoTransaction TransactionRepo
}

func New(
	repoBalance *repository.BalanceRepo,
	repoInventory *repository.InventoryRepo,
	repoTransaction *repository.TransactionRepo,
) *UseCase {
	return &UseCase{
		repoBalance:     repoBalance,
		repoInventory:   repoInventory,
		repoTransaction: repoTransaction,
	}
}

type (
	Info interface {
		GetInfo(ctx context.Context, username string) (*entity.InfoResponse, error)
	}

	BalanceRepo interface {
		GetUserBalance(ctx context.Context, username string) (int, error)
		DecreaseBalance(ctx context.Context, username string, amount int) error
	}

	InventoryRepo interface {
		GetInventory(ctx context.Context, username string) ([]entity.InventoryItem, error)
	}

	TransactionRepo interface {
		GetSentTransactions(ctx context.Context, username string) ([]entity.SentTransaction, error)
		GetReceivedTransactions(ctx context.Context, username string) ([]entity.ReceivedTransaction, error)
	}
)

func (uc *UseCase) GetInfo(ctx context.Context, username string) (*entity.InfoResponse, error) {
	const op = "usecase.info.GetInfo"

	balance, err := uc.repoBalance.GetUserBalance(ctx, username)
	if err != nil {
		return nil, err
	}

	inventory, err := uc.repoInventory.GetInventory(ctx, username)
	if err != nil {
		return nil, err
	}

	sentTxns, err := uc.repoTransaction.GetSentTransactions(ctx, username)
	if err != nil {
		return nil, err
	}

	receivedTxns, err := uc.repoTransaction.GetReceivedTransactions(ctx, username)
	if err != nil {
		return nil, err
	}

	return &entity.InfoResponse{
		Coins:     balance,
		Inventory: inventory,
		CoinHistory: entity.CoinHistory{
			Received: receivedTxns,
			Sent:     sentTxns,
		},
	}, nil
}
