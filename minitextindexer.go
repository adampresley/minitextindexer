package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/adampresley/minitextindexer/catalog"
	"github.com/adampresley/minitextindexer/config"
	"github.com/adampresley/minitextindexer/listener"
	"github.com/adampresley/minitextindexer/middleware"

	"github.com/adampresley/logging"
)

/* Version of the server */
const VERSION string = "v1.0.0"

func main() {
	flag.Parse()
	log := logging.NewLoggerWithMinimumLevel("Mini Text Indexer", logging.StringToLogType(*logLevel))
	var err error

	/*
	 * Load configuration data
	 */
	log.Info("Loading configuration file...")

	configuration, err := config.LoadConfigurationFromFile("./config.json")
	if err != nil {
		log.Fatalf("There was an error loading the configuration file config.json: %s", err.Error())
		os.Exit(1)
	}

	/*
	 * Setup shutdown channel, application context and HTTP listener. Start serving
	 */
	log.Info("Creating index...")

	catalog := catalog.NewCatalog(log, configuration)
	catalog.Index()
	//log.Info(catalog.ToJSON())
	//catalog.DisplayTree()

	appContext := &middleware.AppContext{
		Catalog: catalog,
		Config:  configuration,
		Log:     log,
		Version: VERSION,
	}

	log.Info("Starting HTTP server...")
	httpListener := listener.NewHTTPListenerService(*ip, *port, appContext)

	setupMiddleware(httpListener, appContext)
	setupRoutes(httpListener, appContext)

	go httpListener.StartHTTPListener()

	/*
	 * Block this thread until we receive SIGINT or
	 * SIGTERM
	 */
	doneChannel := make(chan os.Signal)
	signal.Notify(doneChannel, syscall.SIGINT, syscall.SIGTERM)
	log.Info(<-doneChannel)

	log.Info("Shut down.")
	os.Exit(0)
}
