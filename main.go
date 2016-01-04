package main

import (
	"os"
	"flag"
	"strconv"

	"github.com/golang/glog"

	"github.com/ObjectIsAdvantag/answering-machine/service"
)

const version = "0.1.draft"

func main() {
	var showVersion bool
	var port string
	flag.StringVar(&port, "port", "8080", "ip port of the service, defaults to 8080")
	flag.BoolVar(&showVersion, "version", false, "display version")

	flag.Parse()

	if showVersion {
		glog.Infof("SmartProxy version %s\n", version)
		return
	}

	if _, err := strconv.Atoi(port); err != nil {
		glog.Errorf("Invalid port: %s (%s)\n", port, err)
	}

	// [TODO] Initialize from an env variable
	var apiKey="REPLACE ME"

	glog.Infof("Starting Answering Machine, version: %s\n", version)

	if err := service.Run(apiKey, port, version); err != nil {
		glog.Errorf("Service exited with error: %s\n", err)
		glog.Flush()
		os.Exit(255)
		return
	}

	glog.Info("Service exited gracefully\n")
	glog.Flush()
}
