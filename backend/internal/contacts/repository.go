package contacts

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khaingminhtun/realtimechatty/internal/db"
)

// ---------------------------------------------------------------------------
// Interface
// ---------------------------------------------------------------------------

type ContactRepository interface {
	// Transaction support
	WithTransaction(ctx context.Context, fn func(txRepo ContactRepository) error) error

	// Contact log CRUD
	CreateContactLog(ctx context.Context, arg db.CreateContactLogParams) (db.Contact, error)
	GetContactByID(ctx context.Context, arg db.GetContactByIDParams) (db.Contact, error)
	GetContactsByRelationship(ctx context.Context, relationshipID int64) ([]db.Contact, error)
	UpdateContact(ctx context.Context, arg db.UpdateContactParams) (db.Contact, error)
	DeleteContact(ctx context.Context, arg db.DeleteContactParams) error

	// Drift / reminder queries (operate on the relationships table)
	GetPendingDriftReminders(ctx context.Context) ([]db.GetPendingDriftRemindersRow, error)
	MarkReminderAsSent(ctx context.Context, id int64) error
	ListRelationshipsForDrift(ctx context.Context, ownerID int64) ([]db.ListRelationshipsForDriftRow, error)

	// Full-text search (operates on the relationships table)
	SearchRelationships(ctx context.Context, arg db.SearchRelationshipsParams) ([]db.Relationship, error)
}

// ---------------------------------------------------------------------------
// Implementation
// ---------------------------------------------------------------------------

type contactRepository struct {
	dbPool *pgxpool.Pool
	q      *db.Queries
}

func NewContactRepository(dbPool *pgxpool.Pool) ContactRepository {
	return &contactRepository{
		dbPool: dbPool,
		q:      db.New(dbPool),
	}
}

// ---------------------------------------------------------------------------
// Transaction
// ---------------------------------------------------------------------------

func (r *contactRepository) WithTransaction(ctx context.Context, fn func(txRepo ContactRepository) error) error {
	tx, err := r.dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	txRepo := &contactRepository{
		dbPool: r.dbPool,
		q:      db.New(tx),
	}

	err = fn(txRepo)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Contact log CRUD
// ---------------------------------------------------------------------------

func (r *contactRepository) CreateContactLog(ctx context.Context, arg db.CreateContactLogParams) (db.Contact, error) {
	return r.q.CreateContactLog(ctx, arg)
}

func (r *contactRepository) GetContactByID(ctx context.Context, arg db.GetContactByIDParams) (db.Contact, error) {
	return r.q.GetContactByID(ctx, arg)
}

func (r *contactRepository) GetContactsByRelationship(ctx context.Context, relationshipID int64) ([]db.Contact, error) {
	return r.q.GetContactsByRelationship(ctx, relationshipID)
}

func (r *contactRepository) UpdateContact(ctx context.Context, arg db.UpdateContactParams) (db.Contact, error) {
	return r.q.UpdateContact(ctx, arg)
}

func (r *contactRepository) DeleteContact(ctx context.Context, arg db.DeleteContactParams) error {
	return r.q.DeleteContact(ctx, arg)
}

// ---------------------------------------------------------------------------
// Drift / reminder
// ---------------------------------------------------------------------------

func (r *contactRepository) GetPendingDriftReminders(ctx context.Context) ([]db.GetPendingDriftRemindersRow, error) {
	return r.q.GetPendingDriftReminders(ctx)
}

func (r *contactRepository) MarkReminderAsSent(ctx context.Context, id int64) error {
	return r.q.MarkReminderAsSent(ctx, id)
}

func (r *contactRepository) ListRelationshipsForDrift(ctx context.Context, ownerID int64) ([]db.ListRelationshipsForDriftRow, error) {
	return r.q.ListRelationshipsForDrift(ctx, ownerID)
}

// ---------------------------------------------------------------------------
// Search
// ---------------------------------------------------------------------------

func (r *contactRepository) SearchRelationships(ctx context.Context, arg db.SearchRelationshipsParams) ([]db.Relationship, error) {
	return r.q.SearchRelationships(ctx, arg)
}
