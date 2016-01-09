// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

//
// Utility Server to upload Tropo Recordings and servce Audio Files
//
package recorder


import (
	"net/http"
	"os"
	"io"
	"strings"

	"github.com/golang/glog"
)


type Recorder struct {
	UploadRoute			string // URI to store files (POST)
	UploadDirectory		string // directory to which files are stored
	FormIdentifier		string // entry under which the file can be retrieved
	NameFileFromForm	bool // if set, the recorder stores the uploaded file with the name that appears in the form data file
	DownloadRoute		string // URI to retreive files (GET)
	DownloadDirectory   string // directory from which files are served
}


func (app *Recorder) RegisterHandlers() {
	http.HandleFunc(app.UploadRoute, app.uploadHandler)
	http.HandleFunc(app.DownloadRoute, app.downloadHandler)
}

func (app *Recorder) uploadHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("new upload")

	if req.Method != "POST" {
		glog.V(1).Infof("Expecting a POST, not a %s, exiting...", req.Method)
		http.Error(w, "Only POSTs are accepted here", http.StatusBadRequest)
		return
	}

	infile, header, err := req.FormFile(app.FormIdentifier)
	if err != nil {
		glog.V(1).Infof("Parsing error")
		http.Error(w, "Error parsing form: " + err.Error(), http.StatusBadRequest)
		return
	}
	defer infile.Close()

	// filename can be deduced from FormFile input or from Request Path
	var uploadFilename string
	if app.NameFileFromForm {
		uploadFilename = header.Filename
	} else {
		uploadFilename = strings.TrimPrefix(req.URL.Path, app.UploadRoute)
	}

	// Check filename
	if uploadFilename == "" {
		glog.V(1).Infof("Filename is absent, cannot store file", uploadFilename)
		http.Error(w, "No filename " + err.Error(), http.StatusBadRequest)
		return
	}

	// SECURITY sanatize the filename to prevent ../../
	if strings.Contains("/", uploadFilename) {
		glog.V(1).Infof("Security warning, filename contains slashes %s", uploadFilename)
		http.Error(w, "Slashes not supported for file upload " + err.Error(), http.StatusBadRequest)
		return
	}

	outfile, err := os.Create(app.UploadDirectory + "/" + uploadFilename)
	if err != nil {
		glog.V(1).Infof("Error creating file %s", uploadFilename)
		http.Error(w, "Error saving file " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer outfile.Close()

	_, err = io.Copy(outfile, infile)
	if err != nil {
		glog.V(1).Infof("Error saving file")

		http.Error(w, "Error saving file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	glog.V(0).Infof("File saved as: %s", uploadFilename)
	w.WriteHeader(http.StatusOK)
}


func (app *Recorder) downloadHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("new audio")

	if req.Method != "GET" {
		glog.V(1).Infof("Expecting a GET, not a %s, exiting...", req.Method)
		http.Error(w, "Only GET are accepted here", http.StatusBadRequest)
		return
	}

	glog.V(2).Infof("ressource is: %s", req.URL.Path)
	recording := strings.TrimPrefix(req.URL.Path, app.DownloadRoute)
	if recording == "" {
		glog.V(2).Infof("No file specified")
		http.Error(w, "No file specified", http.StatusBadRequest)
		return
	}

	glog.V(2).Infof("serving file %s", recording)

	// TODO: SECURITY sanatize the filename to prevent ../../
	infile, err := os.Open(app.DownloadDirectory + "/" + string(recording))
	defer infile.Close()
	if err != nil {
		glog.V(1).Infof("Error opening file %s", recording)
		http.Error(w, "Error opening file " + err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "audio/wav")
	_, err = io.Copy(w, infile)
	if err != nil {
		glog.V(1).Infof("Error sending file")
		http.Error(w, "Error sending file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	glog.V(0).Infof("Service file: %s", recording)
}