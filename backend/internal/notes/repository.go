package notes

import (
	"context"

	"github.com/khaingminhtun/realtimechatty/internal/db"
)

type NoteRepository interface {
	GetByRelationship(ctx context.Context, relationshipID int64, userID int64) (db.PrivateNote, error)
	Upsert(ctx context.Context, arg db.UpsertPrivateNoteParams) (db.PrivateNote, error)
	Delete(ctx context.Context, relationshipID int64, userID int64) error
}

type noteRepo struct {
	q *db.Queries
}

func NewNoteRepository(queries *db.Queries) NoteRepository {
	return &noteRepo{q: queries}
}

func (r *noteRepo) GetByRelationship(ctx context.Context, relationshipID int64, userID int64) (db.PrivateNote, error) {
	return r.q.GetPrivateNote(ctx, db.GetPrivateNoteParams{RelationshipID: relationshipID, UserID: userID})
}

func (r *noteRepo) Upsert(ctx context.Context, arg db.UpsertPrivateNoteParams) (db.PrivateNote, error) {
	return r.q.UpsertPrivateNote(ctx, arg)
}

func (r *noteRepo) Delete(ctx context.Context, relationshipID int64, userID int64) error {
	return r.q.DeletePrivateNote(ctx, db.DeletePrivateNoteParams{RelationshipID: relationshipID, UserID: userID})
}
