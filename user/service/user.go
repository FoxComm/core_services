package user

import (
	"crypto/md5"
	"fmt"

	"github.com/FoxComm/FoxComm/db/db_switcher"
)

const (
	SecretKey = "HiImASecretD0nk3y"
)

type User struct {
	db_switcher.PG    `sql:"-"`
	Id                int
	Name              string
	Email             string
	EncryptedPassword string
	Role              string
}

func toEncryptedPassword(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str+SecretKey)))
}

func (user *User) UpdatePassword(password string) bool {
	user.EncryptedPassword = toEncryptedPassword(password)
	return user.Save(user).Error == nil
}

func (user *User) IsValidPassword(password string) bool {
	return user.EncryptedPassword == toEncryptedPassword(password)
}
