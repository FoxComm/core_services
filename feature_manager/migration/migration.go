package migration

import (
	"github.com/FoxComm/FoxComm/db/masterdb"
	"github.com/FoxComm/core_services/backups/backuper"
	"github.com/FoxComm/core_services/feature_manager/core"
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
