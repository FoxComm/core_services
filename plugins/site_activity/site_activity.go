package site_activity

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/FoxComm/FoxComm/logger"
	"github.com/FoxComm/FoxComm/models"
	"github.com/FoxComm/FoxComm/repositories"
	"github.com/FoxComm/FoxComm/router/common"
	netutils "github.com/FoxComm/FoxComm/router/netutils"
	"github.com/FoxComm/FoxComm/utils"
	"github.com/jmcvetta/napping"
	oxyutils "github.com/mailgun/oxy/utils"
)

type SiteActivityMiddleware struct {
	next http.Handler
}

func New(next http.Handler) *SiteActivityMiddleware {
	return &SiteActivityMiddleware{next: next}
}

func (m *SiteActivityMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	accept := req.Header.Get("Accept")
	if strings.Contains(accept, "html") || strings.Contains(accept, "json") {
		logger.Debug("[SiteActivity] Looking for tracked action for request %s", req.RequestURI)

		body := &bytes.Buffer{}
		bw := oxyutils.NewBufferWriter(oxyutils.NopWriteCloser(body))

		m.next.ServeHTTP(bw, req)
		respBytes := body.Bytes()

		action := models.GetTrackedAction(req, respBytes, bw.Code)
		if action.Name != "" {
			logger.Debug("[SiteActivity] Action found: %s", action.Name)
		}

		cookies := netutils.FetchResponseCookies(bw.Header())

		logger.Debug("[SiteActivity] Response cookie array has length: %d", len(cookies))
		logger.Debug("[SiteActivity] Response cookie array: %+v", cookies)
		if len(cookies) == 0 {
			logger.Warn("[SiteActivity] Cookie list for HistoryWriter is empty for %s with URL %s", action.Name, req.RequestURI)
		}

		historyWriter := NewHistoryWriter(req, cookies, respBytes)
		historyWriter.CreateSiteActivity(action)

		// Write response to real responseWriter
		oxyutils.CopyHeaders(rw.Header(), bw.Header())
		rw.Header().Set("Content-Length", strconv.Itoa(body.Len()))
		rw.WriteHeader(bw.Code)
		io.Copy(rw, body)
	} else {
		m.next.ServeHTTP(rw, req)
	}
}

type HistoryWriter struct {
	Request       *http.Request
	Cookies       []*http.Cookie
	ResponseBody  []byte
	EntitySession *models.EntitySession
}

func NewHistoryWriter(req *http.Request, cookies []*http.Cookie, respBytes []byte) *HistoryWriter {
	return &HistoryWriter{
		Request:      req,
		Cookies:      cookies,
		ResponseBody: respBytes,
	}
}

func (hw *HistoryWriter) CreateSiteActivity(action models.SiteAction) {
	r := hw.Request
	storeID := common.StoreID(r)
	sArepo, err := repositories.NewSiteActivityRepoWithStoreId(storeID)
	if err != nil {
		logger.Error("[SiteActivity] Can't connect to db store: %v\nAction %+v skipped", err, action)
		return
	}

	sharerToken, tokenErr := sharerToken(r.RequestURI)
	logger.Debug("[SiteActivity] sharerToken is %s", sharerToken)
	if tokenErr != nil {
		logger.Warn("[SiteActivity] Error occured while extracting token %s", tokenErr.Error())
	}

	sharingActivityTag, tagErr := activityTag(r.RequestURI)
	logger.Debug("[SiteActivity] sharingActivityTag is %d", sharingActivityTag)
	if tagErr != nil {
		logger.Warn("[SiteActivity] Error occured while extracting tag %s", tagErr.Error())
	}

	if action.Name != "" {
		logger.Debug("[SiteActivity] Activity found: %s", action.Name)

		var cookie string

		for _, c := range hw.Cookies {
			logger.Debug("[SiteActivity] Processing cookie %s with value %s", c.Name, c.Value)
			if c.Name == common.SiteActivityCookieName {
				cookie = c.Value
				break
			}
		}

		if rc, err := r.Cookie(common.SiteActivityCookieName); cookie == "" && err == nil {
			logger.Debug("[SiteActivity] Using cookie from request: %v", rc.Value)
			cookie = rc.Value
		} else {
			logger.Error("[SiteActivity] HistoryWriter can't get session token from response for site activity: %v", err)
		}

		logger.Debug("[SiteActivity] Cookie value: %v", cookie)

		siteActivity := models.SiteActivity{
			Action:             action.Name,
			RefererURL:         r.Referer(),
			SessionToken:       cookie,
			UserAgent:          r.UserAgent(),
			RemoteIP:           r.RemoteAddr,
			SharerToken:        sharerToken,
			SharingActivityTag: sharingActivityTag,
			StoreURL:           hw.getStoreUrl(),
		}

		siteActivity.AddActivityDetail(hw.ResponseBody)
		siteActivity.AddActivityTypeDetail(r, action)

		if err := sArepo.Create(&siteActivity); err == nil {
			attributionsURL := fmt.Sprintf("%s/foxcomm/social_analytics/attributions", hw.getStoreUrl())

			logger.Debug("[SiteActivity] Sending process request to social analytics %s", attributionsURL)
			logger.Debug("[SiteActivity] Site activity %+v", siteActivity)

			socialAnalyticsResponse, err := sendRequest(attributionsURL, &siteActivity)

			if err != nil || socialAnalyticsResponse == nil || socialAnalyticsResponse.Status() != 202 {
				logger.Warn("[SiteActivity] Failed to create SA attribution")
				if err != nil {
					logger.Warn("[SiteActivity] Error happened while processing social analytics attribution: %s", err.Error())
				}
				go startRetryHandler(attributionsURL, &siteActivity)
			} else {
				logger.Info("[SiteActivity] Social analytics response is %s", socialAnalyticsResponse)
			}
		} else {
			logger.Error("[SiteActivity] An error occurred while saving the activity %s", err.Error())
		}
	}
}

