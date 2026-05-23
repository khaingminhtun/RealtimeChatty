package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khaingminhtun/realtimechatty/internal/db"
)

type AuthRepository interface {
	// Transaction Core
	WithTransaction(ctx context.Context, fn func(txRepo AuthRepository) error) error

	// Auth Queries
	CreateUserAuth(ctx context.Context, arg db.CreateUserAuthParams) (db.UserAuth, error)
	GetUserAuthByUserID(ctx context.Context, id int64) (db.UserAuth, error)

	MarkUserAsVerified(ctx context.Context, email string) error

	CreateRefreshToken(ctx context.Context, arg db.CreateRefreshTokenParams) error

	// Session Queries
	CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error)
}

type authRepository struct {
	dbPool *pgxpool.Pool
	q      *db.Queries
}

func NewAuthRepository(dbPool *pgxpool.Pool) AuthRepository {
	return &authRepository{
		dbPool: dbPool,
		q:      db.New(dbPool),
	}
}

func (r *authRepository) WithTransaction(ctx context.Context, fn func(txRepo AuthRepository) error) error {
	tx, err := r.dbPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	txRepo := &authRepository{
		dbPool: r.dbPool,
		q:      db.New(tx), // Bind all queries to the isolated transaction instance
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

// --- Auth Implementations ---

func (r *authRepository) CreateUserAuth(ctx context.Context, arg db.CreateUserAuthParams) (db.UserAuth, error) {
	return r.q.CreateUserAuth(ctx, arg)
}

func (r *authRepository) GetUserAuthByUserID(ctx context.Context, id int64) (db.UserAuth, error) {
	return r.q.GetUserAuthByUserID(ctx, id)
}

func (r *authRepository) MarkUserAsVerified(ctx context.Context, email string) error {
	return r.q.MarkUserAsVerified(ctx, email)
}

func (r *authRepository) CreateRefreshToken(ctx context.Context, arg db.CreateRefreshTokenParams) error {
	_, err := r.q.CreateRefreshToken(ctx, arg)
	return err
}

// --- Session Implementations ---

func (r *authRepository) CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error) {
	return r.q.CreateSession(ctx, arg)
}
