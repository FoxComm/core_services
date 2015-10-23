package backuper

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/stacktic/dropbox"

	"github.com/FoxComm/core_services/lib/alerts"
	"github.com/FoxComm/libs/logger"
)

func OriginDatabase(settings Settings) error {
	out, err := BackupDatabaseFromConnection(settings.DatabaseUrl)

	if err != nil {
		logger.Error("[backuper] OriginDatabase(): %s", err.Error())
		alerts.Slack(err.Error())
		return err
	}

	db := dropboxClient(settings)

	dstFile := fmt.Sprintf("db_origin_backup_%s.tar", time.Now().Format(DateFormat))
	defer os.Remove(dstFile)

	_, err = db.UploadByChunk(ioutil.NopCloser(out), dropbox.DefaultChunkSize, dstFile, true, "")

	if err != nil {
		logger.Error("[backuper] OriginDatabase(): %s", err.Error())
		alerts.Slack(err.Error())
	}

	return err
}
