package alerts

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"

	"github.com/FoxComm/FoxComm/configs"
	"github.com/FoxComm/libs/logger"
)

var slack_webhook_url string

func init() {
	env := configs.Get("FC_ENV")
	if env != "development" && env != "test" {
		slack_webhook_url = configs.Get("SLACK_WEBHOOK_URL")
	}
}

// Posts error messages to Slack via http.
// It receives a string representing an error message.
// It also specifies file name and line number of the caller.
func Slack(message string) {

	env := configs.Get("FC_ENV")
	if env == "development" || env == "test" {
		return
	}

	payload := fmt.Sprintf(`{
      "attachments":[{
        "color":  "danger",
        "fields": [{
          "title": "%s",
          "value": "File: %s"
        }]
      }]}`,
		message, file(),
	)

	res, err := http.Post(slack_webhook_url, "application/json", bytes.NewBuffer([]byte(
		payload,
	)))

	if err != nil {
		logger.Error("[alerts] Slack()", err.Error())
	}

	defer res.Body.Close()
}

// Returns second level caller's file name and line number joined by a colon.
func file() string {
	_, file, line, ok := runtime.Caller(2)

	if ok {
		return fmt.Sprintf("%s:%d", file, line)
	}

	return "Unknown file."
}
