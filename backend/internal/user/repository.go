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
	GetProfile(ctx context.Context, userID int64) (db.User, error)
	UpdateProfile(ctx context.Context, arg db.UpdateUserProfileParams) (db.User, error)
	Delete(ctx context.Context, id int64) error
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

func (r *userRepository) GetProfile(ctx context.Context, userID int64) (db.User, error) {
	return r.q.GetUserByID(ctx, userID)
}

func (r *userRepository) UpdateProfile(ctx context.Context, arg db.UpdateUserProfileParams) (db.User, error) {
	return r.q.UpdateUserProfile(ctx, arg)
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	return r.q.SoftDeleteUser(ctx, id)
}
