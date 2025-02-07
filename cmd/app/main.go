package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"restapi_project/internal/database"
	"restapi_project/internal/handlers"
	"restapi_project/internal/taskService"
	"restapi_project/internal/userService"
	"restapi_project/internal/web/tasks"
	"restapi_project/internal/web/users"
)

func main() {
	database.InitDB()
	if err := database.DB.AutoMigrate(&taskService.Task{}); err != nil {
		log.Fatalf("Ошибка миграции БД: %v", err)
		return
	}

	repo := taskService.NewTaskRepository(database.DB)
	service := taskService.NewService(repo)

	repoUser := userService.NewUserRepository(database.DB)
	serviceUser := userService.NewService(repoUser)

	handler := handlers.NewHandler(service)
	handlerUser := handlers.NewUserHandler(serviceUser)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	strictHandler := tasks.NewStrictHandler(handler, nil)
	tasks.RegisterHandlers(e, strictHandler)

	strictUserHandler := users.NewStrictHandler(handlerUser, nil)
	users.RegisterHandlers(e, strictUserHandler)

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start with err: %v", err)
	}
}
