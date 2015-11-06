package middleware

import "net/http"

/*
Authorization is a middleware for handling authorization of requests. Here is
where you might authorize a token, user session, or cookie.
*/
func (ctx *AppContext) Authorization(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	})
}
