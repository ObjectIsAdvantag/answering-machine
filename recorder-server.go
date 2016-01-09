// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

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

const version = "v0.2"



func main() {
	var showVersion, nameFromForm bool
	var port, name, upload, download, dir, fileID string
	flag.StringVar(&port, "port", "8081", "ip port of the server, defaults to 8081")
	flag.StringVar(&name, "name", "Recorder", "name of the service, defaults to Recorder")
	flag.StringVar(&upload, "upload", "upload", "route to store files, defaults to /upload")
	flag.StringVar(&download, "download", "download", "route to serve files, defaults to /downlaod")
	flag.StringVar(&dir, "directory", ".", "directory to store and serve files, defaults to .")
	flag.StringVar(&fileID, "formID", ".", "identifier of the file to updload, mandatory in current version")
	flag.BoolVar(&nameFromForm, "nameFromForm", false, "if set to true, names uploaded file with filename found in form, by default, the filename is extracted from URI")

	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.Parse()

	if showVersion {
		glog.Infof("%s version %s\n", name, version)
		return
	}

	if _, err := strconv.Atoi(port); err != nil {
		glog.Errorf("Invalid port: %s (%s)\n", port, err)
		return
	}

	service := &recorder.Recorder{
		UploadRoute: "/"+upload+"/",
		UploadDirectory: dir,
		FormIdentifier: fileID, // as specified in tropo documentation
		NameFileFromForm: nameFromForm,
		DownloadRoute : "/"+download+"/",
		DownloadDirectory: dir,
	}

	glog.V(1).Infof("Recorder configuration %s", service)

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
