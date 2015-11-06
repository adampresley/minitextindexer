package middleware

import (
	"net/http"

	"github.com/adampresley/minitextindexer/catalog"
	"github.com/adampresley/minitextindexer/config"

	"github.com/adampresley/logging"
	"github.com/gorilla/context"
)

/*
AppContext holds context data for the application. This can hold information
such as a database connection, session data, user info, and more. Your middlewares
should attach functions to this structure to pass critical data to request
handlers.
*/
type AppContext struct {
	Catalog *catalog.Catalog
	Config  *config.Configuration
	Log     *logging.Logger
	Version string
}

/*
StartAppContext is a middleware that should be early in the chain. This
sets up the initial context and attaches important data to the Gorilla
Context which comes across in the request.
*/
func (ctx *AppContext) StartAppContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		context.Set(request, "catalog", ctx.Catalog)
		context.Set(request, "config", ctx.Config)
		context.Set(request, "log", ctx.Log)
		context.Set(request, "version", ctx.Version)

		h.ServeHTTP(writer, request)
	})
}
