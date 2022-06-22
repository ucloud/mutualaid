package service

import (
	"context"
	"github.com/ucloud/mutualaid/backend/infra/userutil"
	biz "github.com/ucloud/mutualaid/backend/internal/biz/aid"

	pb "github.com/ucloud/mutualaid/backend/api/mutualaid"
)

type MutualAidQueryService struct {
	pb.UnimplementedMutualAidQueryServer
	uc biz.UseCase
}

func NewMutualAidQueryService(uc biz.UseCase) *MutualAidQueryService {
	return &MutualAidQueryService{uc: uc}
}

func (s *MutualAidQueryService) Discovery(ctx context.Context, req *pb.DiscoveryReq) (*pb.DiscoveryResp, error) {
	list, err := s.uc.Discovery(ctx, userutil.ExtractUID(ctx), float64(req.Latitude), float64(req.Longitude), req.PageNumber, req.PageSize)
	if err != nil {
		return nil, err
	}

	getAid := make([]*pb.Aid, 0, len(list))
	for _, a := range list {
		aid := &pb.Aid{
			Id:            a.ID,
			Type:          a.Type,
			Group:         a.Group,
			Emergency:     a.EmergencyLevel,
			Content:       a.Content,
			Distance:      a.Distance,
			CreateTime:    a.CreateTime,
			Status:        a.Status,
			ExamineStatus: a.ExamineStatus,
			Address:       a.Address,
			MessageCount:  a.MessageCount,
		}
		if a.UserInfo != nil {
			aid.User = &pb.UserInfo{
				Name: a.UserInfo.Name,
				Icon: a.UserInfo.ICon,
			}
		}
		getAid = append(getAid, aid)
	}

	return &pb.DiscoveryResp{
		List:      getAid,
		TotalSize: int64(len(getAid)),
	}, nil
}
