package service

import (
	"context"
	"github.com/ucloud/mutualaid/backend/infra/userutil"
	biz "github.com/ucloud/mutualaid/backend/internal/biz/aid"

	pb "github.com/ucloud/mutualaid/backend/api/mutualaid"
)

type UserAidManagerService struct {
	pb.UnimplementedUserAidManagerServer
	uc biz.UseCase
}

func NewUserAidManagerService(uc biz.UseCase) *UserAidManagerService {
	return &UserAidManagerService{uc: uc}
}

func (s *UserAidManagerService) CreateAid(ctx context.Context, req *pb.CreateAidReq) (*pb.CreateAidResp, error) {
	aid := biz.NewAid(userutil.ExtractUID(ctx), int32(req.Type.Number()), int32(req.Group.Number()), int32(req.Emergency.Number()), req.Content, float64(req.Longitude), float64(req.Latitude), "", req.Addr, req.Phone)
	if err := s.uc.CreateAid(ctx, *aid); err != nil {
		return nil, err
	}
	return &pb.CreateAidResp{Id: aid.ID}, nil
}
func (s *UserAidManagerService) CancelAid(ctx context.Context, req *pb.CancelAidReq) (*pb.CancelAidResp, error) {
	if err := s.uc.CancelAid(ctx, req.Id, userutil.ExtractUID(ctx)); err != nil {
		return nil, err
	}
	return &pb.CancelAidResp{}, nil
}
func (s *UserAidManagerService) FinishAid(ctx context.Context, req *pb.FinishAidReq) (*pb.FinishAidResp, error) {
	if err := s.uc.FinishAid(ctx, req.Id, userutil.ExtractUID(ctx), req.MessageId); err != nil {
		return nil, err
	}

	return &pb.FinishAidResp{}, nil
}
func (s *UserAidManagerService) CreateAidMessage(ctx context.Context, req *pb.CreateAidMessageReq) (*pb.CreateAidMessageResp, error) {
	msg := biz.NewMessage(req.Id, userutil.ExtractUID(ctx), req.Phone, req.Content)
	if err := s.uc.CreateAidMessage(ctx, *msg); err != nil {
		return nil, err
	}
	return &pb.CreateAidMessageResp{}, nil
}
