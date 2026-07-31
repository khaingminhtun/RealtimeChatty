package contacts

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/khaingminhtun/realtimechatty/internal/db"
	"github.com/khaingminhtun/realtimechatty/internal/pkg/dbutils"
)

// ---------------------------------------------------------------------------
// Interface
// ---------------------------------------------------------------------------

type ContactService interface {
	// LogContact inserts a new contact log entry.
	LogContact(ctx context.Context, dto LogContactDTO) (ContactResponseDTO, error)

	// GetContact fetches a single contact log (ownership-checked).
	GetContact(ctx context.Context, id int64, userID int64) (ContactResponseDTO, error)

	// GetContactsByRelationship returns all contact logs for a relationship, newest first.
	GetContactsByRelationship(ctx context.Context, relationshipID int64) ([]ContactResponseDTO, error)

	// UpdateContact performs a partial update on a contact log entry.
	UpdateContact(ctx context.Context, dto UpdateContactDTO) (ContactResponseDTO, error)

	// DeleteContact removes a contact log entry (ownership-checked).
	DeleteContact(ctx context.Context, id int64, userID int64) error

	// ListDriftReminders returns relationships whose next_contact_at has passed.
	ListDriftReminders(ctx context.Context) ([]DriftReminderDTO, error)

	// MarkReminderSent updates last_reminder_sent_at to NOW() for a relationship.
	MarkReminderSent(ctx context.Context, relationshipID int64) error

	// ListRelationshipsForDrift returns all relationships with drift metadata for a user.
	ListRelationshipsForDrift(ctx context.Context, ownerID int64) ([]DriftReminderDTO, error)

	// SearchRelationships performs full-text search over a user's relationships.
	SearchRelationships(ctx context.Context, ownerID int64, query string) ([]SearchResultDTO, error)
}

// ---------------------------------------------------------------------------
// Implementation
// ---------------------------------------------------------------------------

type contactService struct {
	repo ContactRepository
}

func NewContactService(repo ContactRepository) ContactService {
	return &contactService{repo: repo}
}

// ---------------------------------------------------------------------------
// LogContact
// ---------------------------------------------------------------------------

func (s *contactService) LogContact(ctx context.Context, dto LogContactDTO) (ContactResponseDTO, error) {
	contactedAt, err := resolveContactedAt(dto.ContactedAt)
	if err != nil {
		return ContactResponseDTO{}, fmt.Errorf("invalid contacted_at timestamp: %w", err)
	}

	params := db.CreateContactLogParams{
		RelationshipID: dto.RelationshipID,
		UserID:         dto.UserID,
		Channel:        dto.Channel,
		Note:           dbutils.NewText(stringOrEmpty(dto.Note)),
		// Column5 maps to the COALESCE($5, NOW()) param — pass nil to use NOW()
		Column5: pgtype.Timestamptz{Time: contactedAt.UTC(), Valid: true},
	}

	created, err := s.repo.CreateContactLog(ctx, params)
	if err != nil {
		return ContactResponseDTO{}, fmt.Errorf("create contact log: %w", err)
	}

	return mapContactToDTO(created), nil
}

// ---------------------------------------------------------------------------
// GetContact
// ---------------------------------------------------------------------------

func (s *contactService) GetContact(ctx context.Context, id int64, userID int64) (ContactResponseDTO, error) {
	row, err := s.repo.GetContactByID(ctx, db.GetContactByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return ContactResponseDTO{}, fmt.Errorf("get contact: %w", err)
	}
	return mapContactToDTO(row), nil
}

// ---------------------------------------------------------------------------
// GetContactsByRelationship
// ---------------------------------------------------------------------------

func (s *contactService) GetContactsByRelationship(ctx context.Context, relationshipID int64) ([]ContactResponseDTO, error) {
	rows, err := s.repo.GetContactsByRelationship(ctx, relationshipID)
	if err != nil {
		return nil, err
	}

	dtos := make([]ContactResponseDTO, 0, len(rows))
	for _, row := range rows {
		dtos = append(dtos, mapContactToDTO(row))
	}
	return dtos, nil
}

// ---------------------------------------------------------------------------
// UpdateContact
// ---------------------------------------------------------------------------

func (s *contactService) UpdateContact(ctx context.Context, dto UpdateContactDTO) (ContactResponseDTO, error) {
	// Build the pgtype-aware params; COALESCE in SQL means zero-value = keep current.
	params := db.UpdateContactParams{
		ID:     dto.ID,
		UserID: dto.UserID,
	}

	if dto.Channel != nil {
		params.Channel = *dto.Channel
	}

	if dto.Note != nil {
		params.Note = dbutils.NewText(*dto.Note)
	}

	if dto.ContactedAt != nil && *dto.ContactedAt != "" {
		t, err := time.Parse(time.RFC3339, *dto.ContactedAt)
		if err != nil {
			return ContactResponseDTO{}, fmt.Errorf("invalid contacted_at: %w", err)
		}
		params.ContactedAt = pgtype.Timestamptz{Time: t.UTC(), Valid: true}
	}

	updated, err := s.repo.UpdateContact(ctx, params)
	if err != nil {
		return ContactResponseDTO{}, fmt.Errorf("update contact: %w", err)
	}

	return mapContactToDTO(updated), nil
}

// ---------------------------------------------------------------------------
// DeleteContact
// ---------------------------------------------------------------------------

