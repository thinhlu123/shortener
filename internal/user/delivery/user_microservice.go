package delivery

import (
	"context"
	"github.com/thinhlu123/shortener/internal/models"
	"github.com/thinhlu123/shortener/internal/user"
	pd "github.com/thinhlu123/shortener/internal/user/pb"
	"github.com/thinhlu123/shortener/pkg/auth"
	"github.com/thinhlu123/shortener/pkg/logger"
	"github.com/thinhlu123/shortener/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"regexp"
)

func NewUserMicroservice(userUsecase user.UserUsecase, log logger.Logger) *UserMicroservice {
	return &UserMicroservice{
		userUsecase: userUsecase,
		log:         log,
	}
}

type UserMicroservice struct {
	userUsecase user.UserUsecase
	log         logger.Logger
}

// check pwd has whole lowercase, uppercase, special characters
func isPasswordValid(pwd string) bool {
	if len(pwd) < 8 {
		return false
	}

	//ok, err := regexp.MatchString(`^(\d+[a-z]+[A-Z]+\W+)$`, pwd)
	//fmt.Println(err)
	return true
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func (u UserMicroservice) Register(ctx context.Context, req *pd.RegisterReq) (*pd.RegisterResp, error) {

	var (
		usr   = req.GetUsername()
		pwd   = req.GetPassword()
		email = req.GetEmail()
	)

	if len(usr) == 0 || !isPasswordValid(pwd) || !isEmailValid(email) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid: missing user or password")
	}

	err := u.userUsecase.Register(ctx, models.User{
		Usr:         usr,
		Pwd:         pwd,
		Email:       email,
		FullName:    req.GetFullName(),
		PhoneNumber: req.GetPhoneNumber(),
	})
	if err != nil {
		return nil, status.Errorf(utils.ParseGRPCErrStatusCode(err), err.Error())
	}

	return &pd.RegisterResp{
		Message: "Register succeed",
	}, nil
}

func (u UserMicroservice) Login(ctx context.Context, req *pd.LoginReq) (*pd.LoginResp, error) {
	usr := req.GetUsername()
	pwd := req.GetPassword()
	if len(usr) == 0 || !isPasswordValid(pwd) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid: missing user or password")
	}

	auth, err := u.userUsecase.Login(ctx, usr, pwd)
	if err != nil {
		return nil, status.Errorf(utils.ParseGRPCErrStatusCode(err), err.Error())
	}

	return &pd.LoginResp{
		Message: "Login succeed",
		Token:   auth,
	}, nil
}

func (u UserMicroservice) UpdateUser(ctx context.Context, req *pd.UpdateUserReq) (*pd.UpdateUserResp, error) {
	token := utils.GetFromMetadata(ctx, "Authorization")
	usr, err := auth.GetUsernameFromToken(token)
	if err != nil {
		return nil, status.Errorf(utils.ParseGRPCErrStatusCode(err), err.Error())
	}

	filter := models.User{
		Usr: usr,
	}

	// TODO: validate update field
	updater := models.User{}

	if err := u.userUsecase.UpdateUser(ctx, filter, updater); err != nil {
		return nil, status.Errorf(utils.ParseGRPCErrStatusCode(err), err.Error())
	}

	return &pd.UpdateUserResp{}, nil
}

func (u UserMicroservice) Withdraw(ctx context.Context, req *pd.WithdrawReq) (*pd.WithdrawResp, error) {
	//TODO implement me
	panic("implement me")
}
