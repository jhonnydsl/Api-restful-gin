package dtos

type Task struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}