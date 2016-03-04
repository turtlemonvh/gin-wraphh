# WrapHH [![GoDoc](https://godoc.org/github.com/turtlemonvh/gin-wraphh?status.svg)](https://godoc.org/github.com/turtlemonvh/gin-wraphh)

Use this to wrap middleware that accepts and returns `http.Handler` objects for use in gin.

Created to look like these helper methods in gin:

* https://godoc.org/github.com/gin-gonic/gin#WrapF
* https://godoc.org/github.com/gin-gonic/gin#WrapH

## Examples

See the test code.  I have examples wrapping [the NYT gzip library](https://github.com/NYTimes/gziphandler) and [NoSurf](https://github.com/justinas/nosurf), a popular CSRF protection middleware for golang.

## About

Based on this gist: https://gist.github.com/turtlemonvh/6cd23ef13e1e290717ef

## License

MIT
