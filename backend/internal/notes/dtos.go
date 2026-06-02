package notes

import "time"

type NoteResponseDTO struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	RelationshipID int64     `json:"relationship_id"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
