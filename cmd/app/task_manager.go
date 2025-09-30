package app

import (
	"fmt"
	"log"
	"task_manager_server/internal/application/usecase"
	"task_manager_server/internal/config"
	"task_manager_server/internal/domain/repository"
	"task_manager_server/internal/infrastructure/handlers"
	"task_manager_server/internal/infrastructure/persistence/database"
	implrepository "task_manager_server/internal/infrastructure/persistence/repository"
)

type App struct {
	Config   *config.Config
	Repos    *Repositories
	UseCases *UseCases
	Handlers *Handlers
}

type Repositories struct {
	UserRepo repository.UserRepository
	TaskRepo repository.TaskRepository
}

type UseCases struct {
	UserUseCase *usecase.UserUseCase
	TaskUseCase *usecase.TaskUseCase
}

type Handlers struct {
	UserHandler *handlers.UserHandler
	TaskHandler *handlers.TaskHandler
}

func NewApp() *App {
	conf, err := config.NewConfig()

	if err != nil {
		log.Fatalln(err)
	}

	conn, err := database.NewPostgresConnection(conf.DB.DSN())

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Успешное подключение к БД")

	repos := &Repositories{
		UserRepo: implrepository.NewUserRepository(conn),
		TaskRepo: implrepository.NewTaskRepository(conn),
	}

	uc := &UseCases{
		UserUseCase: usecase.NewUserUseCase(repos.UserRepo),
		TaskUseCase: usecase.NewTaskUseCase(repos.TaskRepo),
	}

	h := &Handlers{
		UserHandler: handlers.NewUserHandler(uc.UserUseCase),
		TaskHandler: handlers.NewTaskHandler(uc.TaskUseCase),
	}

	return &App{
		Config:   conf,
		Repos:    repos,
		UseCases: uc,
		Handlers: h,
	}
}
