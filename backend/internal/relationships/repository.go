package relationships

import (
	"context"

	"github.com/khaingminhtun/realtimechatty/internal/db"
)

type RelationshipRepository interface {
	CreateRelationship(ctx context.Context, arg db.CreateRelationshipParams) (db.Relationship, error)
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
