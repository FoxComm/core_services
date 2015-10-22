package migration

import (
	"github.com/FoxComm/core_services/backups/backuper"
	"github.com/FoxComm/core_services/feature_manager/core"
	"github.com/FoxComm/libs/db/masterdb"
)

func Run() {
	masterdb.Db().AutoMigrate(
		&core.Merchant{},
		&core.Store{},
		&core.Domain{},
		&core.Feature{},
		&core.StoreFeature{},
		&backuper.Settings{},
	)
}
