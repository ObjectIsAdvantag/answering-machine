// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

package tropo

import (
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"

	"github.com/golang/glog"
	"errors"
)


type CommunicationHandler struct {
	request			*http.Request
	writer 			http.ResponseWriter
	hasRead			bool // flag to memorize if incoming flow has been read
	hasWritten		bool // flag to memorize if outgoing flow has been written
}

// helper to compose an outgoing payloaqd
type Composer struct {
	commands		[]Command // ordered list of payloads
}

// tagging interface for commands, ie, encoding.json.Marshaler
type Command interface {
	MarshalJSON() ([]byte, error)
}


func NewHandler(w http.ResponseWriter, req *http.Request) *CommunicationHandler {
	return &CommunicationHandler{req, w, false, false}
}

func (d *CommunicationHandler) DecodeSession() (*Session, error) {

	if d.hasRead {
		glog.V(0).Infof("Implementation error : Bad usage, cannot read twice incoming payload\n")
		return nil, errors.New("Bad API usage, incoming flow already read")
	}

	d.hasRead = true

	if (d.request.Method != "POST") {
		glog.V(1).Infof("Unsupported incoming request: %s\n", d.request.Method)
		return nil, errors.New("Only POST is implemented")
	}

	decoder := json.NewDecoder(d.request.Body)
	var sw sessionWrapper
	if err := decoder.Decode(&sw); err != nil {
		glog.V(1).Infof("Could not parse json body")
		return nil, errors.New("Could not parse body")
	}

	return &(sw.Session), nil
}


func (d *CommunicationHandler) DecodeRecordingAnswer() (*RecordingResult, error) {

	if d.hasRead {
		glog.V(0).Infof("Implementation error : Bad usage, cannot read twice incoming payload\n")
		return nil, errors.New("Bad API usage, incoming flow already read")
	}

	d.hasRead = true

	if (d.request.Method != "POST") {
		glog.V(1).Infof("Unsupported incoming request: %s\n", d.request.Method)
		return nil, errors.New("Only POST is implemented")
	}

	decoder := json.NewDecoder(d.request.Body)
	var rw recordingResultWrapper
	if err := decoder.Decode(&rw); err != nil {
		glog.V(1).Infof("Could not parse json body")
		return nil, errors.New("Could not parse body")
	}

	return &(rw.RecordingResult), nil
}

// Starts composing an outgoing payload
func (d *CommunicationHandler) NewComposer() *Composer {
	commands := make([]Command,0,5)
	return &Composer{commands}

}

func (compo *Composer) AddCommand(cmd Command)  {
	// TODO check if at capacity
	compo.commands = append(compo.commands, cmd)
}

// Sends outgoing payload
func (handler *CommunicationHandler) ExecuteComposer(compo *Composer) error {
	if compo == nil {
		glog.V(0).Info("Implementation error : Bad API usage, composition should not be null\n")
		return errors.New("Bad API usage, no payload to transmit")
	}

	// Write outgoing payloads
	glog.V(2).Infof("Composing outgoing payload from %d commands", len(compo.commands))

	var b bytes.Buffer // A Buffer needs no initialization.
	b.Write([]byte(`{"tropo":[`))
	first := true
	for _, cmd := range compo.commands {
		if !first {
			b.WriteByte(',')
		}
		first = false
		encoded, err := cmd.MarshalJSON()
		if err != nil {
			glog.V(0).Infof("Could not encode Commands %v", cmd)
			return errors.New("Could not encode Command")
		}
		b.Write(encoded)
	}
	b.Write([]byte(`]}`))

	handler.SendRawJSON(b.String())

	return nil
}

func (handler *CommunicationHandler) ExecuteCommand(cmd Command) error {
	compo := handler.NewComposer()
	compo.AddCommand(cmd)
	return handler.ExecuteComposer(compo)
}


func (handler *CommunicationHandler) Say(message string, voice *Voice) error {
	// Assert message is not null, nor empty
	if len(message) == 0 {
		glog.V(0).Info("Implementation error : Bad API usage, cannot read twice incoming payload\n")
		return errors.New("Bad API usage, empty message")
	}

	cmd := &SayCommand{Message:message, Voice:voice}
	return handler.ExecuteCommand(cmd)
}

func (handler *CommunicationHandler) SendRawJSON(jsonString string) error {
	if handler.hasWritten {
		glog.V(0).Info("Implementation error : Bad API usage, cannot read twice incoming payload\n")
		return errors.New("Bad API usage, incoming flow already read")
	}
	handler.hasWritten = true

	handler.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(handler.writer, jsonString)
	return nil
}


func (handler *CommunicationHandler) ReplyInternalError() {

	glog.V(2).Infof("Starting\n")

	handler.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	handler.writer.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(handler.writer, `{ "error": { "status":"%d", "reason":"NOT_IMPLEMENTED", "message":"You hitted an endpoint that is not implemented yet, contact the author to speed up devs" } }`, http.StatusInternalServerError)
}


func (handler *CommunicationHandler) ReplyBadInput() {

	glog.V(2).Infof("Starting\n")

	handler.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	handler.writer.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(handler.writer, `{ "error": { "status":"%d", "reason":"NOT_IMPLEMENTED", "message":"You hitted an endpoint that is not implemented yet, contact the author to speed up devs" } }`, http.StatusInternalServerError)
}






