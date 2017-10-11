package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"prat/prat"
)

var name = flag.String("name", "anon", "your name")
var host = flag.String("host", "localhost", "the host to connect to")
var port = flag.Int("port", prat.DefaultPort, "the port to connect to")
var server = flag.Bool("server", false, "starts a server instead of a client")
var logFile = flag.String("log", "", "file to output logging information to")
var silent = flag.Bool("silent", false, "no log output")

func main() {

	flag.Parse()
	if *server {
		var logger *log.Logger
		if *silent {
			// if silent flag has been set, use a discard logger
			logger = log.New(ioutil.Discard, "", 0)
		} else {
			if *logFile != "" {
				// if a logfile has been specified, we create the file
				// and log to it
				file, err := os.Create(*logFile)
				if err != nil {
					panic(err)
				}
				logger = log.New(file, "", log.Ldate|log.Ltime)
			} else {
				// otherwise, output log to stdout
				logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
			}
		}
		server := prat.NewServer(logger)
		server.Start(*port)
	} else {
		address := fmt.Sprintf("%s:%d", *host, *port)
		client := prat.NewClient(address, *name)
		cui := prat.NewClientUI(client)
		if err := cui.Run(); err != nil {
			panic(err)
		}
	}
}
