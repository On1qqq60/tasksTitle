package main

import (
	"log"
	"net/http"
	"taskTitle/internal/tasks/handler"
	"taskTitle/internal/tasks/repository"
	"taskTitle/internal/tasks/service"
	usershandler "taskTitle/internal/users/handler"
	usersrepo "taskTitle/internal/users/repository"
	usersservice "taskTitle/internal/users/service"
)

func main() {
	db, err := repository.OpenDB()
	if err != nil {
		log.Fatalf("Ошибка старта бд: %s", err)
	}
	defer db.Close()

	if err := repository.InitSchemaTask(db); err != nil {
		log.Fatalf("Ошибка схемы: %s", err)
	}

	if err := usersrepo.InitSgemaUsers(db); err != nil {
		log.Fatalf("Ошибка схемы пользователей: %s", err)
	}

	//tasks
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	//user
	usersRepository := usersrepo.NewRepository(db)
	usersSvc := usersservice.NewService(usersRepository)
	usersH := usershandler.NewHandler(usersSvc)

	mux := http.NewServeMux()

	//staf
	mux.HandleFunc("GET /", h.Index)
	mux.HandleFunc("GET /health", h.Health)

	//tasks
	mux.HandleFunc("GET /api/notes", h.ListHandler)
	mux.HandleFunc("POST /api/notes", h.AddHandler)
	mux.HandleFunc("POST /delete", h.DeleteHandler)
	mux.HandleFunc("POST /add", h.AddPageHandler)

	//user
	mux.HandleFunc("POST /register", usersH.RegistrHandler)
	mux.HandleFunc("POST /login", usersH.LoginHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Ошибка: %s", err)
	}
}
