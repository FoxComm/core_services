package models

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/FoxComm/libs/logger"
	"github.com/FoxComm/libs/db/db_switcher"
	"github.com/FoxComm/libs/utils"
)

//Are we going to have one master Struct?  Or will we differentiate for inbound and on-site requests?
type SiteActivity struct {
	SharerToken        string //should be unique
	Timestamp          time.Time
	SharingActivityTag int
	RefererURL         string
	UserAgent          string
	LandingURL         string
	ApiRequestURL      string
	SessionToken       string
	RemoteIP           string
	Entity             Entity
	Action             string
	CheckoutDetails    CheckoutActivityDetails `,omitempty`
	SignupDetails      SignupActivityDetails   `,omitempty`
	Type               string
	StoreURL           string
	db_switcher.Mongo  `json:"-" bson:"-"`
}

func (act *SiteActivity) BeforeSave() {
	act.Timestamp = time.Now().Local()
}

func (sa *SiteActivity) AddActivityDetail(responseBody []byte) {
	switch sa.Action {
	case "signup", "signin":
		err := json.Unmarshal(responseBody, &sa.SignupDetails)
		if err != nil {
			logger.Error("[SiteActivity] Activity signup detail unmarshal error: %v", err)
		} else {
			sa.Entity = NewUserEntity(sa.SignupDetails.User)
		}
	case "checkout":
		err := json.Unmarshal(responseBody, &sa.CheckoutDetails)
		if err != nil {
			logger.Error("[SiteActivity] Activity checkout detail unmarshal error: %v", err)
		} else {
			sa.Entity = NewUserEntity(sa.CheckoutDetails.User)
		}
	}
}

func (sa *SiteActivity) AddActivityTypeDetail(r *http.Request, actionConfig SiteAction) {
	newReferrerUrlPattern := utils.ReferrerUrlPattern()
	if (actionConfig.QueryParams != nil && QueryParamsMatch(actionConfig, r)) || newReferrerUrlPattern.MatchString(r.RequestURI) {
		sa.Type = "inbound"
		sa.LandingURL = r.URL.Path
	} else {
		sa.Type = "onsite"
		sa.ApiRequestURL = r.URL.Path
	}
}