func (s *contactService) DeleteContact(ctx context.Context, id int64, userID int64) error {
	return s.repo.DeleteContact(ctx, db.DeleteContactParams{
		ID:     id,
		UserID: userID,
	})
}

// ---------------------------------------------------------------------------
// ListDriftReminders
// ---------------------------------------------------------------------------

func (s *contactService) ListDriftReminders(ctx context.Context) ([]DriftReminderDTO, error) {
	rows, err := s.repo.GetPendingDriftReminders(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]DriftReminderDTO, 0, len(rows))
	for _, row := range rows {
		dtos = append(dtos, mapPendingReminderToDTO(row))
	}
	return dtos, nil
}

// ---------------------------------------------------------------------------
// MarkReminderSent
// ---------------------------------------------------------------------------

func (s *contactService) MarkReminderSent(ctx context.Context, relationshipID int64) error {
	return s.repo.MarkReminderAsSent(ctx, relationshipID)
}

// ---------------------------------------------------------------------------
// ListRelationshipsForDrift
// ---------------------------------------------------------------------------

func (s *contactService) ListRelationshipsForDrift(ctx context.Context, ownerID int64) ([]DriftReminderDTO, error) {
	rows, err := s.repo.ListRelationshipsForDrift(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	dtos := make([]DriftReminderDTO, 0, len(rows))
	for _, row := range rows {
		dtos = append(dtos, mapDriftRowToDTO(row))
	}
	return dtos, nil
}

// ---------------------------------------------------------------------------
// SearchRelationships
// ---------------------------------------------------------------------------

func (s *contactService) SearchRelationships(ctx context.Context, ownerID int64, query string) ([]SearchResultDTO, error) {
	rows, err := s.repo.SearchRelationships(ctx, db.SearchRelationshipsParams{
		OwnerID:        ownerID,
		PlaintoTsquery: query,
	})
	if err != nil {
		return nil, err
	}

	dtos := make([]SearchResultDTO, 0, len(rows))
	for _, row := range rows {
		dtos = append(dtos, mapRelationshipToSearchDTO(row))
	}
	return dtos, nil
}

// ---------------------------------------------------------------------------
// Mapping helpers
// ---------------------------------------------------------------------------

func mapContactToDTO(c db.Contact) ContactResponseDTO {
	var note *string
	if c.Note.Valid {
		n := c.Note.String
		note = &n
	}

	contactedAt := time.Now()
	if c.ContactedAt.Valid {
		contactedAt = c.ContactedAt.Time
	}

	createdAt := time.Now()
	if c.CreatedAt.Valid {
		createdAt = c.CreatedAt.Time
	}

	return ContactResponseDTO{
		ID:             c.ID,
		RelationshipID: c.RelationshipID,
		UserID:         c.UserID,
		Channel:        c.Channel,
		Note:           note,
		ContactedAt:    contactedAt,
		CreatedAt:      createdAt,
	}
}

func mapPendingReminderToDTO(r db.GetPendingDriftRemindersRow) DriftReminderDTO {
	var lastContact *time.Time
	if r.LastContactAt.Valid {
		t := r.LastContactAt.Time
		lastContact = &t
	}

	var nextContact *time.Time
	if r.NextContactAt.Valid {
		t := r.NextContactAt.Time
		nextContact = &t
	}

	return DriftReminderDTO{
		RelationshipID:     r.ID,
		OwnerID:            r.OwnerID,
		Name:               r.Name,
		Type:               r.Type.String,
		DriftThresholdDays: int(r.DriftThresholdDays),
		LastContactAt:      lastContact,
		NextContactAt:      nextContact,
	}
}

func mapDriftRowToDTO(r db.ListRelationshipsForDriftRow) DriftReminderDTO {
	var lastContact *time.Time
	if r.LastContactAt.Valid {
		t := r.LastContactAt.Time
		lastContact = &t
	}

	var nextContact *time.Time
	if r.NextContactAt.Valid {
		t := r.NextContactAt.Time
		nextContact = &t
	}

	return DriftReminderDTO{
		RelationshipID:     r.ID,
		OwnerID:            r.OwnerID,
		Name:               r.Name,
		Type:               r.Type.String,
		DriftThresholdDays: int(r.DriftThresholdDays),
		LastContactAt:      lastContact,
		NextContactAt:      nextContact,
	}
}

func mapRelationshipToSearchDTO(r db.Relationship) SearchResultDTO {
	var lastContact *time.Time
	if r.LastContactAt.Valid {
		t := r.LastContactAt.Time
		lastContact = &t
	}

	var nextContact *time.Time
	if r.NextContactAt.Valid {
		t := r.NextContactAt.Time
		nextContact = &t
	}

	return SearchResultDTO{
		ID:            r.ID,
		OwnerID:       r.OwnerID,
		Name:          r.Name,
		Type:          r.Type.String,
		DriftStatus:   r.DriftStatus,
		LastContactAt: lastContact,
		NextContactAt: nextContact,
	}
}

// ---------------------------------------------------------------------------
// Utilities
// ---------------------------------------------------------------------------

// resolveContactedAt parses an optional RFC3339 string; returns time.Now() if nil/empty.
func resolveContactedAt(s *string) (time.Time, error) {
	if s == nil || *s == "" {
		return time.Now(), nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
