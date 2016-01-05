package tropo

import (
	"fmt"

	"net/http"
	"encoding/json"

	"github.com/golang/glog"
)


type Application struct {
	DefaultVoice	string // see https://www.tropo.com/docs/webapi/international-features/speaking-multiple-languages
}


var apiKey string
func init() {
	// [TODO] Initialize from an env variable
	apiKey="REPLACE ME"
}

func NewApplication() *Application {
	app := Application{"Audrey"}
	return &app
}


func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	glog.V(2).Infof("Incoming call")

	if (req.Method != "POST") {
		glog.V(1).Infof("Unsupported incoming request: %s\n", req.Method)
		respondInternalError(w, req)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var sw sessionWrapper
	if err := decoder.Decode(&sw); err != nil {
		glog.V(1).Infof("Could not parse json body")
		respondInternalError(w, req)
		return
	}

	// check a human is calling
	if sw.Session.Usertype != "HUMAN" || sw.Session.From.Channel != "VOICE" {
		glog.V(1).Infof("Unsupported incoming request: %s\n", req.Method)
		respondBadInput(&sw, w, req)
		return
	}

	// echo leave a message
	number := sw.Session.From.ID
	glog.V(0).Infof(`SessionID "%s", CallID "%s", From "+%s"`, sw.Session.ID, sw.Session.Callid, number)
	continueWorkflow(&sw, w, req)
}


func respondInternalError(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("respondInternalError\n")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{ "error": { "status":"%d", "reason":"NOT_IMPLEMENTED", "message":"You hitted an endpoint that is not implemented yet, contact the author to speed up devs" } }`, http.StatusInternalServerError)
}

func respondBadInput(sw *sessionWrapper, w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("respondBadInput\n")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{ "error": { "status":"%d", "reason":"NOT_IMPLEMENTED", "message":"You hitted an endpoint that is not implemented yet, contact the author to speed up devs" } }`, http.StatusInternalServerError)
}

func continueWorkflow(sw *sessionWrapper, w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof(`continueWorkflow for callID "%s"`, sw.Session.Callid)

	// Say voice message
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, `{"tropo":[{"say":[{"value":"Bienvenue chez Stève, Valérie, Jeanne et Olivia. Bonne année 2016 ! Laissez votre message.","voice":"Audrey"}]}]}`)
}