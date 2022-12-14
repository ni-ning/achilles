package auth

import (
	"achilles/global"
	"achilles/internal/model"
)

type AuthRegisterRequest struct {
	Username string `form:"username" binding:"required,min=1,max=128"`
	Password string `form:"username" binding:"required,min=1,max=128"`
	Role     uint8  `form:"role,default=1" binding:"oneof=0 1"`
}

type AuthLoginRequest struct {
	Username string `form:"username" binding:"required,min=1,max=128"`
	Password string `form:"username" binding:"required,min=1,max=128"`
}

func AccountRegister(req *AuthRegisterRequest) error {
	db := global.DBEngine.Model(&model.Account{})

	return db.Debug().Create(&model.Account{Username: req.Username,
		Password: req.Password, Role: req.Role}).Error
}

func AccountLogin(req *AuthLoginRequest) (bool, *model.Account) {
	var account = model.Account{}
	db := global.DBEngine.Model(&model.Account{})
	result := db.Where("username = ? and password = ?", req.Username, req.Password).First(&account)
	if result.RowsAffected > 0 {
		return true, &account
	} else {
		return false, nil
	}
}

func GetAccountById(accountId int64) (bool, *model.Account) {
	var account = model.Account{}
	db := global.DBEngine.Model(&model.Account{})
	result := db.Where("id = ?", accountId).First(&account)
	if result.RowsAffected > 0 {
		return true, &account
	} else {
		return false, nil
	}
}
