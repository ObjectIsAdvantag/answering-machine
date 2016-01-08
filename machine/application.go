// Copyright 2015, Stève Sfartz
// Licensed under the MIT License
package machine


import (
	"net/http"

	"github.com/golang/glog"
	"github.com/ObjectIsAdvantag/answering-machine/tropo"
	"fmt"
)


type AnsweringMachine struct {
	welcomeMessage				string     // message played to callers
	defaultVoice				*tropo.Voice // see https://www.tropo.com/docs/webapi/international-features/speaking-multiple-languages
	welcomeMessageRoute			string    // route to the welcome message
	successRoute				string       // invoked after message are recorded
	incompleteRoute				string    // invoked if a timeout occurs
	communicationErrorRoute		string    // invoked if the recording failed due to communication issues between Tropo and the AnsweringMachine
	recorderEndpoint			string       // URI to record the messages
	transcriptsReceiver			string    // email of the transcriptions receiver
	checkerPhoneNumber			string     // phone number to check messages
	checkerFirstName			string       // for greeting purpose

}


func NewAnsweringMachine(welcomeMessage string, welcomeVoice *tropo.Voice, recordingEndpoint string, transcriptsEmail string, checkerPhoneNumber string, checkerFirstName string) *AnsweringMachine {
	if welcomeMessage == "" {
		welcomeMessage = "Laissez votre message après le bip sonore."
	}
	if welcomeVoice == nil {
		welcomeVoice = tropo.VOICE_AUDREY
	}

	app := AnsweringMachine{
		welcomeMessage,
		welcomeVoice,
		"/",
		"/answer",
		"/timeout",
		"/error",
		recordingEndpoint,
		"mailto:"+transcriptsEmail,
		checkerPhoneNumber,
		checkerFirstName,
	}

	glog.V(2).Infof("Created new AnsweringMachine with configuration %v", app)

	return &app
}


func (app *AnsweringMachine) RegisterHandlers() {
	http.HandleFunc(app.welcomeMessageRoute, app.welcomeHandler)
	http.HandleFunc(app.successRoute, app.recordingSuccessHandler)
	http.HandleFunc(app.incompleteRoute, app.recordingIncompleteHandler)
	http.HandleFunc(app.communicationErrorRoute, app.recordingErrorHandler)
}

func (app *AnsweringMachine) welcomeHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("Incoming call")

	tropoHandler := tropo.NewHandler(w, req)

	var session *tropo.Session
	var err error
    if session, err = tropoHandler.DecodeSession(); err != nil {
		glog.V(1).Infof("Cannot process incoming payload\n")
		tropoHandler.ReplyInternalError()
		return
	}

	// check a human issued the call
	if !(session.IsHumanInitiated() && session.IsCall()) {
		glog.V(1).Infof("Unsupported request, a voice call is expected\n")
		tropoHandler.ReplyBadInput()
		return
	}

	caller := session.From.ID
	glog.V(0).Infof(`SessionID "%s", CallID "%s", From "+%s"`, session.ID, session.CallID, caller)

	// redirect to check messages if the answering machine registered owner is calling
	if caller == app.checkerPhoneNumber {
		app.checkMessagesHandlerInternal(tropoHandler, session, w, req)
		return
	}

	// please leave a message, start recording
	compo := tropoHandler.NewComposer()
	compo.AddCommand(&tropo.SayCommand{Message:app.welcomeMessage, Voice:app.defaultVoice})
	choices := tropo.RecordChoices{Terminator:"#"}
	transcript := tropo.RecordTranscription{ID:session.CallID, URL:app.transcriptsReceiver}

	compo.AddCommand(&tropo.RecordCommand{Bargein:true, Attempts:3, Beep:true, Choices:&choices, MaxSilence:3, Timeout:10, MaxTime:60, Name:"recording", URL:"https://recorder.localtunnel.me/recordings", AsyncUpload:true, Transcription:&transcript})
	compo.AddCommand(&tropo.OnCommand{Event:"continue", Next:"/answer", Required:true})
	compo.AddCommand(&tropo.OnCommand{Event:"incomplete", Next:"/timeout", Required:true})
	compo.AddCommand(&tropo.OnCommand{Event:"error", Next:"/error", Required:true})

	tropoHandler.ExecuteComposer(compo)
}


func (app *AnsweringMachine) checkMessagesHandlerInternal(tropoHandler *tropo.CommunicationHandler, session *tropo.Session,w http.ResponseWriter, req *http.Request) {
	// check if new messages
	newMessage := 0
	if newMessage == 0 {
		msg := "Pas de nouveau message, bonne journée"
		if app.checkerFirstName != "" {
			msg = fmt.Sprintf("Bonjour %s, pas de nouveau message, bonne journée !", app.checkerFirstName)
		}

		tropoHandler.Say(msg, app.defaultVoice)
		return
	}

	// TODO: play message
}


func (app *AnsweringMachine) recordingSuccessHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("Recording response")

	tropoHandler := tropo.NewHandler(w, req)

	var answer *tropo.RecordingResult
	var err error
	if answer, err = tropoHandler.DecodeRecordingAnswer(); err != nil {
		glog.V(1).Infof("Cannot process recording result\n")
		tropoHandler.ReplyInternalError()
		return
	}

	glog.V(0).Infof(`SessionID "%s", CallID "%s"\n`, answer.SessionID, answer.CallID)
	glog.V(2).Infof("Recording result details: %s\n", answer)

	// say good bye
	tropoHandler.Say("Votre message est bien enregistré. Bonne journée !", tropo.VOICE_AUDREY)
}

func (app *AnsweringMachine) recordingIncompleteHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("RecordingIncompleteHandler")
}

func (app *AnsweringMachine) recordingErrorHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("RecordingErrordHandler")
}


