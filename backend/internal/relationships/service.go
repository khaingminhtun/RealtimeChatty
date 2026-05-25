package relationships

import (
	"context"
	"time"

	"github.com/khaingminhtun/realtimechatty/internal/db"
	"github.com/khaingminhtun/realtimechatty/internal/pkg/dbutils"
)

type RelationshipService interface {
	AddRelationship(ctx context.Context, dto CreateRelationshipDTO) (RelationshipResponseDTO, error)
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
