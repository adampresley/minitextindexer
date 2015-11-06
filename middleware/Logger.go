package middleware

import (
	"net/http"
	"time"
)

/*
Logger is a middleware which logs requests to the logger. It also includes the
time it takes for the request to complete.
*/
func (ctx *AppContext) Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(writer, request)
		ctx.Log.Infof("%s - %s (%s)", request.Method, request.URL.String(), time.Since(startTime))
	})
}
