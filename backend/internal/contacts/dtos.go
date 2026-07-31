package contacts

import "time"

// ---------------------------------------------------------------------------
// Request DTOs
// ---------------------------------------------------------------------------

// LogContactDTO is the request body for POST /contacts/relationships/:id/contacts
type LogContactDTO struct {
	RelationshipID int64   `json:"-"` // injected from URL path param
	UserID         int64   `json:"-"` // injected from auth middleware

	Channel     string  `json:"channel"`      // required: "call", "message", "email", "in_person"
	Note        *string `json:"note"`         // optional free-text note
	ContactedAt *string `json:"contacted_at"` // optional RFC3339; defaults to NOW()
}

// UpdateContactDTO is the request body for PATCH /contacts/:id
type UpdateContactDTO struct {
	ID     int64 `json:"-"` // injected from URL path param
	UserID int64 `json:"-"` // injected from auth middleware

	Channel     *string `json:"channel"`
	Note        *string `json:"note"`
	ContactedAt *string `json:"contacted_at"` // optional RFC3339 override
}

// ---------------------------------------------------------------------------
// Response DTOs
// ---------------------------------------------------------------------------

// ContactResponseDTO is the canonical API shape for a contact log entry.
type ContactResponseDTO struct {
	ID             int64      `json:"id"`
	RelationshipID int64      `json:"relationship_id"`
	UserID         int64      `json:"user_id"`
	Channel        string     `json:"channel"`
	Note           *string    `json:"note"`
	ContactedAt    time.Time  `json:"contacted_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

// DriftReminderDTO represents a relationship that is overdue for contact
// or is listed for drift scheduling purposes.
type DriftReminderDTO struct {
	RelationshipID     int64      `json:"relationship_id"`
	OwnerID            int64      `json:"owner_id"`
	Name               string     `json:"name"`
	Type               string     `json:"type"`
	DriftThresholdDays int        `json:"drift_threshold_days"`
	LastContactAt      *time.Time `json:"last_contact_at"`
	NextContactAt      *time.Time `json:"next_contact_at"`
}

// SearchResultDTO is a slim view of a relationship returned from full-text search.
type SearchResultDTO struct {
	ID            int64      `json:"id"`
	OwnerID       int64      `json:"owner_id"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	DriftStatus   string     `json:"drift_status"`
	LastContactAt *time.Time `json:"last_contact_at"`
	NextContactAt *time.Time `json:"next_contact_at"`
}
