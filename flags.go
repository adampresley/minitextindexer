package main

import "flag"

var ip = flag.String("ip", "localhost", "IP address/hostname to bind this service to")
var port = flag.Int("port", 8999, "Port number to bind this service to")
var logLevel = flag.String("loglevel", "debug", "Set minimum log level. debug or info")
