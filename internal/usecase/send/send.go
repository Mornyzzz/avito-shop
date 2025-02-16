package send

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"avito-shop/internal/entity"
	"avito-shop/internal/repository"
	"avito-shop/pkg/errors"
)

type UseCase struct {
	repoBalance     BalanceRepo
	repoTransaction TransactionRepo
	trManager       *manager.Manager
}

func New(rb *repository.BalanceRepo, rt *repository.TransactionRepo, trManager *manager.Manager) *UseCase {
	return &UseCase{
		repoBalance:     rb,
		repoTransaction: rt,
		trManager:       trManager,
	}
}

//go:generate mockery --name=Send

type (
	Send interface {
		SendCoin(ctx context.Context, fromUser, toUser string, amount int) error
	}

	BalanceRepo interface {
		GetUserBalance(ctx context.Context, username string) (int, error)
		DecreaseBalance(ctx context.Context, username string, amount int) error
		IncreaseBalance(ctx context.Context, username string, amount int) error
	}

	TransactionRepo interface {
		AddTransaction(ctx context.Context, txn entity.CoinTransaction) error
	}
)

func (uc *UseCase) SendCoin(ctx context.Context, fromUser, toUser string, amount int) error {
	const op = "usecase.SendCoin"

	if amount <= 0 {
		return fmt.Errorf("%s: %w", op, errors.ErrInvalidCredentials)
	}

	err := uc.trManager.Do(ctx, func(ctx context.Context) error {
		balance, err := uc.repoBalance.GetUserBalance(ctx, fromUser)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		_, err = uc.repoBalance.GetUserBalance(ctx, toUser)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if balance < amount {
			return fmt.Errorf("%s: %w", op, errors.ErrInvalidCredentials)
		}

		if err = uc.repoBalance.DecreaseBalance(ctx, fromUser, amount); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if err = uc.repoBalance.IncreaseBalance(ctx, toUser, amount); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		txn := entity.CoinTransaction{
			FromUser: fromUser,
			ToUser:   toUser,
			Amount:   amount,
		}

		if err = uc.repoTransaction.AddTransaction(ctx, txn); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
