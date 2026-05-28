package relationships

import (
	"context"

	"github.com/khaingminhtun/realtimechatty/internal/db"
)

type RelationshipRepository interface {
	CreateRelationship(ctx context.Context, arg db.CreateRelationshipParams) (db.Relationship, error)
	GetByID(ctx context.Context, id int64, ownerID int64) (db.Relationship, error)
	List(ctx context.Context, ownerID int64, relType string) ([]db.ListRelationshipsRow, error)
	Update(ctx context.Context, arg db.UpdateRelationshipParams) (db.Relationship, error)
	Delete(ctx context.Context, id int64, ownerID int64) error
}

type relationshipRepository struct {
	q *db.Queries
}

func NewRelationShipRepository(q *db.Queries) RelationshipRepository {
	return &relationshipRepository{
		q: q,
	}
}

func (r *relationshipRepository) CreateRelationship(ctx context.Context, arg db.CreateRelationshipParams) (db.Relationship, error) {
	return r.q.CreateRelationship(ctx, arg)
}

func (r *relationshipRepository) GetByID(ctx context.Context, id int64, ownerID int64) (db.Relationship, error) {
	return r.q.GetRelationshipByID(ctx, db.GetRelationshipByIDParams{ID: id, OwnerID: ownerID})
}

func (r *relationshipRepository) List(ctx context.Context, ownerID int64, relType string) ([]db.ListRelationshipsRow, error) {
	return r.q.ListRelationships(ctx, db.ListRelationshipsParams{OwnerID: ownerID, Column2: relType})
}

func (r *relationshipRepository) Update(ctx context.Context, arg db.UpdateRelationshipParams) (db.Relationship, error) {
	return r.q.UpdateRelationship(ctx, arg)
}

func (r *relationshipRepository) Delete(ctx context.Context, id int64, ownerID int64) error {
	return r.q.DeleteRelationship(ctx, db.DeleteRelationshipParams{ID: id, OwnerID: ownerID})
}
