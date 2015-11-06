package main

import (
	"github.com/adampresley/minitextindexer/controllers"
	"github.com/adampresley/minitextindexer/listener"
	"github.com/adampresley/minitextindexer/middleware"
)

/*
Add routes here using AddRoute and AddRouteWithMiddleware.
*/
func setupRoutes(httpListener *listener.HTTPListenerService, appContext *middleware.AppContext) {
	httpListener.
		AddRoute("/getterm", controllers.GetSpecificTerm, "GET", "OPTIONS").
		AddRoute("/search", controllers.Search, "GET", "OPTIONS").
		AddRoute("/version", controllers.GetVersion, "GET")
}
