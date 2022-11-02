package service

import (
	"context"
	"github.com/stellarisJAY/goim/pkg/db/dao"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MessageServiceImpl) ListNotifications(_ context.Context, request *pb.ListNotificationRequest) (*pb.ListNotificationResponse, error) {
	var notifications []*model.Notification
	var err error
	if request.NotRead {
		notifications, err = dao.ListNotReadNotifications(request.UserID)
	} else {
		notifications, err = dao.ListAllNotifications(request.UserID)
	}
	if err != nil {
		return &pb.ListNotificationResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	result := make([]*pb.Notification, len(notifications))
	for i, n := range notifications {
		result[i] = &pb.Notification{
			Id:          n.Id,
			Receiver:    n.Receiver,
			TriggerUser: n.TriggerUser,
			Message:     n.Message,
			Read:        n.Read,
			Type:        int32(n.Type),
			Timestamp:   n.Timestamp,
		}
	}
	return &pb.ListNotificationResponse{
		Code:          pb.Success,
		Notifications: result,
	}, nil
}

func (m *MessageServiceImpl) MarkNotificationRead(_ context.Context, request *pb.MarkNotificationReadRequest) (*pb.MarkNotificationReadResponse, error) {
	err := dao.MarkNotificationRead(request.UserID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.MarkNotificationReadResponse{Code: pb.NotFound, Message: "notification not found"}, nil
		} else {
			return &pb.MarkNotificationReadResponse{Code: pb.Error, Message: err.Error()}, nil
		}
	}
	return &pb.MarkNotificationReadResponse{Code: pb.Success}, nil
}

func (m *MessageServiceImpl) AddNotification(_ context.Context, request *pb.AddNotificationRequest) (*pb.AddNotificationResponse, error) {
	notification := model.Notification{
		Id:          request.Notification.Id,
		Receiver:    request.Notification.Receiver,
		TriggerUser: request.Notification.TriggerUser,
		Type:        byte(request.Notification.Type),
		Message:     request.Notification.Message,
		Read:        false,
		Timestamp:   request.Notification.Timestamp,
	}
	err := dao.AddNotification(&notification)
	if err != nil {
		return &pb.AddNotificationResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.AddNotificationResponse{Code: pb.Success}, nil
}

func (m *MessageServiceImpl) GetNotification(ctx context.Context, request *pb.GetNotificationRequest) (*pb.GetNotificationResponse, error) {
	notification, err := dao.GetNotification(request.Id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.GetNotificationResponse{Code: pb.NotFound, Message: "notification not found"}, nil
		}
		return &pb.GetNotificationResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.GetNotificationResponse{Code: pb.Success, Notification: &pb.Notification{
		Id:          notification.Id,
		Receiver:    notification.Receiver,
		TriggerUser: notification.TriggerUser,
		Message:     notification.Message,
		Read:        notification.Read,
		Type:        int32(notification.Type),
		Timestamp:   notification.Timestamp,
	}}, nil
}

func (m *MessageServiceImpl) RemoveNotification(ctx context.Context, request *pb.RemoveNotificationRequest) (*pb.RemoveNotificationResponse, error) {
	err := dao.RemoveNotification(request.Id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.RemoveNotificationResponse{Code: pb.NotFound, Message: "notification not found"}, nil
		}
		return &pb.RemoveNotificationResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	return &pb.RemoveNotificationResponse{Code: pb.Success}, nil
}
