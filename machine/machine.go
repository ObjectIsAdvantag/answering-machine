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

type I18nMessages struct {
	defaultVoice				*tropo.Voice 		// see https://www.tropo.com/docs/webapi/international-features/speaking-multiple-languages
	welcomeMessage				string     			// message played at incoming calls
	welcomeAltMessage			string     			// message played at incoming calls when recording is not active
	checkNoMessage				string
	checkNewMessages			string
	recordingOKMessage			string
	recordingFailedMessage		string
}

type EnvConfiguration struct {
	recorderEndpoint			string       		// URI to record the messages
	recorderUsername			string
	recorderPassword			string
	audioServerEndpoint			string
	transcriptsReceiver			string  			// email of the transcriptions receiver
	checkerPhoneNumber			string    		 	// phone number to check messages
	checkerFirstName			string       		// for greeting purpose
	dbFilename					string
}

type HandlerRoutes struct {

	welcomeMessageRoute			string    			// route to the welcome message
	recordingSuccessRoute		string       		// invoked after message are recorded
	recordingIncompleteRoute	string   	 		// invoked if a timeout occurs
	recordingFailedRoute		string    			// invoked if the recording failed due to communication issues between Tropo and the AnsweringMachine
	adminRoute					string				// webapi to browse voice messages
}

type AnsweringMachine struct {
	routes 						*HandlerRoutes
	messages 					*I18nMessages
	env							*EnvConfiguration
	db 							*VoiceMessageStorage
}


func NewAnsweringMachine(env *EnvConfiguration, routes *HandlerRoutes, messages *I18nMessages) *AnsweringMachine {

	db, err := NewStorage(env.dbFilename)
	if err != nil {
		// TODO switch to ALT mode : say welcome message but do not record
		glog.Fatalf("Coud not create database to store messages states, error: %s", err)
		errors.New("Coud not create database to store messages states, exiting")
		return nil
	}

	app := AnsweringMachine{routes, messages, env, db}

	glog.V(2).Infof("Created new AnsweringMachine with configuration %s", app)

	return &app
}


func (app *AnsweringMachine) RegisterHandlers() {
	http.HandleFunc(app.routes.welcomeMessageRoute, app.incomingCallHandler)
	http.HandleFunc(app.routes.recordingSuccessRoute, app.recordingSuccessHandler)
	http.HandleFunc(app.routes.recordingIncompleteRoute, app.recordingIncompleteHandler)
	http.HandleFunc(app.routes.recordingFailedRoute, app.recordingErrorHandler)

	// Add admin API
	if app.routes.adminRoute != "" {
		CreateAdminWebAPI(app.db, app.routes.adminRoute)
		glog.V(0).Infof("Admin API registered, browse message at URL http://.../admin ")
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
	if session.From.ID == app.env.checkerPhoneNumber {
		app.checkMessagesHandlerInternal(tropoHandler, session, w, req)
		return
	}

	app.welcomeHandlerInternal(tropoHandler, session, w, req)
}


func (app *AnsweringMachine) welcomeHandlerInternal(tropoHandler *tropo.CommunicationHandler, session *tropo.Session,w http.ResponseWriter, req *http.Request) {
	glog.V(3).Infof("welcomeHandlerInternal")

	// if no database to record message, say alternate welcome message
	if app.db == nil {
		tropoHandler.Say(app.messages.welcomeAltMessage, app.messages.defaultVoice)
		return
	}

	// store the new message entry
	voiceMessage := app.db.CreateVoiceMessage(session.CallID, "+" + session.From.ID)
	if err := app.db.Store(voiceMessage); err != nil {
		// say alternate welcome message
		tropoHandler.Say(app.messages.welcomeAltMessage, app.messages.defaultVoice)
		return
	}

	// please leave a message, start recording
	compo := tropoHandler.NewComposer()
	compo.AddCommand(&tropo.SayCommand{Message:app.messages.welcomeMessage, Voice:app.messages.defaultVoice})

	choices := tropo.RecordChoices{Terminator:"#"}
	transcript := tropo.RecordTranscription{ID:session.CallID, URL:app.env.transcriptsReceiver}
    recorderURL := app.env.recorderEndpoint + "/" + session.CallID + ".wav"
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
		Username:app.env.recorderUsername,
		Password:app.env.recorderPassword,
		AsyncUpload:true,
		Transcription:&transcript})
	compo.AddCommand(&tropo.OnCommand{Event:"continue", Next:"/answer", Required:true})
	compo.AddCommand(&tropo.OnCommand{Event:"incomplete", Next:"/timeout", Required:true})
	compo.AddCommand(&tropo.OnCommand{Event:"error", Next:"/error", Required:true})

	tropoHandler.ExecuteComposer(compo)
}


func (app *AnsweringMachine) checkMessagesHandlerInternal(tropoHandler *tropo.CommunicationHandler, session *tropo.Session,w http.ResponseWriter, req *http.Request) {
	glog.V(3).Infof("checkMessagesHandlerInternal")

	// check if new messages
	messages := app.db.FetchNewMessages()
	nbOfNewMessages := len(messages)
	if nbOfNewMessages == 0 {
		msg := fmt.Sprintf(app.messages.checkNoMessage, app.env.checkerFirstName)
		tropoHandler.Say(msg, app.messages.defaultVoice)
		return
	}

	compo := tropoHandler.NewComposer()
	msg := fmt.Sprintf(app.messages.checkNewMessages, app.env.checkerFirstName, nbOfNewMessages)
	compo.AddCommand(&tropo.SayCommand{Message:msg, Voice:app.messages.defaultVoice})

	// TODO: Say when the message was recorded

	// play first message
	firstMessage := messages[0]
	audioFile := firstMessage.CallID + ".wav"
	audioURI:= app.env.audioServerEndpoint + "/" + audioFile
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

		tropoHandler.Say("Désolé, nous n'avons pas pu enregistrer votre message. Merci de ré essayer !", app.messages.defaultVoice)
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
		tropoHandler.Say(app.messages.recordingFailedMessage, app.messages.defaultVoice)
		return
	}

	// say good bye
	tropoHandler.Say(app.messages.recordingOKMessage, app.messages.defaultVoice)
}

func (app *AnsweringMachine) recordingIncompleteHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("RecordingIncompleteHandler")
	tropoHandler := tropo.NewHandler(w, req)
	tropoHandler.Say(app.messages.recordingFailedMessage, app.messages.defaultVoice)
}

func (app *AnsweringMachine) recordingErrorHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("RecordingErrordHandler")
	tropoHandler := tropo.NewHandler(w, req)
	tropoHandler.Say(app.messages.recordingFailedMessage, app.messages.defaultVoice)
}


