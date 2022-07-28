package service

import (
	"context"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
)

type MessageServiceImpl struct {
}

func (m *MessageServiceImpl) SyncOfflineMessages(ctx context.Context, request *pb.SyncMsgRequest) (*pb.SyncMsgResponse, error) {
	//TODO implement me
	panic("implement me")
}
