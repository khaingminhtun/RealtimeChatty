package auth

import (
	"context"

	"github.com/khaingminhtun/realtimechatty/internal/db"
)

type AuthRepository interface {
	CreateUserAuth(
		ctx context.Context,
		arg db.CreateUserAuthParams,
	) (db.UserAuth, error)

	GetUserAuthByUserID(
		ctx context.Context,
		id int64,
	) (db.UserAuth, error)
}

type authRepository struct {
	q *db.Queries
}

func NewAuthRepository(q *db.Queries) AuthRepository {
	return &authRepository{q: q}
}

func (r *authRepository) CreateUserAuth(
	ctx context.Context,
	arg db.CreateUserAuthParams,
) (db.UserAuth, error) {
	return r.q.CreateUserAuth(ctx, arg)
}

func (r *authRepository) GetUserAuthByUserID(
	ctx context.Context,
	id int64,
) (db.UserAuth, error) {
	return r.q.GetUserAuthByUserID(ctx, id)
}
