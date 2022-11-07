package service

import (
	"context"
	"fmt"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"google.golang.org/protobuf/proto"
	"time"
)

func (m *MessageServiceImpl) SendMessage(ctx context.Context, request *pb.SendMsgRequest) (*pb.SendMsgResponse, error) {
	message := request.Msg
	// 为消息添加时间戳 和 ID
	message.Timestamp = time.Now().UnixMilli()
	message.Id = m.idGenerator.NextID()

	// 根据聊天类型，检查是否是群成员或是否是好友关系
	switch message.Flag {
	case pb.MessageFlag_From:
		if ok, err := isFriends(message.From, message.To); err != nil {
			return &pb.SendMsgResponse{Code: pb.Error, Message: err.Error()}, nil
		} else if !ok {
			return &pb.SendMsgResponse{Code: pb.AccessDenied, Message: "can't send message to stranger"}, nil
		}
	case pb.MessageFlag_Group:
		if isMember, normal, err := isGroupMember(message.From, message.To); err != nil {
			return &pb.SendMsgResponse{Code: pb.Error, Message: err.Error()}, nil
		} else if !isMember {
			return &pb.SendMsgResponse{Code: pb.AccessDenied, Message: "not a member in this group chat"}, nil
		} else if !normal {
			return &pb.SendMsgResponse{Code: pb.AccessDenied, Message: "banned to speak in this group chat"}, nil
		}
	default:
		return &pb.SendMsgResponse{Code: pb.Error, Message: "unknown message flag"}, nil
	}
	// 过滤敏感词
	_, replaced := m.wordFilter.DoFilter(message.Content)
	message.Content = replaced
	// 序列化
	marshal, err := proto.Marshal(message)
	if err != nil {
		return &pb.SendMsgResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	// 发送给推送服务
	key := fmt.Sprintf("%x", message.Id)
	err = m.producer.PushMessage(pb.MessageTransferTopic, key, marshal)
	if err != nil {
		return &pb.SendMsgResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	// 回复时带上替换后的消息内容
	return &pb.SendMsgResponse{
		Code:      pb.Success,
		MessageId: message.Id,
		Timestamp: message.Timestamp,
		Content:   replaced,
	}, nil
}

// isFriends 检查消息双方是否是好友关系
func isFriends(from, to int64) (bool, error) {
	service, err := GetFriendService()
	if err != nil {
		return false, nil
	}
	response, err := service.CheckFriendship(context.TODO(), &pb.FriendshipRequest{FriendID: to, UserID: from})
	if err != nil {
		return false, err
	}
	switch response.Code {
	case pb.Error:
		return false, fmt.Errorf("get friendship error %s", response.Message)
	case pb.Success:
		return response.IsFriend, nil
	default:
		return false, nil
	}
}

// isGroupMember 检查用户是否是群成员，是否能够在群聊中发言
func isGroupMember(userID, groupID int64) (bool, bool, error) {
	service, err := GetGroupService()
	if err != nil {
		return false, false, err
	}
	response, err := service.GetGroupMember(context.TODO(), &pb.GetGroupMemberRequest{GroupID: groupID, UserID: userID})
	if err != nil {
		return false, false, err
	}
	switch response.Code {
	case pb.Error:
		return false, false, fmt.Errorf("get group member info error %s", response.Message)
	case pb.NotFound:
		return false, false, nil
	case pb.Success:
		return true, response.Member.Status == pb.GroupMemberStatus_normal, nil
	}
	return false, false, nil
}

func GetFriendService() (pb.FriendClient, error) {
	conn, err := naming.GetClientConn("friend")
	if err != nil {
		return nil, err
	}
	return pb.NewFriendClient(conn), nil
}

func GetGroupService() (pb.GroupClient, error) {
	conn, err := naming.GetClientConn("group")
	if err != nil {
		return nil, err
	}
	return pb.NewGroupClient(conn), nil
}
