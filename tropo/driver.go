package tropo

import (
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/golang/glog"
	"errors"
)


type TropoDriver struct {
	request		*http.Request
	writer 		http.ResponseWriter
}


/*
var apiKey string
func init() {
	// [TODO] Initialize from an env variable
	apiKey="REPLACE ME"
}
*/


func NewDriver(w http.ResponseWriter, req *http.Request) *TropoDriver {
	return &TropoDriver{req, w}
}


func (d *TropoDriver) ReadSession() (*Session, error) {
	_ = "breakpoint"
	if (d.request.Method != "POST") {
		glog.V(1).Infof("Unsupported incoming request: %s\n", d.request.Method)
		return nil, errors.New("Only POST is implemented")
	}

	decoder := json.NewDecoder(d.request.Body)
	var sw SessionWrapper
	if err := decoder.Decode(&sw); err != nil {
		glog.V(1).Infof("Could not parse json body")
		return nil, errors.New("Could not parse body")
	}

	return &(sw.session), nil
}


func (d *TropoDriver) Say(message string, voice string) {
	d.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(d.writer, `{"tropo":[{"say":[{"value":"%s","voice":"%s"}]}]}`, message, voice)
}


func (d *TropoDriver) SendRaw(jsonString string) {
	d.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(d.writer, jsonString)
}


func (d *TropoDriver) ReplyInternalError() {
	glog.V(2).Infof("ReplyInternalError\n")

	d.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	d.writer.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(d.writer, `{ "error": { "status":"%d", "reason":"NOT_IMPLEMENTED", "message":"You hitted an endpoint that is not implemented yet, contact the author to speed up devs" } }`, http.StatusInternalServerError)
}


func (d *TropoDriver) ReplyBadInput() {
	glog.V(2).Infof("ReplyBadInput\n")

	d.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	d.writer.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(d.writer, `{ "error": { "status":"%d", "reason":"NOT_IMPLEMENTED", "message":"You hitted an endpoint that is not implemented yet, contact the author to speed up devs" } }`, http.StatusInternalServerError)
}
