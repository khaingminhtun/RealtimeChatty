package relationships

import (
	"context"
	"time"

	"github.com/khaingminhtun/realtimechatty/internal/db"
	"github.com/khaingminhtun/realtimechatty/internal/pkg/dbutils"
)

type RelationshipService interface {
	AddRelationship(ctx context.Context, dto CreateRelationshipDTO) (RelationshipResponseDTO, error)
	GetRelationship(ctx context.Context, id int64, ownerID int64) (RelationshipResponseDTO, error)
	ListRelationships(ctx context.Context, ownerID int64, relType string) ([]RelationshipResponseDTO, error)
	UpdateRelationship(ctx context.Context, dto UpdateRelationshipDTO) (RelationshipResponseDTO, error)
	DeleteRelationship(ctx context.Context, id int64, ownerID int64) error

	//tags
	ReplaceTags(ctx context.Context, id int64, ownerID int64, tags []string) (TagsUpdateResponseDTO, error)
	AppendTags(ctx context.Context, id int64, ownerID int64, tags []string) (TagsUpdateResponseDTO, error)
	RemoveTag(ctx context.Context, id int64, ownerID int64, tag string) (TagsUpdateResponseDTO, error)
}

type relationshipService struct {
	repo RelationshipRepository
}

func NewRelationshipService(repo RelationshipRepository) RelationshipService {
	return &relationshipService{repo: repo}
}

func (s *relationshipService) AddRelationship(
	ctx context.Context,
	dto CreateRelationshipDTO,
) (RelationshipResponseDTO, error) {

	// =====================================================
	// 1. Resolve drift threshold
	// =====================================================

	driftDays := dto.DriftThresholdDays
	if driftDays == nil {
		def := defaultDriftByType(dto.Type)
		driftDays = &def
	}

	// =====================================================
	// 2. Compute next_contact_at
	// =====================================================

	now := time.Now()

	nextContactAt := now.AddDate(0, 0, *driftDays)

	// =====================================================
	// 3. Build DB insert params
	// =====================================================

	arg := db.CreateRelationshipParams{
		OwnerID: dto.OwnerID,
		Name:    dto.Name,
		Type:    dbutils.NewText(dto.Type),

		HowWeMet: dbutils.NewText(dto.HowWeMet),
		Birthday: dbutils.ParseDateString(dto.Birthday),
		Location: dbutils.NewText(dto.Location),
		Tags:     dto.Tags,

		DriftThresholdDays: int32(*driftDays),
		NextContactAt:      dbutils.NewTimestamp(nextContactAt),
	}

	// =====================================================
	// 4. Insert into DB
	// =====================================================

	dbResult, err := s.repo.CreateRelationship(ctx, arg)
	if err != nil {
		return RelationshipResponseDTO{}, err
	}

	// =====================================================
	// 5. Map response
	// =====================================================

	return mapToResponseDTO(dbResult), nil
}

func (s *relationshipService) GetRelationship(ctx context.Context, id int64, ownerID int64) (RelationshipResponseDTO, error) {
	dbrs, err := s.repo.GetByID(ctx, id, ownerID)
	if err != nil {
		return RelationshipResponseDTO{}, err
	}
	return mapToResponseDTO(dbrs), nil
}

func (s *relationshipService) ListRelationships(ctx context.Context, ownerID int64, relType string) ([]RelationshipResponseDTO, error) {
	dbList, err := s.repo.List(ctx, ownerID, relType)
	if err != nil {
		return nil, err
	}

	// Instantiate the slice to prevent returning an accidental 'null' in JSON response
	dtos := make([]RelationshipResponseDTO, 0, len(dbList))
	for _, row := range dbList {
		dtos = append(dtos, mapListRowToResponseDTO(row))
	}

	return dtos, nil
}

func (s *relationshipService) UpdateRelationship(ctx context.Context, dto UpdateRelationshipDTO) (RelationshipResponseDTO, error) {
	params := db.UpdateRelationshipParams{
		ID:      dto.ID,
		OwnerID: dto.OwnerID,
	}

	// Safely map pointer values into pgtype wrappers for partial updates
	if dto.Name != nil {
		params.Name = dbutils.NewText(*dto.Name)
	}
	if dto.Type != nil {
		params.Type = dbutils.NewText(*dto.Type)
	}
	if dto.HowWeMet != nil {
		params.HowWeMet = dbutils.NewText(*dto.HowWeMet)
	}
	if dto.Location != nil {
		params.Location = dbutils.NewText(*dto.Location)
	}
	if dto.AvatarURL != nil {
		params.AvatarUrl = dbutils.NewText(*dto.AvatarURL)
	}
	if dto.Birthday != nil {
		params.Birthday = dbutils.NewDate(*dto.Birthday)
	}

	updatedRel, err := s.repo.Update(ctx, params)
	if err != nil {
		return RelationshipResponseDTO{}, err
	}

	return mapToResponseDTO(updatedRel), nil
}

