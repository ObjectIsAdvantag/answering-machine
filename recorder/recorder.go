// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

//
// Utility Server to store Tropo recordings
//
package recorder


import (
	"net/http"
	"os"
	"io"

	"github.com/golang/glog"
)


type Recorder struct {
	RecordingRoute				string // invoked after message are recorded
}

func (app *Recorder) RegisterHandlers() {
	http.HandleFunc(app.RecordingRoute, app.recordingHandler)
}

func (app *Recorder) recordingHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("Incoming recording")

	if req.Method != "POST" {
		glog.V(1).Infof("Expecting a POST, not a %s, exiting...", req.Method)
		http.Error(w, "Only POSTs are accepted here", http.StatusBadRequest)
		return
	}

	infile, header, err := req.FormFile("filename")
	if err != nil {
		glog.V(1).Infof("Parsing error")
		http.Error(w, "Error parsing uploads file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer infile.Close()

	// TODO: SECURITY sanatize the filename to prevent ../../
	recording := header.Filename

	outfile, err := os.Create("./uploads/" + recording)
	if err != nil {
		glog.V(1).Infof("Error creating file %s", recording)
		http.Error(w, "Error saving recording " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer outfile.Close()

	_, err = io.Copy(outfile, infile)
	if err != nil {
		glog.V(1).Infof("Error saving recording")

		http.Error(w, "Error saving recording: "+err.Error(), http.StatusInternalServerError)
		return
	}

	glog.V(0).Infof("Recording saved as: %s", recording)
	w.WriteHeader(http.StatusOK)
}
