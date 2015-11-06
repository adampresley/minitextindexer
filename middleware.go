package main

import (
	"github.com/adampresley/minitextindexer/listener"
	"github.com/adampresley/minitextindexer/middleware"
)

func setupMiddleware(httpListener *listener.HTTPListenerService, appContext *middleware.AppContext) {
	httpListener.
		AddMiddleware(appContext.Logger).
		AddMiddleware(appContext.StartAppContext).
		AddMiddleware(appContext.AccessControl).
		AddMiddleware(appContext.OptionsHandler)
}
