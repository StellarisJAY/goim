package dao

import (
	"github.com/stellarisJAY/goim/pkg/db/model"
	"testing"
)

func init() {
	//_ = InsertFriendship(&model.Friend{FriendID: testUser1.ID, OwnerID: testUser2.ID}, &model.Friend{FriendID: testUser2.ID, OwnerID: testUser1.ID})
}

func TestCheckFriendship(t *testing.T) {
	t.Run("existing-friendship", func(t *testing.T) {
		if isFriend, err := CheckFriendship(testUser2.ID, testUser1.ID); err != nil {
			t.Error(err)
			t.FailNow()
		} else if !isFriend {
			t.Logf("user %d %d are supposed to be friends", testUser2.ID, testUser1.ID)
			t.FailNow()
		}
		if isFriend, err := CheckFriendship(testUser1.ID, testUser2.ID); err != nil {
			t.Error(err)
			t.FailNow()
		} else if !isFriend {
			t.Logf("user %d %d are supposed to be friends", testUser2.ID, testUser1.ID)
			t.FailNow()
		}
	})

	t.Run("not-friends", func(t *testing.T) {
		if isFriend, err := CheckFriendship(testUser2.ID, 111); err != nil {
			t.Error(err)
			t.FailNow()
		} else if isFriend {
			t.Logf("user %d %d are not supposed to be friends", testUser2.ID, 111)
			t.FailNow()
		}
	})
}

func TestListFriendIDs(t *testing.T) {
	if friends, err := ListFriendIDs(testUser1.ID); err != nil {
		t.Error(err)
		t.FailNow()
	} else if len(friends) != 1 {
		t.Logf("wrong number of friends")
		t.FailNow()
	} else if friends[0] != testUser2.ID {
		t.Logf("friends[0] is supposed to be %d", testUser2.ID)
		t.FailNow()
	}
}

func TestRemoveFriendship(t *testing.T) {
	f1 := &model.Friend{FriendID: 10001, OwnerID: 10002}
	f2 := &model.Friend{FriendID: 10002, OwnerID: 10001}
	if err := InsertFriendship(f1, f2); err != nil {
		t.Errorf("insert friendship failed %v", err)
		t.FailNow()
	}
	if cf, err := CheckFriendship(f1.FriendID, f1.OwnerID); err != nil {
		t.Errorf("insert friendship failed %v", err)
		t.FailNow()
	} else if !cf {
		t.Error("friendship inserted but can't verify")
	}
	if cf, err := CheckFriendship(f2.FriendID, f2.OwnerID); err != nil {
		t.Errorf("insert friendship failed %v", err)
		t.FailNow()
	} else if !cf {
		t.Error("friendship inserted but can't verify")
		t.FailNow()
	}

	if err := RemoveFriendship(f1.FriendID, f1.OwnerID); err != nil {
		t.Errorf("remove friendship error %v", err)
		t.FailNow()
	}

	isFriend, _ := CheckFriendship(f1.FriendID, f1.OwnerID)
	if isFriend {
		t.Errorf("remove friendship failed")
		t.FailNow()
	}
}
