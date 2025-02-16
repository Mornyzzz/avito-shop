package info

import (
	"avito-shop/internal/entity"
	"avito-shop/internal/repository"
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type UseCase struct {
	repoBalance     BalanceRepo
	repoInventory   InventoryRepo
	repoTransaction TransactionRepo
	trManager       *manager.Manager
}

func New(
	repoBalance *repository.BalanceRepo,
	repoInventory *repository.InventoryRepo,
	repoTransaction *repository.TransactionRepo,
	trManager *manager.Manager,
) *UseCase {
	return &UseCase{
		repoBalance:     repoBalance,
		repoInventory:   repoInventory,
		repoTransaction: repoTransaction,
		trManager:       trManager,
	}
}

//go:generate mockery --name=Inf
//go:generate mockery --name=Info
type (
	Info interface {
		GetInfo(ctx context.Context, username string) (*entity.Info, error)
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

func (uc *UseCase) GetInfo(ctx context.Context, username string) (*entity.Info, error) {
	const op = "usecase.info.GetInfo"

	var (
		balance      int
		inventory    []entity.InventoryItem
		sentTxns     []entity.SentTransaction
		receivedTxns []entity.ReceivedTransaction
		err          error
	)

	err = uc.trManager.Do(ctx, func(ctx context.Context) error {

		balance, err = uc.repoBalance.GetUserBalance(ctx, username)
		if err != nil {
			return err
		}
		inventory, err = uc.repoInventory.GetInventory(ctx, username)
		if err != nil {
			return err
		}
		sentTxns, err = uc.repoTransaction.GetSentTransactions(ctx, username)
		if err != nil {
			return err
		}
		receivedTxns, err = uc.repoTransaction.GetReceivedTransactions(ctx, username)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &entity.Info{
		Coins:     balance,
		Inventory: inventory,
		CoinHistory: entity.CoinHistory{
			Received: receivedTxns,
			Sent:     sentTxns,
		},
	}, nil
}
