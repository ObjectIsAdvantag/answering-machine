package tropo

import (
	"fmt"
	"time"

	"net/http"
	"encoding/json"

	"github.com/golang/glog"
)


type Application struct {

}


var apiKey string
func init() {
	// [TODO] Initialize from an env variable
	apiKey="REPLACE ME"
}

type SessionWrapper struct {
	Session struct {
				ID string `json:"id"`
				Accountid string `json:"accountId"`
				Timestamp time.Time `json:"timestamp"`
				Usertype string `json:"userType"`
				Initialtext interface{} `json:"initialText"`
				Callid string `json:"callId"`
				To struct {
					   ID string `json:"id"`
					   Name interface{} `json:"name"`
					   Channel string `json:"channel"`
					   Network string `json:"network"`
				   } `json:"to"`
				From struct {
					   ID string `json:"id"`
					   Name interface{} `json:"name"`
					   Channel string `json:"channel"`
					   Network string `json:"network"`
				   } `json:"from"`
				Headers struct {
					   MaxForwards string `json:"Max-Forwards"`
					   XSid string `json:"x-sid"`
					   RecordRoute string `json:"Record-Route"`
					   ContentLength string `json:"Content-Length"`
					   Contact string `json:"Contact"`
					   To string `json:"To"`
					   Cseq string `json:"CSeq"`
					   UserAgent string `json:"User-Agent"`
					   Via string `json:"Via"`
					   CallID string `json:"Call-ID"`
					   ContentType string `json:"Content-Type"`
					   From string `json:"From"`
				   } `json:"headers"`
			} `json:"session"`
}


func (app *Application) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	glog.V(2).Infof("Incoming call")

	if (req.Method != "POST") {
		glog.V(1).Infof("Unsupported incoming request: %s\n", req.Method)
		respondInternalError(w, req)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var sw SessionWrapper
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
	glog.V(0).Infof("Session %s, Call %s, From %s\n", sw.Session.ID, sw.Session.Callid, number)
	continueWorkflow(&sw, w, req)
}


func respondInternalError(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("respondInternalError\n")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{ "error": { "status":"%d", "reason":"NOT_IMPLEMENTED", "message":"You hitted an endpoint that is not implemented yet, contact the author to speed up devs" } }`, http.StatusInternalServerError)
}

func respondBadInput(sw *SessionWrapper, w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("respondBadInput\n")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{ "error": { "status":"%d", "reason":"NOT_IMPLEMENTED", "message":"You hitted an endpoint that is not implemented yet, contact the author to speed up devs" } }`, http.StatusInternalServerError)
}

func continueWorkflow(sw *SessionWrapper, w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("continueWorkflow\n")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{ "error": { "status":"%d", "reason":"NOT_IMPLEMENTED", "message":"You hitted an endpoint that is not implemented yet, contact the author to speed up devs" } }`, http.StatusInternalServerError)
}