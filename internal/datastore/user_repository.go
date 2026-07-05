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
		INSERT INTO "users" (fullname, email, description)
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

func (r *UserRepository) List(ctx context.Context) ([]domain.User, error) {
	const query = `
		SELECT uuid, fullname, email, description
		FROM "users"
		ORDER BY fullname ASC, email ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.UUID,
			&user.Fullname,
			&user.Email,
			&user.Description,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
