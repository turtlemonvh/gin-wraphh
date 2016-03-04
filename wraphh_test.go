package wrap_hh

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
)

const (
	testResponse = "cat cat cat cat cat cat cat cat "
)

func newServer() *gin.Engine {
	router := gin.Default()
	router.Use(WrapHH(gzipmiddle.GzipHandler))
	router.GET("/", func(c *gin.Context) {
		c.Header("Content-Length", strconv.Itoa(len(testResponse)))
		c.String(200, testResponse)
	})
	return router
}

// Based off:
// https://github.com/gin-gonic/contrib/blob/master/gzip/gzip_test.go
func TestGzip(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Encoding", "gzip")

	w := httptest.NewRecorder()
	r := newServer()
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

// NEXT: NoSurf
// https://github.com/justinas/nosurf/blob/master/handler.go#L93
