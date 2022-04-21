package usecase

import (
	"context"
	"github.com/thinhlu123/shortener/internal/models"
	"github.com/thinhlu123/shortener/internal/user"
	"github.com/thinhlu123/shortener/pkg/auth"
	"github.com/thinhlu123/shortener/pkg/logger"
)

func NewUserUsecase(repo user.UserRepository, log logger.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo: repo,
		log:      log,
	}
}

type UserUsecase struct {
	userRepo user.UserRepository
	log      logger.Logger
}

func (u *UserUsecase) Login(ctx context.Context, usr string, pwd string) (string, error) {
	_, err := u.userRepo.Login(ctx, models.User{
		Usr: usr,
		Pwd: pwd,
	})
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateAuthToken(usr, pwd)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUsecase) Register(ctx context.Context, user models.User) error {
	if err := user.HashPwd(); err != nil {
		return err
	}

	if err := u.userRepo.Register(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, filter models.User, updater models.User) error {
	if err := u.userRepo.UpdateUser(ctx, filter, updater); err != nil {
		return err
	}

	return nil
}

func (u *UserUsecase) Withdraw(ctx context.Context, amount int64) error {
	//TODO implement me
	panic("implement me")
}
