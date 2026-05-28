package relationships

import "time"

type CreateRelationshipDTO struct {
	OwnerID  int64    `json:"-"` // Handled by auth middleware
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	HowWeMet string   `json:"how_we_met"`
	Birthday string   `json:"birthday"` // YYYY-MM-DD format
	Location string   `json:"location"`
	Tags     []string `json:"tags"`
}

type RelationshipResponseDTO struct {
	ID            int64      `json:"id"`
	OwnerID       int64      `json:"owner_id"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	HowWeMet      string     `json:"how_we_met"`
	Birthday      *string    `json:"birthday"` // Uses a pointer so it can serialize to JSON null if empty
	Location      string     `json:"location"`
	Tags          []string   `json:"tags"`
	LastContactAt *time.Time `json:"last_contact_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type UpdateRelationshipDTO struct {
	ID        int64
	OwnerID   int64
	Name      *string
	Type      *string
	HowWeMet  *string
	Birthday  *time.Time
	Location  *string
	AvatarURL *string
}
