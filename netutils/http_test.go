package netutils

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"strconv"
	"fmt"
)

func TestFetchResponseCookies(t *testing.T) {
	headers := http.Header{}
	headers.Add("Set-Cookie","foo=bar")

	cookies := FetchResponseCookies(headers)

	assert.Equal(t, len(cookies), 1)
}

func TestWriteTextResponse(t *testing.T) {
	writer := httptest.NewRecorder()
	code := 200
	body := "Hello, FoxCommerce!"

	WriteTextResponse(writer, code, body)

	contentType := writer.Header().Get("Content-Type")
	assert.Equal(t, contentType, "text/plain; charset=utf-8")

	contentSize := writer.Header().Get("Content-Length")
	assert.Equal(t, contentSize, strconv.Itoa(len(body)))

	content := writer.Body.String()
	assert.Equal(t, body, content)
}

func TestWriteJsonResponse(t *testing.T) {
	writer := httptest.NewRecorder()
	code := 200
	body := `{'message': 'Hello, FoxCommerce!'}`

	WriteJsonResponse(writer, code, body)

	actualBody := fmt.Sprintf("\"%s\"", body) // it's possibly a hack

	contentType := writer.Header().Get("Content-Type")
	assert.Equal(t, "application/json; charset=utf-8", contentType)

	contentSize := writer.Header().Get("Content-Length")
	assert.Equal(t, strconv.Itoa(len(actualBody)), contentSize)

	content := writer.Body.String()
	assert.Equal(t, actualBody, content)
}



