package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"avito-shop/internal/entity"
	"avito-shop/internal/repository"
	e "avito-shop/pkg/errors"
	"avito-shop/pkg/jwt"
)

const (
	newUserBalance = 1000
)

type UseCase struct {
	repoUser    UserRepo
	repoBalance BalanceRepo
	trManager   *manager.Manager
}

func New(ru *repository.UserRepo, rb *repository.BalanceRepo, trManager *manager.Manager) *UseCase {
	return &UseCase{
		repoUser:    ru,
		repoBalance: rb,
		trManager:   trManager,
	}
}

//go:generate mockery --name=Auth

type (
	Auth interface {
		Login(context.Context, entity.User) (string, error)
	}

	UserRepo interface {
		Get(context.Context, string) (*entity.User, error)
		Add(context.Context, entity.User) error
	}

	BalanceRepo interface {
		InitBalance(ctx context.Context, username string, amount int) error
	}
)

func (uc *UseCase) Login(ctx context.Context, in entity.User) (string, error) {
	const op = "usecase.auth.Login"

	var token string

	var err error

	var user *entity.User

	err = uc.trManager.Do(ctx, func(ctx context.Context) error {
		user, err = uc.repoUser.Get(ctx, in.Username)
		if errors.Is(err, e.ErrNotFound) {
			token, err = uc.register(ctx, in)
			if err != nil {
				return fmt.Errorf("%s: failed to register user: %w", op, err)
			}

			return nil
		} else if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("%s: failed to login: %w", op, err)
	}

	if token != "" {
		return token, nil
	}

	if user.Password != in.Password {
		return "", fmt.Errorf("%s: %w", op, e.ErrInvalidCredentials)
	}

	token, err = jwt.GenerateToken(in.Username)
	if err != nil {
		return "", fmt.Errorf("%s: failed to generate token: %w", op, err)
	}

	return token, nil
}

func (uc *UseCase) register(ctx context.Context, in entity.User) (string, error) {
	const op = "usecase.auth.Register"

	if err := uc.repoUser.Add(ctx, in); err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}

	if err := uc.repoBalance.InitBalance(ctx, in.Username, newUserBalance); err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}

	token, err := jwt.GenerateToken(in.Username)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}

	return token, nil
}
