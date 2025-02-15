package send

import (
	"avito-shop/internal/entity"
	"avito-shop/internal/repository"
	"context"
	"fmt"
)

type UseCase struct {
	repoBalance     BalanceRepo
	repoTransaction TransactionRepo
}

func New(rb *repository.BalanceRepo, rt *repository.TransactionRepo) *UseCase {
	return &UseCase{
		repoBalance:     rb,
		repoTransaction: rt,
	}
}

//go:generate mockgen -source=auth.go -destination=./mocks_test.go -package=usecase_test

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
		return fmt.Errorf("%s: %w", op, "amount must be greater than zero")
	}

	balance, err := uc.repoBalance.GetUserBalance(ctx, fromUser)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if balance < amount {
		return fmt.Errorf("%s: %w", op, "insufficient funds")
	}

	if err = uc.repoBalance.DecreaseBalance(ctx, fromUser, amount); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = uc.repoBalance.IncreaseBalance(ctx, fromUser, amount); err != nil {
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
}