func (s *relationshipService) DeleteRelationship(ctx context.Context, id int64, ownerID int64) error {
	return s.repo.Delete(ctx, id, ownerID)
}

func (s *relationshipService) ReplaceTags(ctx context.Context, id int64, ownerID int64, tags []string) (TagsUpdateResponseDTO, error) {
	if tags == nil {
		tags = []string{}
	}
	res, err := s.repo.ReplaceTags(ctx, tags, id, ownerID)
	if err != nil {
		return TagsUpdateResponseDTO{}, err
	}
	return TagsUpdateResponseDTO{ID: res.ID, Name: res.Name, Tags: res.Tags}, nil
}

func (s *relationshipService) AppendTags(ctx context.Context, id int64, ownerID int64, tags []string) (TagsUpdateResponseDTO, error) {
	res, err := s.repo.AppendTags(ctx, tags, id, ownerID)
	if err != nil {
		return TagsUpdateResponseDTO{}, err
	}
	return TagsUpdateResponseDTO{ID: res.ID, Name: res.Name, Tags: res.Tags}, nil
}

func (s *relationshipService) RemoveTag(ctx context.Context, id int64, ownerID int64, tag string) (TagsUpdateResponseDTO, error) {
	res, err := s.repo.RemoveTag(ctx, tag, id, ownerID)
	if err != nil {
		return TagsUpdateResponseDTO{}, err
	}
	return TagsUpdateResponseDTO{ID: res.ID, Name: res.Name, Tags: res.Tags}, nil
}

func mapToResponseDTO(r db.Relationship) RelationshipResponseDTO {

	var birthdayStr *string
	if r.Birthday.Valid {
		formatted := r.Birthday.Time.Format("2006-01-02")
		birthdayStr = &formatted
	}

	var lastContact *time.Time
	if r.LastContactAt.Valid {
		lastContact = &r.LastContactAt.Time
	}

	var nextContact *time.Time
	if r.NextContactAt.Valid {
		nextContact = &r.NextContactAt.Time
	}

	var lastReminder *time.Time
	if r.LastReminderSentAt.Valid {
		lastReminder = &r.LastReminderSentAt.Time
	}

	return RelationshipResponseDTO{
		ID:       r.ID,
		OwnerID:  r.OwnerID,
		Name:     r.Name,
		Type:     r.Type.String,
		HowWeMet: r.HowWeMet.String,
		Birthday: birthdayStr,
		Location: r.Location.String,
		Tags:     r.Tags,

		// 🔥 Drift system fields
		DriftThresholdDays: int(r.DriftThresholdDays),
		DriftStatus:        r.DriftStatus,
		WarmthScore:        int(r.WarmthScore.Int32),

		LastContactAt:      lastContact,
		NextContactAt:      nextContact,
		LastReminderSentAt: lastReminder,

		CreatedAt: r.CreatedAt.Time,
		UpdatedAt: r.UpdatedAt.Time,
	}
}

func mapListRowToResponseDTO(r db.ListRelationshipsRow) RelationshipResponseDTO {

	var birthdayStr *string
	if r.Birthday.Valid {
		b := r.Birthday.Time.Format("2006-01-02")
		birthdayStr = &b
	}

	var lastContact *time.Time
	if r.LastContactAt.Valid {
		lastContact = &r.LastContactAt.Time
	}

	var nextContact *time.Time
	if r.NextContactAt.Valid {
		nextContact = &r.NextContactAt.Time
	}

	return RelationshipResponseDTO{
		ID:       r.ID,
		OwnerID:  r.OwnerID,
		Name:     r.Name,
		Type:     r.Type.String,
		HowWeMet: r.HowWeMet.String,
		Location: r.Location.String,
		Birthday: birthdayStr,
		Tags:     r.Tags,

		DriftThresholdDays: int(r.DriftThresholdDays),
		DriftStatus:        r.DriftStatus,
		WarmthScore:        int(r.WarmthScore.Int32),

		LastContactAt: lastContact,
		NextContactAt: nextContact,

		CreatedAt: r.CreatedAt.Time,
	}
}

func defaultDriftByType(t string) int {
	switch t {
	case "friend":
		return 7
	case "partner":
		return 3
	case "family":
		return 14
	default:
		return 7
	}
}
