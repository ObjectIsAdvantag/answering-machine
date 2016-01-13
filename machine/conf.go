// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

// Configuration API to read effective runtime configuration
package machine

import (
	"net/http"
	"encoding/json"

	"github.com/golang/glog"
	"github.com/ObjectIsAdvantag/answering-machine/tropo"
)


type I18nMessages struct {
	DefaultVoice				*tropo.Voice 		// see https://www.tropo.com/docs/webapi/international-features/speaking-multiple-languages
	WelcomeMessage				string     			// message played at incoming calls
	WelcomeAltMessage			string     			// message played at incoming calls when recording is not active
	CheckNoMessage				string
	CheckNewMessages			string
	RecordingOKMessage			string
	RecordingFailedMessage		string
}

type EnvConfiguration struct {
	RecorderEndpoint			string       		// URI to record the messages
	RecorderUsername			string
	RecorderPassword			string
	AudioServerEndpoint			string
	TranscriptsReceiver			string  			// email of the transcriptions receiver
	CheckerPhoneNumber			string    		 	// phone number to check messages
	CheckerFirstName			string       		// for greeting purpose
	DBfilename					string
	DBresetAtStartup			bool
}

type HandlerRoutes struct {
	IncomingCallRoute			string    			// route to the welcome message
	RecordingSuccessRoute		string       		// invoked after message are recorded
	RecordingIncompleteRoute	string   	 		// invoked if a timeout occurs
	RecordingFailedRoute		string    			// invoked if the recording failed due to communication issues between Tropo and the AnsweringMachine
	AdminRoute					string				// endpoint to browse voice messages
	ConfigurationRoute			string				// endpoint to read runtime effective configuration
}


func AddConfEndpoint(machine *AnsweringMachine, route string) {

	http.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
		glog.V(3).Infof("Conf API call: %s %s", req.Method, req.URL.String())

		if req.Method != "GET" {
			glog.V(2).Infof("Method %s not supported", req.Method)
			sendBadRequest(w, "only GET requests are supported")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("{"))

		// write env (mask password)
		maskedEnv := *machine.env
		maskedEnv.RecorderPassword = "********"
		w.Write([]byte(`"env":`))
		enc := json.NewEncoder(w)
		enc.Encode(maskedEnv)

		// write messages
		w.Write([]byte(`, "messages":`))
		enc.Encode(*machine.messages)

		// write routes
		w.Write([]byte(`, "routes":`))
		enc.Encode(*machine.routes)

		w.Write([]byte("}"))
	})
}


