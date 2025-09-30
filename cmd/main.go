package main

import (
	"log"
	"net/http"
	service "task_manager_server/cmd/app"
	_ "task_manager_server/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Task Manager API
// @version 1.0
// @description Task Manager REST API
// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	app := service.NewApp()
	r := mux.NewRouter()
	app.Handlers.UserHandler.InitRoutes(r)
	app.Handlers.TaskHandler.InitRoutes(r)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Fatalln(http.ListenAndServe(":8080", r))

}
