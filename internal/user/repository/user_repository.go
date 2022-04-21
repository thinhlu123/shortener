package repository

import (
	"context"
	"errors"
	"github.com/thinhlu123/shortener/internal/models"
)

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

type UserRepository struct{}

func (u *UserRepository) Login(ctx context.Context, user models.User) (string, error) {
	var us models.User
	err := models.UserDB.GetCollection().FindOne(ctx, models.User{
		Usr: user.Usr,
	}).Decode(&u)
	if err != nil {
		return "", err
	}

	if !us.ComparePwd(user.Pwd) {
		return "", errors.New("password not match")
	}

	return "", nil
}

func (u *UserRepository) Register(ctx context.Context, user models.User) error {
	_, err := models.UserDB.GetCollection().InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, filter models.User, updater models.User) error {
	_, err := models.UserDB.GetCollection().UpdateOne(ctx, filter, updater)
	if err != nil {
		return err
	}

	return nil
}
