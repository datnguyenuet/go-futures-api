package models

import (
	authService "go-futures-api/proto/auth"
	"strconv"
)

type User struct {
	Id                        string `json:"id" validate:"required"`
	Email                     string `json:"email" validate:"required"`
	AuthenticatorVerifyStatus int32  `json:"authenticator_verify_status" validate:"required"`
	AccountLv                 int32  `json:"account_lv" validate:"required"`
}

func UserFromProto(user *authService.User) (*User, error) {
	authenticatorVerifyStatus, _ := strconv.ParseInt(user.AuthenticatorVerifyStatus, 10, 32)
	accountLv, _ := strconv.ParseInt(user.AccountLv, 10, 32)
	return &User{
		Id:                        user.Id,
		Email:                     user.Email,
		AuthenticatorVerifyStatus: int32(authenticatorVerifyStatus),
		AccountLv:                 int32(accountLv),
	}, nil
}
