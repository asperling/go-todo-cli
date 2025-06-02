package todos

import "github.com/google/uuid"

type Todo struct {
	ID        uuid.UUID `json:"id"`
	Task      string    `json:"task"`
	Completed bool      `json:"completed"`
}
