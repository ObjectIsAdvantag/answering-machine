// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

package machine


import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/ObjectIsAdvantag/answering-machine/tropo"
)


type AnsweringMachine struct {
	welcomeMessage				string     			// message played at incoming calls
	welcomeAltMessage			string     			// message played at incoming calls when recording is not active
	defaultVoice				*tropo.Voice 		// see https://www.tropo.com/docs/webapi/international-features/speaking-multiple-languages
	welcomeMessageRoute			string    			// route to the welcome message
	successRoute				string       		// invoked after message are recorded
	incompleteRoute				string   	 		// invoked if a timeout occurs
	communicationErrorRoute		string    			// invoked if the recording failed due to communication issues between Tropo and the AnsweringMachine
	recorderEndpoint			string       		// URI to record the messages
	recorderUsername			string
	recorderPassword			string
	audioServerEndpoint			string
	transcriptsReceiver			string  			// email of the transcriptions receiver
	checkerPhoneNumber			string    		 	// phone number to check messages
	checkerFirstName			string       		// for greeting purpose

	db 							*VoiceMessageStorage
	adminRoute					string
}


func NewAnsweringMachine(welcomeMessage string, welcomeAltMessage string, welcomeVoice *tropo.Voice, recordingEndpoint string, recordingUsername string, recordingPassword string, audioEndpoint string, transcriptsEmail string, checkerPhoneNumber string, checkerFirstName string) *AnsweringMachine {
	if welcomeMessage == "" {
		welcomeMessage = "Welcome, please leave a message after the beep"
	}
	if welcomeAltMessage == "" {
		welcomeAltMessage = "Sorry we do not take any message currently, please call again later"
	}
	if welcomeVoice == nil {
		welcomeVoice = tropo.VOICE_AUDREY
	}

	app := AnsweringMachine{
		welcomeMessage,
		welcomeAltMessage,
		welcomeVoice,
		"/",
		"/answer",
		"/timeout",
		"/error",
		recordingEndpoint,
		recordingUsername,
		recordingPassword,
		audioEndpoint,
		"mailto:"+transcriptsEmail,
		checkerPhoneNumber,
		checkerFirstName,
		nil,
		"/admin",
	}

	db, err := NewStorage("messages.db")
	if err != nil {
		// TODO switch to ALT mode : say welcome message but do not record
		glog.V(0).Infof("Coud not create database to store messages states, exiting", app)

	}
	app.db = db

	glog.V(2).Infof("Created new AnsweringMachine with configuration %s", app)

	return &app
}


func (app *AnsweringMachine) RegisterHandlers() {
	http.HandleFunc(app.welcomeMessageRoute, app.incomingCallHandler)
	http.HandleFunc(app.successRoute, app.recordingSuccessHandler)
	http.HandleFunc(app.incompleteRoute, app.recordingIncompleteHandler)
	http.HandleFunc(app.communicationErrorRoute, app.recordingErrorHandler)

	// Add admin API
	if app.adminRoute != "" {
		CreateAdminWebAPI(app.db, app.adminRoute)
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
	if session.From.ID == app.checkerPhoneNumber {
		app.checkMessagesHandlerInternal(tropoHandler, session, w, req)
		return
	}

	app.welcomeHandlerInternal(tropoHandler, session, w, req)
}


func (app *AnsweringMachine) welcomeHandlerInternal(tropoHandler *tropo.CommunicationHandler, session *tropo.Session,w http.ResponseWriter, req *http.Request) {
	glog.V(3).Infof("welcomeHandlerInternal")

	// if no database to record message, say alternate welcome message
	if app.db == nil {
		tropoHandler.Say(app.welcomeAltMessage, app.defaultVoice)
		return
	}

	// store the new message entry
	voiceMessage := app.db.CreateVoiceMessage(session.CallID, "+" + session.From.ID)
	if err := app.db.Store(voiceMessage); err != nil {
		// say alternate welcome message
		tropoHandler.Say(app.welcomeAltMessage, app.defaultVoice)
		return
	}

	// please leave a message, start recording
	compo := tropoHandler.NewComposer()
	compo.AddCommand(&tropo.SayCommand{Message:app.welcomeMessage, Voice:app.defaultVoice})

	choices := tropo.RecordChoices{Terminator:"#"}
	transcript := tropo.RecordTranscription{ID:session.CallID, URL:app.transcriptsReceiver}
    recorderURL := app.recorderEndpoint + "/" + session.CallID + ".wav"
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
		Username:app.recorderUsername,
		Password:app.recorderPassword,
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
	nbOfNewMessages := 1
	if nbOfNewMessages == 0 {
		msg := fmt.Sprintf("Bonjour %s, pas de nouveau message, bonne journée !", app.checkerFirstName)
		tropoHandler.Say(msg, app.defaultVoice)
		return
	}

	compo := tropoHandler.NewComposer()
	msg := fmt.Sprintf("Bonjour %s, vous avez %d nouveaux messages.", app.checkerFirstName, nbOfNewMessages)
	compo.AddCommand(&tropo.SayCommand{Message:msg, Voice:app.defaultVoice})

	// play first message
	audioFile := "8feadb25a73cac2122bab15ebff58788.wav"
	audioURI:= app.audioServerEndpoint + "/" + audioFile
	//audio := fmt.Sprintf("ftp://%s:%s@ftp.tropo.com/recordings/e238bf666b523148830648743e8df485807310670638548414.wav", )
	compo.AddCommand(&tropo.SayCommand{Message:audioURI})

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

		tropoHandler.Say("Désolé, nous n'avons pas pu enregistrer votre message. Merci de ré essayer !", tropo.VOICE_AUDREY)
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
		tropoHandler.Say("Désolé, nous n'avons pas pu enregistrer votre message. Merci de ré essayer !", tropo.VOICE_AUDREY)
		return
	}

	// say good bye
	tropoHandler.Say("Votre message est bien enregistré. Bonne journée !", tropo.VOICE_AUDREY)
}

func (app *AnsweringMachine) recordingIncompleteHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("RecordingIncompleteHandler")
}

func (app *AnsweringMachine) recordingErrorHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("RecordingErrordHandler")
}


