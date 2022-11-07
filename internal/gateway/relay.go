package gateway

import (
	"context"
	"encoding/json"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// PushMessage 消息下行 RPC 服务
func (s *Server) PushMessage(ctx context.Context, request *pb.PushMsgRequest) (*pb.PushMsgResponse, error) {
	// 获取目标 channel
	load, ok := s.wsServer.Channels.Load(request.Channel)
	resp := new(pb.PushMsgResponse)
	if ok {
		// 将 message 编码后向channel发送
		channel := load.(*websocket.Channel)
		marshal, err := marshalMessage(request.Message)
		if err != nil {
			return nil, err
		}
		err = channel.Push(marshal)
		if err != nil {
			return nil, err
		}
		resp.Base = new(pb.BaseResponse)
		resp.Base.Code = pb.Success
		return resp, nil
	} else {
		resp.Base.Code = pb.NotFound
		resp.Base.Message = "channel not found"
		return resp, nil
	}
}

// Broadcast 消息广播下行，向多个channel发送相同的消息
func (s *Server) Broadcast(ctx context.Context, request *pb.BroadcastRequest) (*pb.BroadcastResponse, error) {
	channelIds := request.Channels
	channels := make([]*websocket.Channel, 0, len(channelIds))
	fails := make([]string, 0, len(channelIds))
	// 获取channel对象，如果channel不存在或不可用则将ID加入到失败列表
	for _, id := range channelIds {
		if c, ok := s.wsServer.Channels.Load(id); ok {
			channel := c.(*websocket.Channel)
			if channel.Available() {
				channels = append(channels, channel)
				continue
			}
		}
		fails = append(fails, id)
	}
	// 编码消息
	marshal, err := marshalMessage(request.Message)
	if err != nil {
		return nil, err
	}
	// 向每个channel推送消息
	for _, channel := range channels {
		err := channel.Push(marshal)
		if err != nil {
			// 推送失败，将channelID添加到失败列表
			log.Warn("push message to channel failed",
				zap.String("channel", channel.ID()),
				zap.Error(err))
			fails = append(fails, channel.ID())
		}
	}
	return &pb.BroadcastResponse{Fails: fails}, nil
}

func marshalMessage(message *pb.BaseMsg) ([]byte, error) {
	if config.Config.Gateway.UseJsonMsg {
		return json.Marshal(message)
	} else {
		return proto.Marshal(message)
	}
}
