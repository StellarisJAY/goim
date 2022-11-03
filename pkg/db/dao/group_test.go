package dao

import (
	"github.com/stellarisJAY/goim/pkg/db/model"
	"testing"
	"time"
)

var testGroup1 = &model.Group{
	ID:           1001001,
	Name:         "test-group-1",
	CreateTime:   time.Now().UnixMilli(),
	Description:  "this is a test group",
	OwnerID:      testUser1.ID,
	OwnerAccount: testUser1.Account,
}

func init() {
	//_ = InsertGroup(testGroup1)
	//_ = AddGroupMember(&model.GroupMember{GroupID: testGroup1.ID, UserID: 1001})
	//_ = AddGroupMember(&model.GroupMember{GroupID: testGroup1.ID, UserID: 1002})
	//_ = AddGroupMember(&model.GroupMember{GroupID: testGroup1.ID, UserID: 1003})
}

func TestAddGroupMember(t *testing.T) {

}

func TestFindGroupInfo(t *testing.T) {
	t.Run("existing-group-info", func(t *testing.T) {
		if info, err := FindGroupInfo(testGroup1.ID); err != nil {
			t.Error(err)
			t.FailNow()
		} else if info == nil {
			t.Logf("can't find existing group %d", testGroup1.ID)
			t.FailNow()
		} else if info.Name != testGroup1.Name {
			t.Logf("found wrong group info: %#v", info)
			t.FailNow()
		}
	})

	t.Run("non-existing-group-info", func(t *testing.T) {
		if info, err := FindGroupInfo(100); err != nil {
			t.Error(err)
			t.FailNow()
		} else if info != nil {
			t.Logf("group not supposed to exist")
			t.FailNow()
		}
	})

}

func TestGroupMemberExists(t *testing.T) {
	t.Run("existing-member", func(t *testing.T) {
		for i := 1001; i <= 1003; i++ {
			if isMember, err := GroupMemberExists(testGroup1.ID, int64(i)); err != nil {
				t.Error(err)
				t.FailNow()
			} else if !isMember {
				t.Logf("user %d supposed to exist", i)
				t.FailNow()
			}
		}
	})
	t.Run("non-existing-member", func(t *testing.T) {
		if isMember, err := GroupMemberExists(testGroup1.ID, int64(1004)); err != nil {
			t.Error(err)
			t.FailNow()
		} else if isMember {
			t.Logf("user %d not supposed to exist", 1004)
			t.FailNow()
		}
	})
}

func TestListGroupInfos(t *testing.T) {

}

func TestListUserJoinedGroupIds(t *testing.T) {
	t.Run("no-joined-groups", func(t *testing.T) {
		if groups, err := ListUserJoinedGroupIds(1004); err != nil {
			t.Error(err)
			t.FailNow()
		} else if groups != nil && len(groups) != 0 {
			t.Logf("user not supposed to have joined groups")
			t.FailNow()
		}
	})

	t.Run("has-joined-groups", func(t *testing.T) {
		for i := 1001; i <= 1003; i++ {
			if groups, err := ListUserJoinedGroupIds(int64(i)); err != nil {
				t.Error(err)
				t.FailNow()
			} else if groups == nil || len(groups) == 0 {
				t.Logf("user %d should have joined groups", i)
				t.FailNow()
			}
		}
	})
}

func TestListGroupMemberIDs(t *testing.T) {
	t.Run("has-member", func(t *testing.T) {
		if members, err := ListGroupMemberIDs(testGroup1.ID); err != nil {
			t.Error(err)
			t.FailNow()
		} else if members == nil || len(members) == 0 {
			t.Logf("group should have members")
			t.FailNow()
		}
	})

	t.Run("no-member", func(t *testing.T) {
		if members, err := ListGroupMemberIDs(101); err != nil {
			t.Error(err)
			t.FailNow()
		} else if members != nil && len(members) != 0 {
			t.Log("group should be empty")
			t.FailNow()
		}
	})
}

func TestPrepareData(t *testing.T) {
	for i := 1012; i <= 1999; i++ {
		//err := InsertUser(&model.User{
		//	ID:           int64(i),
		//	Account:      fmt.Sprintf("test-user-%d", i),
		//	Password:     "d232ca6a624b8a9cac6e152de5db10dc",
		//	NickName:     fmt.Sprintf("test-user-%d", i),
		//	Salt:         "xxvi7hvszwk1b182",
		//	RegisterTime: time.Now().UnixMilli(),
		//})
		err := AddGroupMember(&model.GroupMember{
			GroupID:  1001001,
			UserID:   int64(i),
			JoinTime: time.Now().UnixMilli(),
			Status:   0,
			Role:     6,
		})
		if err != nil {
			t.Error(err)
		}
	}
}
