package controllers

import (
	"net/http"

	"github.com/adampresley/GoHttpService"
	"github.com/gorilla/context"
)

/*
GetVersion writes the current server version to the response writer

GET /version
*/
func GetVersion(writer http.ResponseWriter, request *http.Request) {
	version := (context.Get(request, "version")).(string)
	GoHttpService.WriteJson(writer, version, 200)
}
