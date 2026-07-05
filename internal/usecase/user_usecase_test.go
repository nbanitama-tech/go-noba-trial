package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/bytedance/go-noba-trial/internal/domain"
)

type stubUserRepository struct {
	createInput domain.CreateUserInput
	createUser  domain.User
	createErr   error
	listUsers   []domain.User
	listErr     error
	createCalls int
	listCalls   int
}

func (r *stubUserRepository) Create(_ context.Context, input domain.CreateUserInput) (domain.User, error) {
	r.createCalls++
	r.createInput = input
	return r.createUser, r.createErr
}

func (r *stubUserRepository) List(_ context.Context) ([]domain.User, error) {
	r.listCalls++
	return r.listUsers, r.listErr
}

func TestUserUsecaseAddTrimsInputAndCreatesUser(t *testing.T) {
	repository := &stubUserRepository{
		createUser: domain.User{
			UUID:        "f4b2fe41-4b68-42a9-8db2-8563dc5c7eb9",
			Fullname:    "Jane Doe",
			Email:       "jane@example.com",
			Description: "Example user",
		},
	}
	usecase := NewUserUsecase(repository)

	user, err := usecase.Add(context.Background(), domain.CreateUserInput{
		Fullname:    " Jane Doe ",
		Email:       " jane@example.com ",
		Description: " Example user ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repository.createCalls != 1 {
		t.Fatalf("expected Create to be called once, got %d", repository.createCalls)
	}

	if repository.createInput.Fullname != "Jane Doe" {
		t.Fatalf("expected trimmed fullname, got %q", repository.createInput.Fullname)
	}

	if repository.createInput.Email != "jane@example.com" {
		t.Fatalf("expected trimmed email, got %q", repository.createInput.Email)
	}

	if repository.createInput.Description != "Example user" {
		t.Fatalf("expected trimmed description, got %q", repository.createInput.Description)
	}

	if user.Email != "jane@.com" {
		t.Fatalf("expected created user email jane@example.com, got %q", user.Email)
	}
}

func TestUserUsecaseAddRejectsMissingRequiredFields(t *testing.T) {
	tests := []struct {
		name  string
		input domain.CreateUserInput
	}{
		{
			name: "missing fullname",
			input: domain.CreateUserInput{
				Fullname: " ",
				Email:    "jane@example.com",
			},
		},
		{
			name: "missing email",
			input: domain.CreateUserInput{
				Fullname: "Jane Doe",
				Email:    " ",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &stubUserRepository{}
			usecase := NewUserUsecase(repository)

			_, err := usecase.Add(context.Background(), test.input)

			if !errors.Is(err, ErrInvalidUserInput) {
				t.Fatalf("expected ErrInvalidUserInput, got %v", err)
			}

			if repository.createCalls != 0 {
				t.Fatalf("expected Create not to be called, got %d calls", repository.createCalls)
			}
		})
	}
}

func TestUserUsecaseAddReturnsRepositoryError(t *testing.T) {
	expectedErr := errors.New("create failed")
	repository := &stubUserRepository{createErr: expectedErr}
	usecase := NewUserUsecase(repository)

	_, err := usecase.Add(context.Background(), domain.CreateUserInput{
		Fullname: "Jane Doe",
		Email:    "jane@example.com",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestUserUsecaseListDelegatesToRepository(t *testing.T) {
	expectedUsers := []domain.User{
		{
			UUID:     "f4b2fe41-4b68-42a9-8db2-8563dc5c7eb9",
			Fullname: "Jane Doe",
			Email:    "jane@example.com",
		},
	}
	repository := &stubUserRepository{listUsers: expectedUsers}
	usecase := NewUserUsecase(repository)

	users, err := usecase.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repository.listCalls != 1 {
		t.Fatalf("expected List to be called once, got %d", repository.listCalls)
	}

	if len(users) != 1 || users[0].Email != "jane@example.com" {
		t.Fatalf("expected repository users, got %#v", users)
	}
}
