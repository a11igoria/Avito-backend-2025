package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"avito-2025/internal/api/handlers"
	"avito-2025/internal/service"
	"avito-2025/internal/storage"

	_ "github.com/lib/pq"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "avito_db"
	}

	// Строка подключения
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	log.Println("Подключение к БД успешно!")

	userRepo := storage.NewUserRepository(db)
	teamRepo := storage.NewTeamRepository(db)
	prRepo := storage.NewPRRepository(db)
	prReviewerRepo := storage.NewPRReviewerRepository(db)

	userService := service.NewUserService(userRepo, teamRepo)
	teamService := service.NewTeamService(teamRepo, userRepo)
	prService := service.NewPRService(prRepo, prReviewerRepo, userRepo)

	userHandler := handlers.NewUserHandler(userService)
	teamHandler := handlers.NewTeamHandler(teamService)
	prHandler := handlers.NewPRHandler(prService)

	mux := http.NewServeMux()

	// User endpoints
	mux.HandleFunc("POST /api/users", userHandler.CreateUser)
	mux.HandleFunc("GET /api/users/{id}", userHandler.GetUser)
	mux.HandleFunc("GET /api/teams/{teamName}/members", userHandler.GetTeamMembers)
	mux.HandleFunc("PATCH /api/users/{id}/activate", userHandler.ActivateUser)
	mux.HandleFunc("PATCH /api/users/{id}/deactivate", userHandler.DeactivateUser)

	// Team endpoints
	mux.HandleFunc("POST /api/teams", teamHandler.CreateTeam)
	mux.HandleFunc("GET /api/teams/{teamName}", teamHandler.GetTeamByName)
	mux.HandleFunc("GET /api/teams", teamHandler.ListTeams)

	// PR endpoints
	mux.HandleFunc("POST /api/pull-requests", prHandler.CreatePR)
	mux.HandleFunc("GET /api/pull-requests/{id}", prHandler.GetPR)
	mux.HandleFunc("PATCH /api/pull-requests/{id}/status", prHandler.UpdatePRStatus)
	mux.HandleFunc("POST /api/pull-requests/{id}/reviewers", prHandler.AssignReviewer)
	mux.HandleFunc("DELETE /api/pull-requests/{id}/reviewers/{reviewerId}", prHandler.RemoveReviewer)
	mux.HandleFunc("GET /api/users/{userId}/reviews", prHandler.GetPRsWhereUserIsReviewer)

	port := ":8080"
	log.Printf("Сервер запущен на http://localhost%s", port)
	log.Printf("Документация: http://localhost%s/swagger", port)

	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
