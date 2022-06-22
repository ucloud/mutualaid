package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"

	biz "github.com/ucloud/mutualaid/backend/internal/biz/aid"
	userUc "github.com/ucloud/mutualaid/backend/internal/biz/user"

	pb "github.com/ucloud/mutualaid/backend/api/mutualaid"
)

type ExamineAidService struct {
	pb.UnimplementedExamineAidServer
	uc          biz.UseCase
	userUsecase *userUc.UserUsecase
}

func NewExamineAidService(uc biz.UseCase, userUsecase *userUc.UserUsecase) *ExamineAidService {
	return &ExamineAidService{uc: uc, userUsecase: userUsecase}

}

func (s *ExamineAidService) ExamineLogin(ctx context.Context, req *pb.ExamineLoginReq) (*pb.ExamineLoginResp, error) {
	jwt, err := s.uc.ExamineLogin(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	if tr, ok := transport.FromServerContext(ctx); ok {
		tr.RequestHeader().Set("JWT-Token", jwt)
	}
	return &pb.ExamineLoginResp{}, nil
}

func (s *ExamineAidService) GetExamineList(ctx context.Context, req *pb.GetExamineListReq) (*pb.ExamineListResp, error) {
	list, totalArray, err := s.uc.GetExamineList(ctx, int32(req.ExamineStatus), string(req.ExamineStatusOrder), string(req.CreateTimeOrder), string(req.UpdateTimeOrder), int32(req.PageNumber), int32(req.PageSize), string(req.VagueSearch))
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
			ExamineStatus: a.ExamineStatus,
			CreateTime:    a.CreateTime,
			Status:        a.Status,
			UpdateTime:    a.UpdateTime,
			MessageCount:  a.MessageCount,
			UserId:        int64(a.UserID),
			Address:       a.Address,
		}

		if a.UserInfo != nil {
			aid.User = &pb.UserInfo{
				Name: a.UserInfo.Name,
				Icon: a.UserInfo.ICon}
		}
		getAid = append(getAid, aid)
	}

	// 统计各数值
	return &pb.ExamineListResp{
		List:        getAid,
		TotalSize:   findValue(totalArray, 0),
		WaitingSize: findValue(totalArray, biz.StatusExamineWait),
		PassSize:    findValue(totalArray, biz.StatusExamineFinish),
		BlockSize:   findValue(totalArray, biz.StatusExamineBlock),
	}, nil
}

func (s *ExamineAidService) ExamineAid(ctx context.Context, req *pb.ExamineAidReq) (*pb.ExamineAidResp, error) {
	err := s.uc.ExamineAid(ctx, uint64(req.Id), string(req.ExamineAction))

	if err != nil {
		return nil, err
	}

	return &pb.ExamineAidResp{}, nil
}

func (s *ExamineAidService) GetBlockUserList(ctx context.Context, req *pb.GetBlockUserListReq) (*pb.GetBlockUserListResp, error) {
	// 用了userUserCase不知道是否有啥风险
	uList, err := s.userUsecase.ListUser(ctx, "", biz.UserIsBlock)
	if err != nil {
		return nil, err
	}
	getUserList := make([]*pb.BlockUserInfo, 0, len(uList))

	for _, a := range uList {
		thisUser := &pb.BlockUserInfo{
			Id:         int64(a.ID),
			Addr:       a.Addr,
			Name:       a.Name,
			Phone:      a.Phone,
			Icon:       a.Icon,
			CreateTime: a.CreateTime,
			Status:     int64(a.Status),
		}

		getUserList = append(getUserList, thisUser)
	}

	return &pb.GetBlockUserListResp{List: getUserList}, nil
}

// BlockUser 封禁用户接口
func (s *ExamineAidService) BlockUser(ctx context.Context, req *pb.BlockUserReq) (*pb.BlockUserResp, error) {
	// 确定用户是否存在

	u, err := s.userUsecase.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 执行封禁
	u.Status = biz.UserIsBlock
	updateUserErr := s.userUsecase.UpdateUser(ctx, u)
	if updateUserErr != nil {
		return nil, err
	}

	// 执行对求助与求助消息的封禁
	updateAidErr := s.uc.UpdateAllAidAndMessageExamineStatus(ctx, uint64(req.Id), int32(biz.StatusExamineBlock))
	if updateAidErr != nil {
		return nil, err
	}
	return &pb.BlockUserResp{}, nil
}

// PassUser 解封用户接口
func (s *ExamineAidService) PassUser(ctx context.Context, req *pb.PassUserReq) (*pb.PassUserResp, error) {
	// 确定用户是否存在
	u, err := s.userUsecase.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	// 执行解封
	u.Status = 2
	updateUserErr := s.userUsecase.UpdateUser(ctx, u)
	if updateUserErr != nil {
		return nil, err
	}
	// 执行对求助与求助消息的封禁

	updateAidErr := s.uc.UpdateAllAidAndMessageExamineStatus(ctx, uint64(req.Id), int32(biz.StatusExamineWait))
	if updateAidErr != nil {
		return nil, err
	}
	return &pb.PassUserResp{}, nil
}

// 根据入参Map与需要得到的类型，给出计数
func findValue(typeList []biz.ExamineTypeMap, typeId int) (ret int64) {
	for i := 0; i < len(typeList); i++ {
		if typeList[i].ExamineStatus == typeId {
			return typeList[i].Count
		}
	}
	return 0
}
