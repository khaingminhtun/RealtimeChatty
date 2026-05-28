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
}

type relationshipService struct {
	repo RelationshipRepository
}

func NewRelationshipService(repo RelationshipRepository) RelationshipService {
	return &relationshipService{repo: repo}
}

func (s *relationshipService) AddRelationship(ctx context.Context, dto CreateRelationshipDTO) (RelationshipResponseDTO, error) {
	arg := db.CreateRelationshipParams{
		OwnerID:  dto.OwnerID,
		Name:     dto.Name,
		Type:     dbutils.NewText(dto.Type),
		HowWeMet: dbutils.NewText(dto.HowWeMet),
		Birthday: dbutils.ParseDateString(dto.Birthday),
		Location: dbutils.NewText(dto.Location),
		Tags:     dto.Tags,
	}

	// 2. Call repository layer
	dbResult, err := s.repo.CreateRelationship(ctx, arg)
	if err != nil {
		return RelationshipResponseDTO{}, err
	}

	// 3. Map db.Relationship to clean RelationshipResponseDTO
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

// Internal mapping logic to isolate pgx/pgtype specifics from the controller
func mapToResponseDTO(r db.Relationship) RelationshipResponseDTO {
	var birthdayStr *string
	if r.Birthday.Valid {
		// Formats the pgtype.Date time back to YYYY-MM-DD
		formatted := r.Birthday.Time.Format("2006-01-02")
		birthdayStr = &formatted
	}

	var lastContact *time.Time
	if r.LastContactAt.Valid {
		lastContact = &r.LastContactAt.Time
	}

	return RelationshipResponseDTO{
		ID:            r.ID,
		OwnerID:       r.OwnerID,
		Name:          r.Name,
		Type:          r.Type.String,     // Safe access: extracts empty string if SQL NULL
		HowWeMet:      r.HowWeMet.String, // Safe access
		Birthday:      birthdayStr,       // Handled cleanly as a pointer string
		Location:      r.Location.String, // Safe access
		Tags:          r.Tags,
		LastContactAt: lastContact,
		CreatedAt:     r.CreatedAt.Time,
		UpdatedAt:     r.UpdatedAt.Time,
	}
}

func mapListRowToResponseDTO(r db.ListRelationshipsRow) RelationshipResponseDTO {
	var birthdayStr *string
	if r.Birthday.Valid {
		birthdayStr = new(string)
		*birthdayStr = r.Birthday.Time.Format("2006-01-02")
	}

	var lastContact *time.Time
	if r.LastContactAt.Valid {
		lastContact = &r.LastContactAt.Time
	}

	return RelationshipResponseDTO{
		ID:            r.ID,
		OwnerID:       r.OwnerID,
		Name:          r.Name,
		Type:          r.Type.String,
		HowWeMet:      r.HowWeMet.String, // Now safely populated!
		Location:      r.Location.String, // Now safely populated!
		Birthday:      birthdayStr,
		Tags:          r.Tags,
		LastContactAt: lastContact,
		CreatedAt:     r.CreatedAt.Time,
		UpdatedAt:     time.Time{}, // Still zeroed out unless you add updated_at to the SELECT statement too
	}
}
