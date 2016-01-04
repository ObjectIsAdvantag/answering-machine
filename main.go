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
		glog.Info("SmartProxy version %s", version)
		return
	}

	if _, err := strconv.Atoi(port); err != nil {
		glog.Error("Invalid port: %s (%s)\n", port, err)
	}

	// [TODO] Initialize from an env variable
	var apiKey="REPLACE ME"

	glog.Info("Starting Answering Machine, version: %s", version)

	if err := service.Run(apiKey, port, version); err != nil {
		glog.Error("Service exited with error: %s\n", err)
		glog.Flush()
		os.Exit(255)
		return
	}

	glog.Info("Service exited gracefully")
	glog.Flush()
}
