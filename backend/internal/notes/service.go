package notes

import (
	"context"

	"github.com/khaingminhtun/realtimechatty/internal/db"
)

type NoteService interface {
	GetNote(ctx context.Context, relationshipID int64, userID int64) (NoteResponseDTO, error)
	UpsertNote(ctx context.Context, relationshipID int64, userID int64, content string) (NoteResponseDTO, error)
	DeleteNote(ctx context.Context, relationshipID int64, userID int64) error
}

type noteService struct {
	repo NoteRepository
}

func NewNoteService(repo NoteRepository) NoteService {
	return &noteService{repo: repo}
}

func (s *noteService) GetNote(ctx context.Context, relationshipID int64, userID int64) (NoteResponseDTO, error) {
	note, err := s.repo.GetByRelationship(ctx, relationshipID, userID)
	if err != nil {
		return NoteResponseDTO{}, err
	}
	return mapNoteToDTO(note), nil
}

func (s *noteService) UpsertNote(ctx context.Context, relationshipID int64, userID int64, content string) (NoteResponseDTO, error) {
	note, err := s.repo.Upsert(ctx, db.UpsertPrivateNoteParams{
		UserID:         userID,
		RelationshipID: relationshipID,
		Content:        content,
	})
	if err != nil {
		return NoteResponseDTO{}, err
	}
	return mapNoteToDTO(note), nil
}

func (s *noteService) DeleteNote(ctx context.Context, relationshipID int64, userID int64) error {
	return s.repo.Delete(ctx,relationshipID, userID)
}

func mapNoteToDTO(n db.PrivateNote) NoteResponseDTO {
	return NoteResponseDTO{
		ID:             n.ID,
		UserID:         n.UserID,
		RelationshipID: n.RelationshipID,
		Content:        n.Content,
		CreatedAt:      n.CreatedAt.Time,
		UpdatedAt:      n.UpdatedAt.Time,
	}
}
