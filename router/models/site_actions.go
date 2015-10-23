package models

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
)

type ResponseChecks map[string]interface{}

type SiteAction struct {
	Name           string
	Path           string
	HTTPMethod     string
	QueryParams    []string
	StatusCode     int
	ResponseChecks ResponseChecks
}

//TODO: These actions are too specfic to BeautyKind, we should move this to a per store config
var SiteActions = []SiteAction{
	SiteAction{
		Name:       "signup",
		Path:       "/app/signup.js",
		HTTPMethod: "POST",
		StatusCode: 200,
	},
	SiteAction{
		Name:           "signup",
		Path:           "^/app/auth_sessions/facebook/([0-9a-zA-Z])+$",
		HTTPMethod:     "GET",
		StatusCode:     200,
		ResponseChecks: ResponseChecks{"sign_up": true},
	},
	SiteAction{
		Name:           "signin",
		Path:           "^/app/auth_sessions/facebook/([0-9a-zA-Z])+$",
		HTTPMethod:     "GET",
		StatusCode:     200,
		ResponseChecks: ResponseChecks{"sign_up": false},
	},
	SiteAction{
		Name:       "signin",
		Path:       "/app/login.js",
		HTTPMethod: "POST",
		StatusCode: 200,
	},
	SiteAction{
		Name:       "signout",
		Path:       "/app/logout.js",
		HTTPMethod: "GET",
		StatusCode: 204,
	},
	SiteAction{
		Name:           "checkout",
		Path:           "^/app/api/checkouts/([0-9a-zA-Z])+$",
		HTTPMethod:     "PUT",
		StatusCode:     200,
		ResponseChecks: ResponseChecks{"state": "complete"},
	},
	SiteAction{
		Name:        "referrer",
		QueryParams: []string{"fc_ut", "fc_at"},
	},
	SiteAction{
		Name: "referrer",
		Path: "/s/([0-9a-zA-Z_=-])+/([0-9])+$",
	},
}

var SessionActions = []SiteAction{
	SiteAction{
		Name: "signup",
		Path: "/app/signup.js",
	},
	SiteAction{
		Name: "signup",
		Path: "^/app/auth_sessions/facebook/([0-9a-zA-Z])+$",
	},
	SiteAction{
		Name: "signin",
		Path: "/app/login.js",
	},
	SiteAction{
		Name: "signout",
		Path: "/app/logout.js",
	},
	SiteAction{
		Name:        "referrer",
		QueryParams: []string{"fc_ut", "fc_at"},
	},
	SiteAction{
		Name: "referrer",
		Path: "/s/([0-9a-zA-Z_=-])+/([0-9])+$",
	},
}

func IsTrackedAction(r *http.Request, resp []byte, statusCode int) bool {
	action := GetTrackedAction(r, resp, statusCode)
	return action.Name != ""
}

func IsSessionAction(r *http.Request) bool {
	action := GetSessionAction(r, []byte{}, 0)
	return action.Name != ""
}

func GetSessionAction(r *http.Request, resp []byte, statusCode int) SiteAction {
	for _, action := range SessionActions {
		if actionMatch(action, r, resp, statusCode) {
			return action
		}
	}
	return SiteAction{}
}

func GetTrackedAction(r *http.Request, resp []byte, statusCode int) SiteAction {
	for _, action := range SiteActions {
		if actionMatch(action, r, resp, statusCode) {
			return action
		}
	}
	return SiteAction{}
}

func actionMatch(action SiteAction, r *http.Request, resp []byte, statusCode int) bool {
	return methodMatch(action, r) &&
		pathMatch(action, r) &&
		QueryParamsMatch(action, r) &&
		statusCodeMatch(action, statusCode) &&
		responseChecksMatch(action, resp)
}

func methodMatch(action SiteAction, r *http.Request) bool {
	if action.HTTPMethod == "" {
		return true
	} else {
		match := r.Method == action.HTTPMethod
		return match
	}
}

func pathMatch(action SiteAction, r *http.Request) bool {
	if action.Path == "" {
		return true
	} else {
		uri := mustParseURL(r.RequestURI)
		match := regexp.MustCompile(action.Path).MatchString(uri.Path)
		return match
	}
}

func QueryParamsMatch(action SiteAction, r *http.Request) bool {
	if action.QueryParams == nil {
		return true
	} else {
		uri := mustParseURL(r.RequestURI)
		requestParams := uri.Query()
		for _, param := range action.QueryParams {
			if requestParams.Get(param) != "" {
				return true
			}
		}
	}
	return false
}

func statusCodeMatch(action SiteAction, statusCode int) bool {
	if action.StatusCode == 0 {
		return true
	} else {
		match := action.StatusCode == statusCode
		return match
	}
}

func responseChecksMatch(action SiteAction, resp []byte) bool {
	if action.ResponseChecks == nil {
		return true
	} else {
		parsedResp := map[string]interface{}{}
		json.Unmarshal(resp, &parsedResp)
		for key, value := range action.ResponseChecks {
			if parsedResp[key] != value {
				return false
			}
		}
		return true
	}
}

func mustParseURL(uri string) *url.URL {
	parsedURL, _ := url.ParseRequestURI(uri)
	return parsedURL
}
