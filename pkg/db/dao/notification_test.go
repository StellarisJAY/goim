package dao

import (
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/snowflake"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"testing"
	"time"
)

var notificationTestID = snowflake.NewSnowflake(10)

func TestAddNotification(t *testing.T) {
	id := notificationTestID.NextID()
	err := AddNotification(&model.Notification{
		Id:          id,
		Receiver:    1001,
		TriggerUser: 1002,
		Type:        model.NotificationFriendRequest,
		Message:     "I am user_1001",
		Read:        false,
		Timestamp:   time.Now().UnixMilli(),
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	notification, err := GetNotification(id)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if notification.Receiver != 1001 || notification.TriggerUser != 1002 {
		t.Log("wrong notification after add")
		t.FailNow()
	}
}

func TestListNotificationOfType(t *testing.T) {
	notifications := make(map[int64]*model.Notification)
	for i := 1; i <= 10; i++ {
		id := notificationTestID.NextID()
		notifications[id] = &model.Notification{
			Id:          id,
			Receiver:    1999,
			TriggerUser: int64(1000 + i),
			Type:        model.NotificationGroupInvitation,
			Message:     "add friend",
			Timestamp:   time.Now().UnixMilli(),
		}
	}
	for _, n := range notifications {
		err := AddNotification(n)
		if err != nil && !mongo.IsDuplicateKeyError(err) {
			t.Errorf("add notification error %v", err)
			t.FailNow()
		}
	}
	// 删除测试数据
	defer func() {
		for id, _ := range notifications {
			_ = RemoveNotification(id)
		}
	}()
	result, err := ListNotificationOfType(1999, model.NotificationGroupInvitation)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, n := range result {
		if original, ok := notifications[n.Id]; !ok {
			t.Logf("list operation id not found %d, n.receiver: %d", n.Id, n.Receiver)
			t.FailNow()
		} else {
			if !reflect.DeepEqual(original, n) {
				t.Logf("notification not equal")
				t.FailNow()
			}
		}
	}
}

func TestRemoveNotification(t *testing.T) {
	id := notificationTestID.NextID()
	addNote := &model.Notification{
		Id:          id,
		Receiver:    1100,
		TriggerUser: 1102,
		Type:        model.NotificationFriendRequest,
		Message:     "add friend",
		Read:        false,
		Timestamp:   time.Now().UnixMilli(),
	}
	err := AddNotification(addNote)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	notification, err := GetNotification(id)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if !reflect.DeepEqual(addNote, notification) {
		t.Errorf("get after add, wrong notification")
		t.FailNow()
	}
	err = RemoveNotification(id)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = GetNotification(id)
	if err == nil || err != mongo.ErrNoDocuments {
		t.FailNow()
	}
}

func TestListAllNotifications(t *testing.T) {
	notifications := make(map[int64]*model.Notification)
	for i := 1; i <= 10; i++ {
		id := notificationTestID.NextID()
		notifications[id] = &model.Notification{
			Id:          id,
			Receiver:    1999,
			TriggerUser: int64(1000 + i),
			Type:        model.NotificationFriendRequest,
			Message:     "add friend",
			Timestamp:   time.Now().UnixMilli(),
		}
	}
	for _, n := range notifications {
		err := AddNotification(n)
		if err != nil && !mongo.IsDuplicateKeyError(err) {
			t.Errorf("add notification error %v", err)
			t.FailNow()
		}
	}
	// 删除测试数据
	defer func() {
		for id, _ := range notifications {
			_ = RemoveNotification(id)
		}
	}()
	result, err := ListAllNotifications(1999)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, n := range result {
		if original, ok := notifications[n.Id]; !ok {
			t.Logf("list operation id not found %d, n.receiver: %d", n.Id, n.Receiver)
			t.FailNow()
		} else {
			if !reflect.DeepEqual(original, n) {
				t.Logf("notification not equal")
				t.FailNow()
			}
		}
	}
}

func TestListNotReadNotifications(t *testing.T) {
	notifications := make(map[int64]*model.Notification)
	for i := 1; i <= 10; i++ {
		id := notificationTestID.NextID()
		notifications[id] = &model.Notification{
			Id:          id,
			Receiver:    1999,
			TriggerUser: int64(1000 + i),
			Type:        model.NotificationFriendRequest,
			Message:     "add friend",
			Read:        i > 5,
			Timestamp:   time.Now().UnixMilli(),
		}
	}
	for _, n := range notifications {
		err := AddNotification(n)
		if err != nil && !mongo.IsDuplicateKeyError(err) {
			t.Errorf("add notification error %v", err)
			t.FailNow()
		}
	}
	// 删除测试数据
	defer func() {
		for id, _ := range notifications {
			_ = RemoveNotification(id)
		}
	}()
	result, err := ListNotReadNotifications(1999)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, n := range result {
		if original, ok := notifications[n.Id]; !ok {
			t.Logf("list operation id not found %d, n.receiver: %d", n.Id, n.Receiver)
			t.FailNow()
		} else {
			if !reflect.DeepEqual(original, n) {
				t.Logf("notification not equal")
				t.FailNow()
			}
		}
	}
}
