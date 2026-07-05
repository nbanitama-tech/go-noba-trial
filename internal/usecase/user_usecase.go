package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/bytedance/go-noba-trial/internal/domain"
)

var ErrInvalidUserInput = errors.New("invalid user input")

type UserRepository interface {
	Create(ctx context.Context, input domain.CreateUserInput) (domain.User, error)
}

type UserUsecase interface {
	Add(ctx context.Context, input domain.CreateUserInput) (domain.User, error)
}

type userUsecase struct {
	repository UserRepository
}

func NewUserUsecase(repository UserRepository) UserUsecase {
	return &userUsecase{
		repository: repository,
	}
}

func (u *userUsecase) Add(ctx context.Context, input domain.CreateUserInput) (domain.User, error) {
	input.Fullname = strings.TrimSpace(input.Fullname)
	input.Email = strings.TrimSpace(input.Email)
	input.Description = strings.TrimSpace(input.Description)

	if input.Fullname == "" || input.Email == "" {
		return domain.User{}, ErrInvalidUserInput
	}

	return u.repository.Create(ctx, input)
}
