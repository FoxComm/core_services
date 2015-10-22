package session

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/gob"
	"net/http"

	"github.com/FoxComm/FoxComm/db/db_switcher"
	"github.com/FoxComm/FoxComm/endpoints"
	"github.com/FoxComm/FoxComm/logger"
	"github.com/FoxComm/FoxComm/models"
	"github.com/FoxComm/router/common"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/kidstuff/mongostore"
)

func init() {
	gob.Register(Session{})
}

func NewSession() *Session {
	newToken := generateToken()
	return &Session{Token: newToken}
}

type Session struct {
	Token string
}

type SessionMiddleware struct {
	next http.Handler
}

func New(next http.Handler) *SessionMiddleware {
	return &SessionMiddleware{next: next}
}

func (s *SessionMiddleware) ServeHTTP(rw http.ResponseWriter, httpReq *http.Request) {
	s.next.ServeHTTP(rw, httpReq)
	action := models.GetSessionAction(httpReq, []byte{}, 0)

	// We don't care about non tracked actions
	if action.Name == "" {
		return
	}

	mongo := &db_switcher.Mongo{}
	err := mongo.InitializeWithStoreID(common.StoreID(httpReq), "sessions", endpoints.SocialAnalyticsAPI)
	if err != nil {
		logger.Error("[session] Can't connect to db_store: %s", err.Error())
		return
	}

	sessionCollection := mongo.Collection
	cookieStore := mongostore.NewMongoStore(sessionCollection, 86400*30, true, []byte("super-secret-donkey"))
	cookieStore.Options.Path = "/"

	cookie, err := cookieStore.Get(httpReq, common.SiteActivityCookieName)

	// If there is an error accessing the site activity cookie
	if err != nil || cookie.Values["session"] == nil {
		session := NewSession()
		cookie, _ = cookieStore.New(httpReq, common.SiteActivityCookieName)
		cookie.Values["session"] = session
		cookie.Save(httpReq, rw)
		return
	}

	// If action is signup, expire cookies
	if action.Name == "signout" {
		cookie.Options = &sessions.Options{MaxAge: -1, Path: "/"}
		cookie.Save(httpReq, rw)
		return
	}
}

func generateToken() string {
	randomKey := securecookie.GenerateRandomKey(20)
	hasher := sha1.New()
	hasher.Write(randomKey)
	token := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return token
}
