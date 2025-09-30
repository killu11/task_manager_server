package usecase

import (
	"context"
	"errors"
	"log"
	"regexp"
	"task_manager_server/internal/domain/entites"
	"task_manager_server/internal/domain/repository"
	implrepository "task_manager_server/internal/infrastructure/persistence/repository"
	"task_manager_server/pkg/security"
)

var PasswordRegExp = regexp.MustCompile(`^[A-Za-z0-9]{7,}$`)

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (c *UserUseCase) Register(ctx context.Context, username, password string) (string, error) {
	// Блок подготовки данных (валидация + хеширование)
	if username == "" {
		return "", entites.InvalidLoginErr
	}

	if !PasswordRegExp.MatchString(password) {
		return "", entites.InvalidPasswordErr
	}

	hp, err := security.HashPassword(password)

	if err != nil {
		log.Printf("hash password err: %v", err)
		return "", security.HashPassErr
	}

	// Создание и сохранение пользователя в БД
	user := entites.NewUser(username, hp)
	id, err := c.repo.Save(ctx, user)

	if err != nil {
		log.Printf("save user failed: %v", err)
		return "", err
	}
	// Генерация access-токена
	token, err := security.GenerateJWT(id, username)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return token, nil
}

func (c *UserUseCase) Auth(
	ctx context.Context,
	username string,
	password string,
) (string, error) {
	id, hash, err := c.repo.GetUserIDAndHash(ctx, username)

	if err != nil {
		if errors.Is(err, implrepository.RepoErr) {
			log.Println(err)
			return "", implrepository.RepoErr
		}

		if errors.Is(err, security.HashPassErr) {
			log.Println(err)
			return "", security.HashPassErr
		}
		return "", entites.UserNotExistErr
	}

	if !security.VerifyPassword(hash, password) {
		return "", security.VerifyPassErr
	}

	token, err := security.GenerateJWT(id, username)

	if err != nil {
		log.Printf("%v: %v", security.GenerateJWTErr, err)
		return "", security.GenerateJWTErr
	}
	return token, nil
}
