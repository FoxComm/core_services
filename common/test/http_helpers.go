package test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/FoxComm/FoxComm/configs"
)

type SimpleCodeHandler struct {
	Code int
}

func (h *SimpleCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(h.Code)
	w.Write([]byte(""))
}

var portCounter = 35107

func MakeHttpServer(handler http.Handler) *http.Server {
	port := configs.GetSafeRunPortStringFromString(fmt.Sprintf("%d", portCounter))
	portCounter += 1
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        handler,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go s.ListenAndServe()

	return s
}
