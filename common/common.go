package common

import (
	"net/http"
	"strconv"
)

const SiteActivityCookieName = "fc-site-activity-session"

func StoreID(r *http.Request) int {
	storeID, _ := strconv.Atoi(r.Header.Get("FC-Store-ID"))
	return storeID
}
