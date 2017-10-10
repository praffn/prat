package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"prat/prat"
)

var name = flag.String("name", "anon", "your name")
var host = flag.String("host", "localhost", "the host to connect to")
var port = flag.Int("port", prat.DefaultPort, "the port to connect to")
var server = flag.Bool("server", false, "starts a server instead of a client")
var logFile = flag.String("log", "prat.log", "file to output logging information to")

func main() {
	flag.Parse()
	if *server {
		file, err := os.Create(*logFile)
		if err != nil {
			panic(err)
		}
		logger := log.New(file, "", 0)
		server := prat.NewServerWithLogger(logger)
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
