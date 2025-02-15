package auth

import (
	"avito-shop/internal/entity"
	e "avito-shop/internal/errors"
	"avito-shop/internal/repository"
	"avito-shop/pkg/jwt"
	"context"
	"errors"
	"fmt"
)

const (
	newUserBalance = 1000
)

type AuthUseCase struct {
	repoUser    UserRepo
	repoBalance BalanceRepo
}

func New(ru *repository.UserRepo, rb *repository.BalanceRepo) *AuthUseCase {
	return &AuthUseCase{
		repoUser:    ru,
		repoBalance: rb,
	}
}

//go:generate mockgen -source=auth.go -destination=./mocks_test.go -package=usecase_test

type (
	Auth interface {
		Authenticate(context.Context, entity.User) (string, error)
	}

	UserRepo interface {
		Get(context.Context, string) (*entity.User, error)
		Add(context.Context, entity.User) error
	}

	BalanceRepo interface {
		InitBalance(ctx context.Context, username string, amount int) error
	}
)

func (uc *AuthUseCase) Authenticate(ctx context.Context, in entity.User) (string, error) {
	const op = "usecase.Authenticate"

	var token string

	user, err := uc.repoUser.Get(ctx, in.Username)
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			token, err = jwt.GenerateToken(in.Username, in.Password)
			if err != nil {
				return "", fmt.Errorf("%s:%w", op, err)
			}
			if err = uc.repoUser.Add(ctx, in); err != nil {
				return "", fmt.Errorf("%s:%w", op, err)
			}
			if err = uc.repoBalance.InitBalance(ctx, in.Username, newUserBalance); err != nil {
				return "", fmt.Errorf("%s:%w", op, err)
			}
			return token, nil
		}
		return "", fmt.Errorf("%s:%w", op, err)
	}

	if in.Password != user.Password {
		return "", fmt.Errorf("%s:%w", op, e.ErrInvalidPassword)
	}

	valid, err := jwt.ValidateToken(in.Username, in.Password, in.Token)
	if err != nil {
		return "", err
	}

	if !valid {
		return "", fmt.Errorf("%s:%w", op, e.ErrInvalidCredentials)
	}

	return in.Token, nil

}
