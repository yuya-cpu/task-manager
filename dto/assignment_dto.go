package dto

type CreateAssignmentInput struct {
	Title       string  `json:"title" binding:"required,min=1,max=200"`
	Description string  `json:"description"`
	DueDate     *string `json:"due_date" binding:"omitempty,datetime=2006-01-02"`
	Priority    string  `json:"priority" binding:"required,oneof=low medium high"`
	Status      string  `json:"status" binding:"omitempty,oneof=todo in_progress done"`
}

type UpdateAssignmentInput struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=200"`
	Description *string `json:"description"`
	DueDate     *string `json:"due_date" binding:"omitempty,datetime=2006-01-02"`
	Priority    *string `json:"priority" binding:"omitempty,oneof=low medium high"`
	Status      *string `json:"status" binding:"omitempty,oneof=todo in_progress done"`
}
