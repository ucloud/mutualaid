package model

import biz "github.com/ucloud/mutualaid/backend/internal/biz/aid"

func AidDO2PO(do biz.Aid) *Aid {
	po := &Aid{
		ID:             do.ID,
		UserID:         do.UserID,
		Type:           do.Type,
		Group:          do.Group,
		EmergencyLevel: do.EmergencyLevel,
		ExamineStatus:  do.ExamineStatus,
		Status:         do.Status,
		FinishUserID:   do.FinishUserID,
		FinishTime:     do.FinishTime,
		MessageCount:   do.MessageCount,
		Content:        do.Content,
		Longitude:      do.Longitude,
		Latitude:       do.Latitude,
		Phone:          do.Phone,
		District:       do.District,
		Address:        do.Address,
		CreateTime:     do.CreateTime,
		Version:        do.Version,
	}
	return po
}

func AidPO2DO(po Aid) *biz.Aid {
	do := &biz.Aid{
		ID:             po.ID,
		UserID:         po.UserID,
		Type:           po.Type,
		Group:          po.Group,
		EmergencyLevel: po.EmergencyLevel,
		Status:         po.Status,
		ExamineStatus:  po.ExamineStatus,
		FinishUserID:   po.FinishUserID,
		FinishTime:     po.FinishTime,
		MessageCount:   po.MessageCount,
		Content:        po.Content,
		Longitude:      po.Longitude,
		Latitude:       po.Latitude,
		Phone:          po.Phone,
		District:       po.District,
		Address:        po.Address,
		CreateTime:     po.CreateTime,
		Version:        po.Version,
		Messages:       nil,
		UserInfo:       nil,
	}
	return do
}

func AidWithUserPO2DO(po AidWithUser) *biz.Aid {
	do := &biz.Aid{
		ID:             po.ID,
		UserID:         po.UserID,
		Type:           po.Type,
		Group:          po.Group,
		EmergencyLevel: po.EmergencyLevel,
		Status:         po.Status,
		ExamineStatus:  po.ExamineStatus,
		FinishUserID:   po.FinishUserID,
		FinishTime:     po.FinishTime,
		MessageCount:   po.MessageCount,
		Content:        po.Content,
		Longitude:      po.Longitude,
		Latitude:       po.Latitude,
		Phone:          po.Phone,
		District:       po.District,
		Address:        po.Address,
		CreateTime:     po.CreateTime,
		UpdateTime:     po.UpdateTime,
		Version:        po.Version,
		UserInfo: &biz.UserInfo{
			Name: po.Name,
			ICon: po.Icon,
		},
	}
	return do
}

func AidMessageDO2PO(do biz.Message) *AidMessage {
	po := &AidMessage{
		ID:         do.ID,
		AidID:      do.AidID,
		Status:     do.Status,
		UserID:     do.UserID,
		UserPhone:  do.UserPhone,
		Content:    do.Content,
		CreateTime: do.CreateTime,
		Version:    do.Version,
	}
	return po
}

func AidMessagePO2DO(po AidMessage) *biz.Message {
	do := &biz.Message{
		ID:         po.ID,
		AidID:      po.AidID,
		Status:     po.Status,
		UserID:     po.UserID,
		UserPhone:  po.UserPhone,
		Content:    po.Content,
		CreateTime: po.CreateTime,
		Version:    po.Version,
	}
	return do
}
