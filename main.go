package main

import (
	"log"
	"net/http"
	"time"

	"task-manager/data"
	"task-manager/dto"
	"task-manager/handlers"
	"task-manager/middlewares"
	"task-manager/models"
	"task-manager/repositories"
	"task-manager/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) *gin.Engine {
	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService)

	assignmentRepository := repositories.NewAssignmentRepository(db)
	assignmentService := services.NewAssignmentService(assignmentRepository)
	assignmentHandler := handlers.NewAssignmentHandler(assignmentService)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:8080", "http://localhost:8080", "http://127.0.0.1:5173", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	authRouter := router.Group("/auth")
	authRouter.POST("/signup", authHandler.SignUp)
	authRouter.POST("/login", authHandler.Login)

	assignmentRouter := router.Group("/assignments", middlewares.AuthMiddleware(authService))
	assignmentRouter.GET("", assignmentHandler.FindAll)
	assignmentRouter.GET("/:id", assignmentHandler.FindByID)
	assignmentRouter.POST("", assignmentHandler.Create)
	assignmentRouter.PUT("/:id", assignmentHandler.Update)
	assignmentRouter.DELETE("/:id", assignmentHandler.Delete)

	router.Static("/web", "./web")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/index.html")
	})

	return router
}

func main() {
	_ = godotenv.Load()

	db := data.SetupDB()

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	assignmentRepository := repositories.NewAssignmentRepository(db)
	assignmentService := services.NewAssignmentService(assignmentRepository)

	seedDemoUser(authService, assignmentService)

	router := setupRouter(db)
	log.Println("server started at http://127.0.0.1:8080")
	log.Println("frontend: http://127.0.0.1:8080/web/index.html")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func seedDemoUser(authService services.AuthService, assignmentService services.AssignmentService) {
	_ = authService.SignUp("demo@example.com", "password123")

	assignments, err := assignmentService.FindAll(1)
	if err != nil || len(assignments) > 0 {
		return
	}

	dueTomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	dueNextWeek := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	_, _ = assignmentService.Create(1, dto.CreateAssignmentInput{
		Title:       "Goの課題を提出する",
		Description: "task-manager APIの実装を完了する",
		DueDate:     &dueTomorrow,
		Priority:    models.PriorityHigh,
		Status:      models.StatusTodo,
	})
	_, _ = assignmentService.Create(1, dto.CreateAssignmentInput{
		Title:       "買い物リスト",
		Description: "牛乳とパンを買う",
		DueDate:     &dueNextWeek,
		Priority:    models.PriorityLow,
		Status:      models.StatusInProgress,
	})
	_, _ = assignmentService.Create(1, dto.CreateAssignmentInput{
		Title:       "読書",
		Description: "Go公式ドキュメントを読む",
		Priority:    models.PriorityMedium,
		Status:      models.StatusDone,
	})
}
