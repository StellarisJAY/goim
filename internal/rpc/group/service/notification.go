package service

import (
	"context"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/snowflake"
	"go.uber.org/zap"
	"time"
)

var notificationId = snowflake.NewSnowflake(config.Config.MachineID)
var notificationChan = make(chan *pb.Notification, 1024)

func NotifyUser(userID int64, triggerUser int64, nType byte, message string) {
	notification := &pb.Notification{
		Id:          notificationId.NextID(),
		Receiver:    userID,
		TriggerUser: triggerUser,
		Message:     message,
		Read:        false,
		Type:        int32(nType),
		Timestamp:   time.Now().UnixMilli(),
	}
	notificationChan <- notification
}

func AsyncPushSendNotification() {
	for notification := range notificationChan {
		go func(notification *pb.Notification) {
			conn, err := naming.GetClientConn("message")
			if err != nil {
				log.Warn("get notification service failed", zap.Error(err))
				return
			}
			service := pb.NewMessageClient(conn)
			response, err := service.AddNotification(context.TODO(), &pb.AddNotificationRequest{Notification: notification})
			if err != nil || response.Code == pb.Error {
				log.Warn("push notification error", zap.Error(err))
			}
		}(notification)
	}
}
