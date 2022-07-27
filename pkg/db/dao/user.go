package dao

import (
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"gorm.io/gorm"
)

func FindUserByAccount(account string) (*model.User, bool, error) {
	user := new(model.User)
	tx := db.DB.MySQL.Where("account=?", account).First(&user)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, tx.Error
	}
	return user, true, nil
}

func InsertUser(user *model.User) error {
	tx := db.DB.MySQL.Create(user)
	return tx.Error
}

func InsertUserLoginLog(login *model.DeviceLogin) error {
	tx := db.DB.MySQL.Create(login)
	return tx.Error
}
