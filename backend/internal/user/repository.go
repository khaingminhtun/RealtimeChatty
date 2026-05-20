package user

import (
	"context"

	"github.com/khaingminhtun/realtimechatty/internal/db"
)

type UserRepository interface {
	CreateUser(
		ctx context.Context,
		arg db.CreateUserParams,
	) (db.User, error)
	GetUserByEmail(
		ctx context.Context,
		email string,
	) (db.User, error)

	GetUserByUsername(
		ctx context.Context,
		username string,
	) (db.User, error)
}

type userRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) UserRepository {
	return &userRepository{q: q}
}

func (r *userRepository) CreateUser(
	ctx context.Context,
	arg db.CreateUserParams,
) (db.User, error) {
	return r.q.CreateUser(ctx, arg)
}

func (r *userRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (db.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *userRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (db.User, error) {
	return r.q.GetUserByUsername(ctx, username)
}
