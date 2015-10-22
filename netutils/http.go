package netutils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteTextResponse(w http.ResponseWriter, code int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	w.WriteHeader(code)
	w.Write([]byte(body))
}

func WriteJsonResponse(w http.ResponseWriter, code int, message interface{}) {
	bytes, err := json.Marshal(message)
	if err != nil {
		bytes = []byte("{}")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bytes)))
	w.WriteHeader(code)
	w.Write(bytes)
}

func FetchResponseCookies(header http.Header) []*http.Cookie {
	resp := http.Response{}
	resp.Header = header
	return resp.Cookies()
}
