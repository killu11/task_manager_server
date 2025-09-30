package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"task_manager_server/internal/domain/entites"
	_ "task_manager_server/internal/domain/entites"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(
	ctx context.Context,
	user *entites.User,
) (int, error) {
	var id int

	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO users("username", "hash_password")
				VALUES ($1, $2) RETURNING id`,
		user.Name(),
		user.Hash(),
	).Scan(&id)

	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	if err != nil {
		if isUniqueViolation(err, "username") {
			return 0, entites.NameAlreadyUsingErr
		}

		return 0, fmt.Errorf("%w: database create failed: %v", RepoErr, err)
	}
	return id, nil
}

func (r *UserRepository) GetUserIDAndHash(
	ctx context.Context,
	username string,
) (int, []byte, error) {
	var id int
	hash := make([]byte, 0)
	
	err := r.db.QueryRowContext(ctx,
		`SELECT id, hash_password FROM users
          WHERE username = $1`,
		username).Scan(&id, &hash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil, entites.UserNotExistErr // заменить позже на ошибку-маркер
		}
		return 0, nil, fmt.Errorf("%w: database find user failed: %v", RepoErr, err)
	}
	return id, hash, nil
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" && strings.Contains(pgErr.Message, constraint)
	}
	return false
}
