package wraphh

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	// Example complex middlewares
	gzipmiddle "github.com/NYTimes/gziphandler"
	"github.com/justinas/nosurf"
)

const (
	testResponse = "cat cat cat cat cat cat cat cat "
)

var middlewareOptions = make(map[string]func(http.Handler) http.Handler)

func init() {
	middlewareOptions["gzip"] = gzipmiddle.GzipHandler
	middlewareOptions["nosurf"] = nosurf.NewPure
}

func newServer(mo string) *gin.Engine {
	router := gin.Default()

	if middlewareOptions[mo] != nil {
		router.Use(WrapHH(middlewareOptions[mo]))
	}

	router.GET("/", func(c *gin.Context) {
		c.Header("Content-Length", strconv.Itoa(len(testResponse)))
		c.String(200, testResponse)
	})

	router.POST("/", func(c *gin.Context) {
		c.Header("Content-Length", strconv.Itoa(len(testResponse)))
		c.String(200, testResponse)
	})

	return router
}

// Based off:
// https://github.com/gin-gonic/contrib/blob/master/gzip/gzip_test.go
func TestNYTGzip(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Encoding", "gzip")

	w := httptest.NewRecorder()
	r := newServer("gzip")
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Header().Get("Content-Encoding"), "gzip")
	assert.Equal(t, w.Header().Get("Vary"), "Accept-Encoding")
	assert.Equal(t, w.Header().Get("Content-Length"), "32")
	assert.NotEqual(t, w.Body.Len(), 32)
	assert.True(t, w.Body.Len() < 32, fmt.Sprintf("body length is %d, not <32", w.Body.Len()))

	gr, err := gzip.NewReader(w.Body)
	assert.NoError(t, err)
	defer gr.Close()

	body, _ := ioutil.ReadAll(gr)
	assert.Equal(t, string(body), testResponse)
}

// Should return a 400 because CSRF token is mising for POST request
func TestNoSurf(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", nil)

	w := httptest.NewRecorder()
	r := newServer("nosurf")
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, nosurf.FailureCode)
}

// Should return a 200
func TestNotNoSurf(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", nil)

	w := httptest.NewRecorder()
	r := newServer("")
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)
}
