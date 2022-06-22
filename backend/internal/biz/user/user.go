package user

import (
	"context"
	"errors"

	userutil2 "github.com/ucloud/mutualaid/backend/infra/userutil"

	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/ucloud/mutualaid/backend/api/mutualaid"
	"github.com/ucloud/mutualaid/backend/internal/data/model"
	"github.com/ucloud/mutualaid/backend/internal/proxy/wechat"
)

type UserRepo interface {
	GetUser(context.Context, *model.User) (*model.User, error)
	CreateUser(context.Context, *model.User) (*model.User, error)
	UpdateUser(context.Context, *model.User) error
	DeleteUser(context.Context, *model.User) error
	ListUser(context.Context, string, int8) ([]*model.User, error)
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUsecase) WxOAuth2(ctx context.Context, req *pb.WxOAuth2Req) (string, error) {
	u := &model.User{}
	if req.Openid != "" {
		unionID, err := wechat.UserInfo(req.Openid)
		if err != nil {
			return "", err
		}
		u.MpOpenid = req.Openid
		u.Unionid = unionID
	} else if req.Code != "" {
		openID, oauthToken, err := wechat.OAuth2(req.Code)
		if err != nil {
			return "", err
		}

		unionID, nickname, headImgUrl, err := wechat.SnsUserInfo(openID, oauthToken)
		if err != nil {
			return "", err
		}

		u.MpOpenid = openID
		u.Unionid = unionID
		u.Name = nickname
		u.Icon = headImgUrl
	} else {
		return "", errors.New("无法完成微信授权")
	}

	u, err := uc.repo.CreateUser(ctx, u)
	if err != nil {
		return "", err
	}

	jwtToken, err := userutil2.NewJWT().Auth(u.ID, u.Openid)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (uc *UserUsecase) WxLogin(ctx context.Context, req *pb.WxLoginReq) (string, bool, error) {
	// 调用微信接口取得openid
	openID, unionID, _, err := wechat.Login(req.LoginCode)
	if err != nil {
		return "", false, err
	}

	u := &model.User{
		Openid:  openID,
		Unionid: unionID,
		Name:    req.Name,
	}
	u, err = uc.repo.CreateUser(ctx, u)
	if err != nil {
		return "", false, err
	}

	jwtToken, err := userutil2.NewJWT().Auth(u.ID, u.Openid)
	if err != nil {
		return "", false, err
	}

	needActive := u.Phone == ""

	return jwtToken, needActive, nil
}

func (uc *UserUsecase) WxPhoneNumber(ctx context.Context, req *pb.WxPhoneNumberReq) (string, error) {
	_, purePhoneNumber, err := wechat.GetPhone(req.PhoneCode)
	if err != nil {
		return "", err
	}

	u, err := uc.repo.GetUser(ctx, &model.User{Openid: userutil2.ExtractOpenID(ctx)})
	if err != nil {
		return "", err
	}
	u.Phone = purePhoneNumber
	err = uc.repo.UpdateUser(ctx, u)
	if err != nil {
		return "", err
	}

	return purePhoneNumber, nil
}

func (uc *UserUsecase) Active(ctx context.Context, req *pb.ActiveUserReq) (string, error) {
	// 调用微信接口取得openid
	openID, unionID, _, err := wechat.Login(req.LoginCode)
	if err != nil {
		return "", err
	}
	if openID == "" {
		return "", errors.New("logincode 无效，无法取得openid")
	}

	_, purePhoneNumber, err := wechat.GetPhone(req.PhoneCode)
	if err != nil {
		return "", err
	}

	u := &model.User{
		Openid:  openID,
		Unionid: unionID,
		Name:    req.Name,
		Icon:    req.Icon,
	}
	if purePhoneNumber != "" {
		u.Phone = purePhoneNumber
	}
	u, err = uc.repo.CreateUser(ctx, u)
	if err != nil {
		return "", err
	}

	jwtToken, err := userutil2.NewJWT().Auth(u.ID, u.Openid)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (uc *UserUsecase) GetUser(ctx context.Context, user_id uint64) (*model.User, error) {
	return uc.repo.GetUser(ctx, &model.User{ID: user_id})
}

func (uc *UserUsecase) ListUser(ctx context.Context, community string, status int8) ([]*model.User, error) {
	return uc.repo.ListUser(ctx, community, status)
}
func (uc *UserUsecase) UpdateUser(ctx context.Context, u *model.User) error {
	return uc.repo.UpdateUser(ctx, u)
}

func (uc *UserUsecase) JSAPISign(ctx context.Context, req *pb.JSAPISignReq) (*pb.JSAPISignResp, error) {
	nonce, ts, sign := wechat.JSSign(req.ApiUrl)
	return &pb.JSAPISignResp{
		Noncestr:  nonce,
		Timestamp: ts,
		Sign:      sign,
	}, nil
}
