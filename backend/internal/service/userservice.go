package service

import (
	"context"
	"errors"

	pb "github.com/ucloud/mutualaid/backend/api/mutualaid"
	"github.com/ucloud/mutualaid/backend/infra/userutil"
	"github.com/ucloud/mutualaid/backend/internal/biz/user"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

type UserService struct {
	pb.UnimplementedUserServiceServer

	userUsecase *user.UserUsecase
	logger      *log.Helper
}

func NewUserService(userUsecase *user.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		logger:      log.NewHelper(logger),
		userUsecase: userUsecase,
	}
}

func (s *UserService) WxOAuth2(ctx context.Context, req *pb.WxOAuth2Req) (*pb.WxOAuth2Resp, error) {
	jwt, err := s.userUsecase.WxOAuth2(ctx, req)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if tr, ok := transport.FromServerContext(ctx); ok {
		tr.RequestHeader().Set("JWT-Token", jwt)
	}
	return &pb.WxOAuth2Resp{}, nil
}

func (s *UserService) WxLogin(ctx context.Context, req *pb.WxLoginReq) (*pb.WxLoginResp, error) {
	if req.LoginCode == "" {
		return nil, errors.New("报歉，未登录不能提供服务")
	}

	jwt, needActive, err := s.userUsecase.WxLogin(ctx, req)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if tr, ok := transport.FromServerContext(ctx); ok {
		tr.RequestHeader().Set("JWT-Token", jwt)
	}
	return &pb.WxLoginResp{NeedActive: needActive}, nil
}

func (s *UserService) WxPhoneNumber(ctx context.Context, req *pb.WxPhoneNumberReq) (*pb.WxPhoneNumberResp, error) {
	if req.PhoneCode == "" {
		return nil, errors.New("请授权获取您的手机号")
	}

	phone, err := s.userUsecase.WxPhoneNumber(ctx, req)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &pb.WxPhoneNumberResp{Phone: phone}, nil
}

func (s *UserService) ActiveUser(ctx context.Context, req *pb.ActiveUserReq) (*pb.ActiveUserResp, error) {
	//if req.PhoneCode == "" {
	//	return nil, errors.New("报歉，必须授权手机号才能提供服务")
	//}

	jwt, err := s.userUsecase.Active(ctx, req)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if tr, ok := transport.FromServerContext(ctx); ok {
		tr.RequestHeader().Set("JWT-Token", jwt)
	}

	return &pb.ActiveUserResp{}, nil
}
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserResp, error) {
	u, err := s.userUsecase.GetUser(ctx, userutil.ExtractUID(ctx))
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &pb.GetUserResp{
		User: &pb.UserInfo{
			Name:  u.Name,
			Phone: u.Phone,
			Icon:  u.Icon,
		},
	}, nil
}
func (s *UserService) JSAPISign(ctx context.Context, req *pb.JSAPISignReq) (*pb.JSAPISignResp, error) {
	return s.userUsecase.JSAPISign(ctx, req)
}
