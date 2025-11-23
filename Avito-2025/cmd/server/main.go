// Создаёшь БД
db, err := storage.NewDB("postgres://...")
if err != nil {
	panic(err)
}

// Создаёшь репозитории
userRepo := storage.NewUserRepository(db)
teamRepo := storage.NewTeamRepository(db)
prRepo := storage.NewPRRepository(db)
prReviewerRepo := storage.NewPRReviewerRepository(db)

// Создаёшь сервисы
userService := service.NewUserService(userRepo, teamRepo)
teamService := service.NewTeamService(teamRepo, userRepo)
prService := service.NewPRService(prRepo, prReviewerRepo, userRepo)

// Используешь сервисы (потом в хендлерах)
user, err := userService.CreateUser(ctx, "john", 1)