func (hw *HistoryWriter) getStoreUrl() string {
	return hw.Request.Header.Get("FC-Store-Host")
}

func sharerToken(request string) (string, error) {
	newReferrerUrlPattern := utils.ReferrerUrlPattern()

	switch {
	case newReferrerUrlPattern.MatchString(request):
		logger.Debug("[SiteActivity] #getSharerToken New referrer url detected.")
		referralURI := newReferrerUrlPattern.FindString(request)
		return parameterFromRequestURL(referralURI, 2)
	default:
		logger.Debug("[SiteActivity] #getSharerToken Default url pattern detected.")
		return stringParameterFromRequest(request, "fc_ut")
	}
}

func activityTag(request string) (int, error) {
	newReferrerUrlPattern := utils.ReferrerUrlPattern()

	switch {
	case newReferrerUrlPattern.MatchString(request):
		logger.Debug("[SiteActivity] #getActivityTag New referrer url detected.")
		referralURI := newReferrerUrlPattern.FindString(request)
		return intParameterFromRequestURL(referralURI, 3)
	default:
		logger.Debug("[SiteActivity] #getActivityTag Default url pattern detected.")
		return intParameterFromRequest(request, "fc_at")
	}
}

func stringParameterFromRequest(request string, parameter string) (string, error) {
	urlWithQuery, err := url.ParseRequestURI(request)
	if err != nil {
		return "", err
	}
	queryValues := urlWithQuery.Query()
	return queryValues.Get(parameter), nil
}

func intParameterFromRequest(request string, parameter string) (int, error) {
	param, paramErr := stringParameterFromRequest(request, parameter)
	if paramErr != nil {
		return 0, paramErr
	}
	if param != "" {
		return strconv.Atoi(param)
	} else {
		return 0, nil
	}
}

func parameterFromRequestURL(requestURI string, index int) (string, error) {
	pathParts := strings.Split(requestURI, "/")
	if len(pathParts) <= index {
		return "", fmt.Errorf("Unable to grab the path from the request %s at index %d", requestURI, index)
	}
	return pathParts[index], nil
}

func intParameterFromRequestURL(request string, index int) (int, error) {
	param, paramErr := parameterFromRequestURL(request, index)
	if paramErr != nil {
		return 0, paramErr
	}
	if param != "" {
		return strconv.Atoi(param)
	} else {
		return 0, nil
	}
}

func sendRequest(attributionsURL string, siteActivity *models.SiteActivity) (*napping.Response, error) {
	logger.Debug("[SiteActivity] AttributionRetryMiddleware: Sending process request to social analytics %s", attributionsURL)
	session := &napping.Session{Client: utils.GetHttpSslFlexibleClient()}
	return session.Post(attributionsURL, siteActivity, nil, nil)
}

func retry(now time.Time, attributionsURL string, siteActivity *models.SiteActivity) (*napping.Response, error) {
	logger.Debug("[SiteActivity] AttributionRetryMiddleware: Running retry %s", now)
	return sendRequest(attributionsURL, siteActivity)
}

func startRetryHandler(attributionsURL string, siteActivity *models.SiteActivity) {
	c := time.Tick(1 * time.Minute)
	timeout := time.After(15 * time.Minute)
	for {
		select {
		case <-timeout:
			logger.Error("[SiteActivity] AttributionRetryMiddleware: Retries timed out for attribution: %+v", siteActivity)
			return
		case now := <-c:
			response, err := retry(now, attributionsURL, siteActivity)
			if err == nil && response != nil && response.Status() == 201 {
				logger.Info("[SiteActivity] AttributionRetryMiddleware: Successfully created attribution")
				break
			} else if err != nil {
				logger.Warn("[SiteActivity] AttributionRetryMiddleware: Error happened while processing social analytics attribution: %s", err.Error())
			}
		}
	}
	return
}
