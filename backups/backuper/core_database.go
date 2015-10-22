package backuper

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/stacktic/dropbox"

	"github.com/FoxComm/FoxComm/alerts"
	"github.com/FoxComm/FoxComm/configs"
	"github.com/FoxComm/FoxComm/logger"
)

func CoreDatabase(settings Settings) error {
	out, err := BackupDatabaseFromConnection(configs.Get("FC_CORE_DB_URL"))

	if err != nil {
		logger.Error("[backuper] CoreDatabase(): %s", err.Error())
		alerts.Slack(err.Error())
		return err
	}

	db := dropboxClient(settings)

	dstFile := fmt.Sprintf("db_core_backup_%s.tar", time.Now().Format(DateFormat))
	defer os.Remove(dstFile)

	_, err = db.UploadByChunk(ioutil.NopCloser(out), dropbox.DefaultChunkSize, dstFile, true, "")

	if err != nil {
		logger.Error("[backuper] CoreDatabase(): %s", err.Error())
		alerts.Slack(err.Error())
	}

	return err
}
