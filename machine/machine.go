// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

package machine


import (
	"fmt"
	"net/http"
	"errors"

	"github.com/golang/glog"
	"github.com/ObjectIsAdvantag/answering-machine/tropo"
)


type AnsweringMachine struct {
	env							*EnvConfiguration
	messages 					*I18nMessages
	routes 						*HandlerRoutes
	db 							*VoiceMessageStorage
}


func NewAnsweringMachine(env *EnvConfiguration, messages *I18nMessages) *AnsweringMachine {

	db, err := NewStorage(env.DBfilename, env.DBresetAtStartup)
	if err != nil {
		// TODO switch to ALT mode : say welcome message but do not record
		glog.Fatalf("Coud not create database to store messages states, error: %s", err)
		errors.New("Coud not create database to store messages states, exiting")
		return nil
	}

	routes := &HandlerRoutes{ "/tropo", "/tropo/recordingSuccess", "/tropo/recordingIncomplete", "/tropo/recordingFailed", "/messages", "/conf" }
	app := AnsweringMachine{env, messages, routes, db}

	glog.V(2).Infof("Created new AnsweringMachine with configuration %s", app)

	return &app
}


func (app *AnsweringMachine) RegisterHandlers() {
	http.HandleFunc(app.routes.IncomingCallRoute, app.incomingCallHandler)
	http.HandleFunc(app.routes.RecordingSuccessRoute, app.recordingSuccessHandler)
	http.HandleFunc(app.routes.RecordingIncompleteRoute, app.recordingIncompleteHandler)
	http.HandleFunc(app.routes.RecordingFailedRoute, app.recordingErrorHandler)

	// Add admin Endpoint
	if app.routes.AdminRoute != "" {
		AddAdminEndpoint(app.db, app.routes.AdminRoute)
		glog.V(0).Infof("Administration endpoint registered at %s", app.routes.AdminRoute)
	}

	// Add conf Endpoint
	if app.routes.ConfigurationRoute != "" {
		AddConfEndpoint(app, app.routes.ConfigurationRoute)
		glog.V(0).Infof("Configuration endpoint registered at %s", app.routes.ConfigurationRoute)
	}

}

func (app *AnsweringMachine) incomingCallHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("Incoming call")

	tropoHandler := tropo.NewHandler(w, req)
	if req.Method != "POST" {
		glog.V(0).Infof("POST expected, not a %s", req.Method)
		tropoHandler.ReplyBadRequest("Expecting a POST, please check tropo documentation")
		return
	}

	var session *tropo.Session
	var err error
    if session, err = tropoHandler.DecodeSession(); err != nil {
		glog.V(1).Infof("Cannot process incoming payload")
		tropoHandler.ReplyInternalError("DECODE FAILED", "Cannot process incoming payload")
		return
	}

	// check a human issued the call
	if !(session.IsHumanInitiated() && session.IsCall()) {
		glog.V(1).Infof("Unsupported request, a voice call is expected")
		tropoHandler.ReplyBadRequest("An incoming voice session is expected, not a M2M dialog")
		return
	}
	glog.V(0).Infof(`SessionID "%s", CallID "%s", From "+%s"`, session.ID, session.CallID, session.From.ID)

	// redirect to check messages if the answering machine registered checker is calling
	if session.From.ID == app.env.CheckerPhoneNumber {
		app.checkMessagesHandlerInternal(tropoHandler, session, w, req)
		return
	}

	app.welcomeHandlerInternal(tropoHandler, session, w, req)
}


func (app *AnsweringMachine) welcomeHandlerInternal(tropoHandler *tropo.CommunicationHandler, session *tropo.Session,w http.ResponseWriter, req *http.Request) {
	glog.V(3).Infof("welcomeHandlerInternal")

	// if no database to record message, say alternate welcome message
	if app.db == nil {
		tropoHandler.Say(app.messages.WelcomeAltMessage, app.messages.DefaultVoice)
		return
	}

	// store the new message entry
	voiceMessage := app.db.CreateVoiceMessage(session.CallID, "+" + session.From.ID)
	if err := app.db.Store(voiceMessage); err != nil {
		// say alternate welcome message
		tropoHandler.Say(app.messages.WelcomeAltMessage, app.messages.DefaultVoice)
		return
	}

	// please leave a message, start recording
	compo := tropoHandler.NewComposer()
	compo.AddCommand(&tropo.SayCommand{Message:app.messages.WelcomeMessage, Voice:app.messages.DefaultVoice})

	choices := tropo.RecordChoices{Terminator:"#"}
	transcript := tropo.RecordTranscription{ID:session.CallID, URL:app.env.TranscriptsReceiver}
    recorderURL := app.env.RecorderEndpoint + "/" + session.CallID + ".wav"
	compo.AddCommand(&tropo.RecordCommand{
		Bargein:true,
		Attempts:3,
		Beep:true,
		Choices:&choices,
		MaxSilence:3,
		Timeout:10,
		MaxTime:60,
		Name:"recording",
		URL:recorderURL,
		Username:app.env.RecorderUsername,
		Password:app.env.RecorderPassword,
		AsyncUpload:true,
		Transcription:&transcript})
	compo.AddCommand(&tropo.OnCommand{Event:"continue", Next:app.routes.RecordingSuccessRoute, Required:true})
	compo.AddCommand(&tropo.OnCommand{Event:"incomplete", Next:app.routes.RecordingIncompleteRoute, Required:true})
	compo.AddCommand(&tropo.OnCommand{Event:"error", Next:app.routes.RecordingFailedRoute, Required:true})

	tropoHandler.ExecuteComposer(compo)
}


