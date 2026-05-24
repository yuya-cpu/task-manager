package models

import "time"

// 優先度
const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)

// ステータス
const (
	StatusTodo       = "todo"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
)

// Assignment はタスク（課題）を表すモデルです。
type Assignment struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Priority    string     `json:"priority" gorm:"not null;default:medium"`
	Status      string     `json:"status" gorm:"not null;default:todo"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ValidPriority(p string) bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}

func ValidStatus(s string) bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
