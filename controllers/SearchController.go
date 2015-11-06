package controllers

import (
	"net/http"

	"github.com/adampresley/GoHttpService"
	"github.com/adampresley/logging"
	"github.com/adampresley/minitextindexer/catalog"
	"github.com/gorilla/context"
)

/*
GetSpecificTerm tries to find nodes that match a specific term

GET /getterm?term=[searchTerm]
*/
func GetSpecificTerm(writer http.ResponseWriter, request *http.Request) {
	log := (context.Get(request, "log")).(*logging.Logger)
	catalog := (context.Get(request, "catalog")).(*catalog.Catalog)
	term := request.URL.Query().Get("term")

	if len(term) <= 0 {
		log.Error("User provided blank term in /getterm")
		GoHttpService.BadRequest(writer, "Please provide a search term")
		return
	}

	log.Infof("Getting term for [%s]", term)

	matchedTerm := catalog.FindTerm(term)
	if matchedTerm == nil {
		GoHttpService.NotFound(writer, "Term "+term+" not found")
		return
	}

	GoHttpService.WriteJson(writer, matchedTerm, 200)
}

/*
Search tries to find nodes that contain a term

GET /search?term=[searchTerm]
*/
func Search(writer http.ResponseWriter, request *http.Request) {
	log := (context.Get(request, "log")).(*logging.Logger)
	catalog := (context.Get(request, "catalog")).(*catalog.Catalog)
	term := request.URL.Query().Get("term")

	if len(term) <= 0 {
		log.Error("User provided blank term in /search")
		GoHttpService.BadRequest(writer, "Please provide a search term")
		return
	}

	log.Infof("Searching for [%s]", term)

	matches := catalog.Search(term)
	if matches == nil {
		GoHttpService.NotFound(writer, "Term "+term+" not found")
		return
	}

	GoHttpService.WriteJson(writer, matches, 200)
}