func (app *AnsweringMachine) checkMessagesHandlerInternal(tropoHandler *tropo.CommunicationHandler, session *tropo.Session,w http.ResponseWriter, req *http.Request) {
	glog.V(3).Infof("checkMessagesHandlerInternal")

	// check if new messages
	messages := app.db.FetchNewMessages()
	nbOfNewMessages := len(messages)
	if nbOfNewMessages == 0 {
		msg := fmt.Sprintf(app.messages.CheckNoMessage, app.env.CheckerFirstName)
		tropoHandler.Say(msg, app.messages.DefaultVoice)
		return
	}

	compo := tropoHandler.NewComposer()
	msg := fmt.Sprintf(app.messages.CheckNewMessages, app.env.CheckerFirstName, nbOfNewMessages)
	compo.AddCommand(&tropo.SayCommand{Message:msg, Voice:app.messages.DefaultVoice})

	// TODO: Say when the message was recorded

	// play first message
	firstMessage := messages[0]
	audioFile := firstMessage.CallID + ".wav"
	audioURI:= app.env.AudioServerEndpoint + "/" + audioFile
	compo.AddCommand(&tropo.SayCommand{Message:audioURI})

	// TODO: play next messages
	//     - register event with messageID

	//     - Mark latest message as read (place code below in the commands callback)
	app.db.MarkMessageAsRead(firstMessage)

	tropoHandler.ExecuteComposer(compo)
}


func (app *AnsweringMachine) recordingSuccessHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("Recording response")

	tropoHandler := tropo.NewHandler(w, req)
	var answer *tropo.RecordingResult
	var err error
	if answer, err = tropoHandler.DecodeRecordingAnswer(); err != nil {
		glog.V(1).Infof("Cannot process recording result")
		tropoHandler.ReplyInternalError("DECODING ERROR", "Cannot process recording result")
		return
	}

	glog.V(0).Infof(`SessionID "%s", CallID "%s"`, answer.SessionID, answer.CallID)
	glog.V(2).Infof("Recording result details: %s\n", answer)

	// Store the recording success
	var vm *VoiceMessage
	vm, err = app.db.GetVoiceMessageForCallID(answer.CallID)
	if err != nil {
		glog.V(2).Infof("Cannot find message with callID: %s", answer.CallID)
		// TODO Analyse how often this case would happen, by default we'll fail but alternatively we could not create a brand new message

		tropoHandler.Say(app.messages.RecordingFailedMessage, app.messages.DefaultVoice)
		return
	}

	vm.Progress = RECORDED
	vm.Recording = answer.Actions.URL
	vm.Duration = answer.Actions.Duration
	vm.Status = NEW
	if err := app.db.Store(vm); err != nil {
		glog.V(2).Infof("Cannot update message with callID: %s", answer.CallID)
		// TODO Analyse how often this case would happen, we should at a minimum update the message state to FAILED

		// say alternate welcome message
		tropoHandler.Say(app.messages.RecordingFailedMessage, app.messages.DefaultVoice)
		return
	}

	// say good bye
	tropoHandler.Say(app.messages.RecordingOKMessage, app.messages.DefaultVoice)
}

func (app *AnsweringMachine) recordingIncompleteHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("RecordingIncompleteHandler")
	tropoHandler := tropo.NewHandler(w, req)
	tropoHandler.Say(app.messages.RecordingFailedMessage, app.messages.DefaultVoice)
}

func (app *AnsweringMachine) recordingErrorHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("RecordingErrordHandler")
	tropoHandler := tropo.NewHandler(w, req)
	tropoHandler.Say(app.messages.RecordingFailedMessage, app.messages.DefaultVoice)
}


