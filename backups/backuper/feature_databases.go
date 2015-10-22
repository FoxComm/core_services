package backuper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/stacktic/dropbox"

	"github.com/FoxComm/FoxComm/alerts"
	"github.com/FoxComm/libs/logger"
	"github.com/FoxComm/core_services/feature_manager/core"
	. "github.com/FoxComm/libs/db/masterdb"
)

func FeatureDatabases(settings Settings) error {
	store := core.Store{}
	Db().First(&store, core.Store{Id: settings.StoreId})

	storeFeatures := []core.StoreFeature{}
	Db().Model(&store).Related(&storeFeatures)

	var out *bytes.Buffer
	var err error

	for _, storeFeature := range storeFeatures {
		if strings.Contains(storeFeature.Datasource, "dbname") {
			out, err = BackupDatabaseFromConnection(storeFeature.Datasource)
		} else {
			out, err = BackupMongoDatabaseFromConnection(storeFeature)
		}

		if err != nil {
			logger.Error("[backuper] FeatureDatabases(): %s", err.Error())
			alerts.Slack(err.Error())
			return err
		}

		db := dropboxClient(settings)

		feature := core.Feature{}
		Db().Model(&storeFeature).Related(&feature, "FeatureID")

		dstFile := fmt.Sprintf("db_%s_backup_%s.tar", feature.Name, time.Now().Format(DateFormat))
		_, err = db.UploadByChunk(ioutil.NopCloser(out), dropbox.DefaultChunkSize, dstFile, true, "")

		defer os.Remove(dstFile)

		if err != nil {
			logger.Error("[backuper] FeatureDatabases(): %s", err.Error())
			alerts.Slack(err.Error())
			return err
		}
	}

	return nil
}
