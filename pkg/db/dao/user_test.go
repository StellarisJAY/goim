package dao

import (
	"github.com/stellarisJAY/goim/pkg/db/model"
	"testing"
	"time"
)

var testUser1 = &model.User{
	ID:           1001001001,
	Account:      "test-user-001",
	Password:     "12345678",
	NickName:     "test-user1",
	Salt:         "1a2b3c4d5e",
	RegisterTime: time.Now().UnixMilli(),
}

var testUser2 = &model.User{
	ID:           1001001002,
	Account:      "test-user-002",
	Password:     "12345678",
	NickName:     "test-user2",
	Salt:         "1a2b3c4d5e",
	RegisterTime: time.Now().UnixMilli(),
}

func TestFindUserInfo(t *testing.T) {
	t.Run("existing-user", func(t *testing.T) {
		if userInfo, err := FindUserInfo(testUser1.ID); err != nil {
			t.Error(err)
			t.FailNow()
		} else {
			t.Logf("find user info success: %#v", userInfo)
		}
	})

	t.Run("not-existing-user", func(t *testing.T) {
		info, err := FindUserInfo(0)
		if err != nil {
			t.Error(err)
			t.FailNow()
		} else if info != nil {
			t.Logf("user not supposed to exists: %#v", info)
			t.FailNow()
		}
	})
}

func TestFindUserByAccount(t *testing.T) {
	_, ok, err := FindUserByAccount(testUser1.Account)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if !ok {
		t.Logf("user is supposed to exist")
		t.FailNow()
	}
}

func TestUpdateUserNickname(t *testing.T) {
	t.Run("do-update-nickname", func(t *testing.T) {
		if err := UpdateUserNickname(&model.UserInfo{ID: testUser1.ID, NickName: "changed"}); err != nil {
			t.Error(err)
			t.FailNow()
		}
	})

	t.Run("check-updated-user-info", func(t *testing.T) {
		if info, err := FindUserInfo(testUser1.ID); err != nil {
			t.Error(err)
			t.FailNow()
		} else if info == nil {
			t.Log("user is supposed to exist")
			t.FailNow()
		} else if info.NickName != "changed" {
			t.Log("change nickname failed")
			t.FailNow()
		}
	})
}
