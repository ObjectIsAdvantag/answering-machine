//
// Utility Server to store Tropo recordings
//
package main


import (
	"flag"
	"os"
	"strconv"

	"github.com/golang/glog"

	"github.com/ObjectIsAdvantag/answering-machine/server"
	"github.com/ObjectIsAdvantag/answering-machine/recorder"
)

const version = "v0.1"



func main() {
	var showVersion bool
	var port, name string
	flag.StringVar(&port, "port", "8081", "ip port of the server, defaults to 8081")
	flag.StringVar(&name, "name", "Recorder", "name of the service, defaults to Recorder")

	flag.BoolVar(&showVersion, "version", false, "display version")

	flag.Parse()

	if showVersion {
		glog.Infof("%s version %s\n", name, version)
		return
	}

	if _, err := strconv.Atoi(port); err != nil {
		glog.Errorf("Invalid port: %s (%s)\n", port, err)
	}

	service := &recorder.Recorder{"/recordings"}
	glog.Infof("Starting %s, version: %s\n", name, version)

	if err := server.Run(port, service, version, name); err != nil {
		glog.Errorf("Service exited with error: %s\n", err)
		glog.Flush()
		os.Exit(255)
		return
	}

	glog.Info("Service exited gracefully\n")
	glog.Flush()
}
