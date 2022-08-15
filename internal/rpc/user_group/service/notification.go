package service

import (
	"context"
	"fmt"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/snowflake"
	"time"
)

var notificationId = snowflake.NewSnowflake(config.Config.MachineID)

// Notification 向某个用户发送一条通知
func Notification(userID int64, from int64, flag pb.MessageFlag, message string) error {
	conn, err := naming.GetClientConn("message")
	if err != nil {
		return err
	}
	client := pb.NewMessageClient(conn)
	response, err := client.SendMessage(context.TODO(), &pb.SendMsgRequest{Msg: &pb.BaseMsg{
		From:      from,
		To:        userID,
		Content:   message,
		Flag:      flag,
		Timestamp: time.Now().UnixMilli(),
		Id:        notificationId.NextID(),
	}})
	if err != nil {
		return fmt.Errorf("send notification to user: %s failed, error: %w", userID, err)
	}
	if response.Code != pb.Success {
		return fmt.Errorf("send notification to user: %s failed: %s", userID, response.Message)
	}
	return nil
}
