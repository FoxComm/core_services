package recover_middleware

import (
	"fmt"
	"net/http"

	"github.com/FoxComm/libs/logger"
	"github.com/FoxComm/vulcand/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/FoxComm/vulcand/plugin"
)

const Type = "recover"

func DefaultRecoverHandler(w http.ResponseWriter, r *http.Request) {
	body := "bad gateway"
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(body))
}

type RecoverMiddleware struct {
	next           http.Handler
	RecoverHandler http.HandlerFunc
}

func (m *RecoverMiddleware) NewHandler(next http.Handler) (http.Handler, error) {
	return &RecoverMiddleware{
		next:           next,
		RecoverHandler: DefaultRecoverHandler,
	}, nil
}

func (m *RecoverMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic in middleware: %+v", r)
			m.RecoverHandler(w, req)
		}
	}()
	m.next.ServeHTTP(w, req)
}

func GetSpec() *plugin.MiddlewareSpec {
	return &plugin.MiddlewareSpec{
		Type:      Type,
		FromOther: FromOther,
		FromCli:   FromCli,
		CliFlags:  []cli.Flag{},
	}
}

func FromOther(h RecoverMiddleware) (plugin.Middleware, error) {
	recoverHandler := h.RecoverHandler
	if recoverHandler == nil {
		recoverHandler = DefaultRecoverHandler
	}
	return &RecoverMiddleware{RecoverHandler: recoverHandler}, nil
}

func New() plugin.Middleware {
	m, _ := FromOther(RecoverMiddleware{})
	return m
}

func FromCli(c *cli.Context) (plugin.Middleware, error) {
	return FromOther(RecoverMiddleware{})
}
