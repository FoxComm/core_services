package backuper

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/FoxComm/core_services/lib/alerts"
	"github.com/FoxComm/libs/logger"
	"github.com/stacktic/dropbox"
)

func Assets(settings Settings) error {
	bucket := bucket(settings)
	content, err := bucket.GetBucketContents()

	if err != nil {
		logger.Error("[backups] Assets(): %s", err.Error())
		alerts.Slack(err.Error())
		return err
	}

	tempDir, err := ioutil.TempDir(".", "s3_backup_")

	if err != nil {
		logger.Error("[backups] Assets(): %s", err.Error())
		alerts.Slack(err.Error())
		return err
	}

	defer os.RemoveAll(tempDir)

	keys := make(chan string, 10)
	wg := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		go downloadFileFromS3(*bucket, keys, tempDir, wg)
	}

	for key, _ := range *content {
		wg.Add(1)
		keys <- key
	}

	wg.Wait()

	err = CompressBackup(tempDir)

	if err != nil {
		logger.Error("[backups] Assets(): %s", err.Error())
		alerts.Slack(err.Error())
		return err
	}

	db := dropboxClient(settings)

	dstFile := fmt.Sprintf("%s.tar", tempDir)
	defer os.Remove(dstFile)

	file, err := os.Open(dstFile)

	defer file.Close()

	if err != nil {
		logger.Error("[backups] Assets(): %s", err.Error())
		alerts.Slack(err.Error())
		return err
	}

	_, err = db.UploadByChunk(file, dropbox.DefaultChunkSize, dstFile, true, "")

	if err != nil {
		logger.Error("[backups] Assets(): %s", err.Error())
		alerts.Slack(err.Error())
	}

	return err
}
