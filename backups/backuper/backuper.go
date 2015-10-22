package backuper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/FoxComm/FoxComm/alerts"
	"github.com/FoxComm/FoxComm/configs"
	"github.com/FoxComm/FoxComm/logger"
	"github.com/FoxComm/core_services/feature_manager/core"
	. "github.com/FoxComm/libs/db/masterdb"
	"github.com/FoxComm/libs/utils"
	_ "github.com/FoxComm/libs/utils/ssl"
	"github.com/jrallison/go-workers"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/pivotal-golang/archiver/compressor"
	"github.com/stacktic/dropbox"
)

const DateFormat = "2006_01_02_15:04"

func Configure() {
	workers.Configure(map[string]string{
		// location of redis instance
		"server": configs.Get("RedisHost"),
		// instance of the database
		"database": "0",
		// number of connections to keep open with redis
		"pool": "30",
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process": "1",
	})
}

func Start(concurrency int) {
	workers.Process("backups", job, concurrency)
	workers.Run()
}

type Args struct {
	Job Job `json:"args"`
}

func job(message *workers.Msg) {
	jsonMessage := message.ToJson()
	args := Args{}

	err := json.Unmarshal([]byte(jsonMessage), &args)

	if err != nil {
		logger.Error("[backuper] job(): %s", err.Error())
		alerts.Slack(err.Error())
		return
	}

	if action, ok := Strategies[args.Job.Action]; ok {
		err := action(args.Job.Settings)

		if err != nil {
			logger.Error("[Backuper] %s", err.Error())
		}
	}
}

func BackupDatabaseFromConnection(connString string) (*bytes.Buffer, error) {
	cmd := exec.Command("pg_dump", "--format", "t", connString)
	out := bytes.NewBufferString("")
	cmd.Stdout = out
	err := cmd.Run()
	return out, err
}

func BackupMongoDatabaseFromConnection(storeFeature core.StoreFeature) (*bytes.Buffer, error) {
	feature := core.Feature{}
	DB.Model(&storeFeature).Related(&feature, "FeatureID")

	dstFolder := fmt.Sprintf("mongo_db_%s_backup_%s", feature.Name, time.Now().Format(DateFormat))
	dstFile := dstFolder + ".tar"

	datasource := strings.Split(storeFeature.Datasource, "#")
	dbHost := datasource[0]
	dbName := datasource[1]

	cmd := exec.Command("mongodump", "-d", dbName, "-h", dbHost, "-o", dstFolder)
	err := cmd.Run()

	if err != nil {
		return nil, err
	}

	err = CompressBackup(dstFolder)

	if err != nil {
		return nil, err
	}

	var out []byte

	out, err = ioutil.ReadFile(dstFile)

	if err != nil {
		return nil, err
	}

	os.RemoveAll(dstFolder)
	os.Remove(dstFile)

	return bytes.NewBuffer(out), nil
}

func bucket(settings Settings) *s3.Bucket {
	auth, err := aws.GetAuth(settings.S3Id, settings.S3Secret)

	if err != nil {
		alerts.Slack(err.Error())
		logger.Error("[backuper] bucket=", err.Error())
	}

	logger.Debug("S3 region: %+v", aws.Regions[settings.S3Region])

	httpClient := utils.GetHttpSslFlexibleClient()
	client := &s3.S3{
		Auth:   auth,
		Region: aws.Regions[settings.S3Region],
		HTTPClient: func() *http.Client {
			return httpClient
		},
	}
	return client.Bucket(settings.S3Bucket)
}

func dropboxClient(settings Settings) *dropbox.Dropbox {
	clientid := settings.DropboxId
	clientsecret := settings.DropboxSecret
	token := settings.DropboxToken

	db := dropbox.NewDropbox()
	db.SetAppInfo(clientid, clientsecret)
	db.SetAccessToken(token)
	return db
}

func downloadFileFromS3(bucket s3.Bucket, keys chan string, tempDir string, wg *sync.WaitGroup) {
	for key := range keys {
		newDir := fmt.Sprintf("%s/%s", tempDir, filepath.Dir(key))
		os.MkdirAll(newDir, 0755)
		reader, err := bucket.GetReader(key)

		if err != nil {
			logger.Error("[backuper] downloadFileFromS3(): %s", err.Error())
			alerts.Slack(err.Error())
		}

		logger.Debug("[backuper] Downloading: %s", key)
		data, err := ioutil.ReadAll(reader)

		if err != nil {
			logger.Error("[backuper] downloadFileFromS3(): %s", err.Error())
			alerts.Slack(err.Error())
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", tempDir, key), data, 0644)

		if err != nil {
			logger.Error("[backuper] downloadFileFromS3(): %s", err.Error())
			alerts.Slack(err.Error())
		}

		wg.Done()
		logger.Debug("Wait count: %+v", wg)
	}
}

func CompressBackup(tempDir string) error {
	out := bytes.NewBufferString("")
	err := compressor.WriteTar(tempDir, out)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s.tar", tempDir), out.Bytes(), 0644)
	return err
}
