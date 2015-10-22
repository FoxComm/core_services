package migration

import (
	"fmt"

	"github.com/FoxComm/FoxComm/db/db_switcher"
	"github.com/FoxComm/FoxComm/user/service"
)

type Migration struct {
	db_switcher.PG
}

func (m Migration) Run() {
	if err := m.AutoMigrate(&user.User{}).Error; err != nil {
		fmt.Println("Failed to run user migration: " + err.Error())
	}
}
