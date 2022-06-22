package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/ucloud/mutualaid/backend/infra/userutil"
	"github.com/ucloud/mutualaid/backend/internal/biz"
	"github.com/ucloud/mutualaid/backend/internal/biz/aid"
	"github.com/ucloud/mutualaid/backend/internal/biz/user"

	pb "github.com/ucloud/mutualaid/backend/api/mutualaid"
)

type UserAidQueryService struct {
	pb.UnimplementedUserAidQueryServer

	uc          aid.UseCase
	userUsecase *user.UserUsecase
	logger      *log.Helper
}

func NewUserAidQueryService(uc aid.UseCase, userUsecase *user.UserUsecase, logger log.Logger) *UserAidQueryService {
	return &UserAidQueryService{
		uc:          uc,
		userUsecase: userUsecase,
		logger:      log.NewHelper(logger),
	}
}

func (s *UserAidQueryService) ListAidOffered(ctx context.Context, req *pb.ListAidOfferedReq) (*pb.ListAidOfferedResp, error) {
	offset, limit := pageArgs(req.PageNumber, req.PageSize)

	aidlist, total, err := s.uc.ListUserAid(ctx, userutil.ExtractUID(ctx), true, offset, limit)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	resp := &pb.ListAidOfferedResp{TotalSize: total}
	for _, a := range aidlist {
		u, err := s.userUsecase.GetUser(ctx, a.UserID)
		if err != nil {
			s.logger.Warnf("not found user info of %d", a.UserID)
			continue
		}
		pbaid := s.aid2pb(ctx, a, float64(req.Latitude), float64(req.Longitude))
		pbaid.User = &pb.UserInfo{
			Name: u.Name,
			Icon: u.Icon,
			Addr: a.Address,
		}
		pbaid.DisplayPim = true

		resp.List = append(resp.List, pbaid)
	}
	return resp, nil
}
func (s *UserAidQueryService) ListAidNeeds(ctx context.Context, req *pb.ListAidNeedsReq) (*pb.ListAidNeedsResp, error) {
	offset, limit := pageArgs(req.PageNumber, req.PageSize)

	aidlist, total, err := s.uc.ListUserAid(ctx, userutil.ExtractUID(ctx), false, offset, limit)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	resp := &pb.ListAidNeedsResp{TotalSize: total}
	for _, a := range aidlist {
		u, err := s.userUsecase.GetUser(ctx, a.UserID)
		if err != nil {
			s.logger.Warnf("not found user info of %d", a.UserID)
			continue
		}
		pbaid := s.aid2pb(ctx, a, float64(req.Latitude), float64(req.Longitude))
		pbaid.User = &pb.UserInfo{
			Name: u.Name,
			Icon: u.Icon,
			Addr: a.Address,
		}
		pbaid.DisplayPim = true

		resp.List = append(resp.List, pbaid)
	}
	return resp, nil
}
func (s *UserAidQueryService) GetAidDetail(ctx context.Context, req *pb.GetAidDetailReq) (*pb.GetAidDetailResp, error) {
	a, isMyAid, isMyHelp, err := s.uc.GetAidDetail(ctx, req.Id, userutil.ExtractUID(ctx))
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	pbaid := s.aid2pb(ctx, a, float64(req.Latitude), float64(req.Longitude))
	pbaid.Message = s.msg2pb(ctx, a, isMyAid)

	u, err := s.userUsecase.GetUser(ctx, a.UserID)
	if err != nil {
		s.logger.Errorf("not found user info of %d, %+v", a.UserID, err)
		return nil, err
	}

	pbaid.User = &pb.UserInfo{
		Name: u.Name,
		Icon: u.Icon,
	}

	// 检查当前用户和信息之间的关系
	if isMyAid || isMyHelp {
		pbaid.User.Phone = a.Phone
		pbaid.User.Addr = a.Address
		pbaid.DisplayPim = true
	}

	resp := &pb.GetAidDetailResp{
		Aid:      pbaid,
		IsMyAid:  isMyAid,
		IsMyHelp: isMyHelp,
	}
	return resp, nil
}

//-------
func (s *UserAidQueryService) aid2pb(ctx context.Context, a *aid.Aid, lat, long float64) *pb.Aid {
	pbaid := &pb.Aid{
		Id:           a.ID,
		Type:         a.Type,
		Group:        a.Group,
		Emergency:    a.EmergencyLevel,
		Content:      a.Content,
		Distance:     int64(biz.GetDistanceReturnMeter(lat, long, a.Latitude, a.Longitude)),
		CreateTime:   a.CreateTime,
		Status:       a.Status,
		MessageCount: a.MessageCount,
		Address:      a.Address,
	}

	return pbaid
}

func (s *UserAidQueryService) msg2pb(ctx context.Context, a *aid.Aid, isMyAid bool) []*pb.Message {
	var msgList []*pb.Message
	for _, msg := range a.Messages {
		m := &pb.Message{
			Id:         msg.ID,
			MaskPhone:  msg.UserPhone,
			Content:    msg.Content,
			CreateTime: msg.CreateTime,
			Status:     (msg.Status),
		}

		u, err := s.userUsecase.GetUser(ctx, msg.UserID)
		if err != nil {
			s.logger.Errorf("not found user info of %d, %+v", a.UserID, err)
			return msgList
		}
		m.User = &pb.UserInfo{
			Name: u.Name,
			Icon: u.Icon,
		}
		m.DisplayPim = isMyAid || msg.UserID == userutil.ExtractUID(ctx)
		if m.DisplayPim {
			m.User.Phone = msg.UserPhone
		}

		msgList = append(msgList, m)
	}

	return msgList
}

func pageArgs(pageNo, pageSize int64) (int64, int64) {
	offset := pageNo * pageSize
	if pageNo < 0 {
		offset = 0
	}
	limit := pageSize

	return offset, limit
}
