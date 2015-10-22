package migration

import (
	"fmt"

	"github.com/FoxComm/core_services/user/service"
	"github.com/FoxComm/libs/db/db_switcher"
)

type Migration struct {
	db_switcher.PG
}

func (m Migration) Run() {
	if err := m.AutoMigrate(&user.User{}).Error; err != nil {
		fmt.Println("Failed to run user migration: " + err.Error())
	}
}
