// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

// Configuration API to read effective runtime configuration
package machine

import (
	"net/http"
	"encoding/json"
	"os"

	"github.com/golang/glog"
	"github.com/paked/configure"
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

func GetDefaultRoutes() (routes *HandlerRoutes) {
	return 	&HandlerRoutes{ "/tropo", "/tropo/recordingSuccess", "/tropo/recordingIncomplete", "/tropo/recordingFailed", "/messages", "/conf" }
}

//
func ReadConfiguration(envProperties string, messagesProperties string) (*EnvConfiguration, *I18nMessages) {

	conf := configure.New()
	glog.V(0).Infof("Loading configuration from 1. environnement")
	conf.Use(configure.NewEnvironment())
	if envProperties != "" {
		glog.V(0).Infof("Loading configuration from 2. json file: %s", envProperties)
		if file, err := os.Open(envProperties); err != nil {
			// TODO What if we wanted to pass all configuration via env variables, not possible as of today
			glog.Fatalf("Could not open Environment configuration: %s", envProperties)
		} else {
			file.Close()
			conf.Use(configure.NewJSONFromFile(envProperties))
		}
	}
	if messagesProperties != "" {
		glog.V(0).Infof("Loading messages from: %s", messagesProperties)
		if file, err := os.Open(messagesProperties); err != nil {
			glog.Warningf("Could not open Messages definition, switching to default values", messagesProperties)
		} else {
			file.Close()
			conf.Use(configure.NewJSONFromFile(messagesProperties))
		}
	}

	// environment specifics
	checkerPhoneNumber := conf.String("GOLAM_CHECKER_NUMBER", "", "the checker phone number to automate new messages check")
	checkerName := conf.String("GOLAM_CHECKER_NAME", "", "to enhance the welcome message of the new messages checker")
	recorderEndpoint := conf.String("GOLAM_RECORDER_ENDPOINT", "", "to receive the recordings")
	recorderUsername := conf.String("GOLAM_RECORDER_USERNAME", "", "credentials to the recorder endpoint")
	recorderPassword := conf.String("GOLAM_RECORDER_PASSWORD", "", "credentials to the recorder endpoint")
	audioEndpoint := conf.String("GOLAM_AUDIO_ENDPOINT", "", "audio files server")
	transcriptsEmail := conf.String("GOLAM_TRANSCRIPTS_EMAIL", "", "to receive transcripts via email")
	dbFilename := conf.String("GOLAM_DATABASE_PATH", "messages.db", "path to the messages database")
	dbResetAtStartup := conf.Bool("GOLAM_DATABASE_RESET", true, "flag to empty messages at startup")

	// messages
	defaultVoice := conf.String("GOLAM_VOICE", "Vanessa", "defaults to English")
	welcome := conf.String("GOLAM_WELCOME", "Welcome, please leave a message after the beep", "to enhance the welcome message of the new messages checker")
	welcomeAlt := conf.String("GOLAM_WELCOME_ALT", "Sorry we do not take any message currently, please call again later", "alternative message if storage service could not be started")
	checkNoMessage := conf.String("GOLAM_CHECK_NO_MESSAGE", "Hello %s, no new messages. Have a good day !", "")
	checkNewMessage := conf.String("GOLAM_CHECK_NEW_MESSAGES", "Hello %s, you have %d new messages", "")
	recordingOK := conf.String("GOLAM_RECORDING_OK", "Your message is recorded. Have a great day !", "")
	recordingFailed := conf.String("GOLAM_RECORDING_FAILED", "Sorry, we could not record your message. Please try again later", "")

	conf.Parse()


	var env EnvConfiguration
	env.RecorderEndpoint = *recorderEndpoint
	env.RecorderUsername= *recorderUsername
	env.RecorderPassword = *recorderPassword
	env.AudioServerEndpoint= *audioEndpoint
	env.TranscriptsReceiver= "mailto:" + *transcriptsEmail
	env.CheckerPhoneNumber= *checkerPhoneNumber
	env.CheckerFirstName = *checkerName
	env.DBfilename= *dbFilename
	env.DBresetAtStartup = *dbResetAtStartup

	var messages I18nMessages
	messages.DefaultVoice = tropo.GetVoice(*defaultVoice)
	messages.WelcomeMessage = *welcome
	messages.WelcomeAltMessage = *welcomeAlt
	messages.CheckNewMessages = *checkNewMessage
	messages.CheckNoMessage = *checkNoMessage
	messages.RecordingOKMessage = *recordingOK
	messages.RecordingFailedMessage = *recordingFailed

	return &env, &messages
}


// Adds an endpoint to display the active runtime configuration of the AnsweringMachine.
// By default, the configuration is accessible at /conf
// Note that the configuration is defined at startup and cannot be changed afterwards
func AddConfEndpoint(machine *AnsweringMachine, route string) {

	// default route
	if route == nil {
		route = "/conf"
	}

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
		publicEnv := *machine.env
		publicEnv.RecorderPassword = "********"
		w.Write([]byte(`"env":`))
		enc := json.NewEncoder(w)
		enc.Encode(publicEnv)

		// write messages
		w.Write([]byte(`, "messages":`))
		enc.Encode(*machine.messages)

		// write routes
		w.Write([]byte(`, "routes":`))
		enc.Encode(*machine.routes)

		w.Write([]byte("}"))
	})
}


