package main

import (
	"log"
	"time"

	"task-manager/data"
	"task-manager/dto"
	"task-manager/handlers"
	"task-manager/models"
	"task-manager/repositories"
	"task-manager/services"

	"github.com/gin-gonic/gin"
)

func setupRouter(assignmentHandler handlers.AssignmentHandler) *gin.Engine {
	router := gin.Default()

	group := router.Group("/assignments")
	group.GET("", assignmentHandler.FindAll)
	group.GET("/:id", assignmentHandler.FindByID)
	group.POST("", assignmentHandler.Create)
	group.PUT("/:id", assignmentHandler.Update)
	group.DELETE("/:id", assignmentHandler.Delete)

	return router
}

func seedIfEmpty(service services.AssignmentService) {
	assignments, err := service.FindAll()
	if err != nil || len(assignments) > 0 {
		return
	}

	dueTomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	dueNextWeek := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	_, _ = service.Create(dto.CreateAssignmentInput{
		Title:       "Goの課題を提出する",
		Description: "task-manager APIの実装を完了する",
		DueDate:     &dueTomorrow,
		Priority:    models.PriorityHigh,
		Status:      models.StatusTodo,
	})
	_, _ = service.Create(dto.CreateAssignmentInput{
		Title:       "買い物リスト",
		Description: "牛乳とパンを買う",
		DueDate:     &dueNextWeek,
		Priority:    models.PriorityLow,
		Status:      models.StatusInProgress,
	})
	_, _ = service.Create(dto.CreateAssignmentInput{
		Title:       "読書",
		Description: "Go公式ドキュメントを読む",
		Priority:    models.PriorityMedium,
		Status:      models.StatusDone,
	})
}

func main() {
	db := data.SetupDB()

	repository := repositories.NewAssignmentRepository(db)
	service := services.NewAssignmentService(repository)
	assignmentHandler := handlers.NewAssignmentHandler(service)

	seedIfEmpty(service)

	router := setupRouter(assignmentHandler)
	log.Println("server started at http://127.0.0.1:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
