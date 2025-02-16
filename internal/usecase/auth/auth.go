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

type UseCase struct {
	repoUser    UserRepo
	repoBalance BalanceRepo
}

func New(ru *repository.UserRepo, rb *repository.BalanceRepo) *UseCase {
	return &UseCase{
		repoUser:    ru,
		repoBalance: rb,
	}
}

//go:generate mockgen -source=auth.go -destination=./mocks_test.go -package=usecase_test

type (
	Auth interface {
		Login(context.Context, entity.User) (string, error)
		Register(context.Context, entity.User) (string, error)
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

	// Начинаем транзакцию
	tx, err := uc.repoUser.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}

	// Откат транзакции в случае ошибки
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("%s: failed to rollback transaction: %v", op, rollbackErr)
			}
		}
	}()

	// Получаем пользователя в рамках транзакции
	user, err := uc.repoUser.Get(ctx, in.Username)
	if errors.Is(err, e.ErrNotFound) {
		// Если пользователь не найден, регистрируем его в рамках той же транзакции
		token, err := uc.Register(ctx, in)
		if err != nil {
			return "", fmt.Errorf("%s: failed to register user: %w", op, err)
		}

		// Фиксируем транзакцию, если регистрация прошла успешно
		if commitErr := tx.Commit(); commitErr != nil {
			return "", fmt.Errorf("%s: failed to commit transaction: %w", op, commitErr)
		}

		return token, nil
	}

	if err != nil {
		return "", fmt.Errorf("%s: failed to get user: %w", op, err)
	}

	// Проверяем пароль
	if user.Password != in.Password {
		return "", fmt.Errorf("%s: %w", op, e.ErrInvalidPassword)
	}

	// Генерируем токен
	token, err := jwt.GenerateToken(in.Username)
	if err != nil {
		return "", fmt.Errorf("%s: failed to generate token: %w", op, err)
	}

	// Фиксируем транзакцию, если все прошло успешно
	if commitErr := tx.Commit(); commitErr != nil {
		return "", fmt.Errorf("%s: failed to commit transaction: %w", op, commitErr)
	}

	return token, nil
}

func (uc *UseCase) Register(ctx context.Context, in entity.User) (string, error) {
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
