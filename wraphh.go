package wraphh

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// A wrapper that turns a http.ResponseWriter into a gin.ResponseWriter, given an existing gin.ResponseWriter
// Needed if the middleware you are using modifies the writer it passes downstream
// FIXME: Wrap more methods: https://golang.org/pkg/net/http/#ResponseWriter
type wrappedResponseWriter struct {
	gin.ResponseWriter
	writer http.ResponseWriter
}

func (w *wrappedResponseWriter) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

func (w *wrappedResponseWriter) WriteString(s string) (n int, err error) {
	return w.writer.Write([]byte(s))
}

// An http.Handler that passes on calls to downstream middlewares
type nextRequestHandler struct {
	c      *gin.Context
	isNext bool //flag with next call
}

// Run the next request in the middleware chain and return
func (h *nextRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.c.Writer = &wrappedResponseWriter{h.c.Writer, w}
	h.c.Request = r
	h.c.Next()
	h.isNext = true
}

// Wrap something that accepts an http.Handler, returns an http.Handler
func WrapHH(hh func(h http.Handler) http.Handler) gin.HandlerFunc {
	// Steps:
	// - create an http handler to pass `hh`
	// - call `hh` with the http handler, which returns a function
	// - call the ServeHTTP method of the resulting function to run the rest of the middleware chain

	return func(c *gin.Context) {
		h := &nextRequestHandler{c, false}
		hh(h).ServeHTTP(c.Writer, c.Request)
		if h.isNext == false {
			c.Abort()
		}
	}
}
