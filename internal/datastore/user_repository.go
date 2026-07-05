package datastore

import (
	"context"
	"database/sql"

	"github.com/bytedance/go-noba-trial/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, input domain.CreateUserInput) (domain.User, error) {
	const query = `
		INSERT INTO "user" (fullname, email, description)
		VALUES ($1, $2, $3)
		RETURNING uuid, fullname, email, description
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, input.Fullname, input.Email, input.Description).Scan(
		&user.UUID,
		&user.Fullname,
		&user.Email,
		&user.Description,
	)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
