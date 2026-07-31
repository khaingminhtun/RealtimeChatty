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
	// Drift system (IMPORTANT)
	DriftThresholdDays *int `json:"drift_threshold_days,omitempty"`
}

type RelationshipResponseDTO struct {
	ID            int64     `json:"id"`
	OwnerID       int64     `json:"owner_id"`

	Name          string    `json:"name"`
	Type          string    `json:"type"`
	HowWeMet      string    `json:"how_we_met"`

	Birthday      *string   `json:"birthday"`
	Location      string    `json:"location"`
	Tags          []string  `json:"tags"`

	// Drift system (core)
	DriftThresholdDays  int        `json:"drift_threshold_days"`
	DriftStatus         string     `json:"drift_status"`
	WarmthScore         int        `json:"warmth_score"`

	LastContactAt       *time.Time `json:"last_contact_at"`
	NextContactAt       *time.Time `json:"next_contact_at"`
	LastReminderSentAt  *time.Time `json:"last_reminder_sent_at"`

	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UpdateRelationshipDTO struct {
	ID      int64
	OwnerID int64

	Name     *string
	Type     *string
	HowWeMet *string
	Birthday *time.Time
	Location *string
	AvatarURL *string
	Tags     []string

	// Drift system updates (optional overrides)
	DriftThresholdDays *int
}

type TagsUpdateResponseDTO struct {
	ID   int64    `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
