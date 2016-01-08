package machine


import (
	"net/http"

	"github.com/golang/glog"
	"github.com/ObjectIsAdvantag/answering-machine/tropo"
)


type AnsweringMachine struct {
	Voice						*tropo.TropoVoice // see https://www.tropo.com/docs/webapi/international-features/speaking-multiple-languages
	WelcomeMessageRoute			string // route to the welcome message
	SuccessRoute				string // invoked after message are recorded
	IncompleteRoute				string // invoked if a timeout occurs
	CommunicationErrorRoute		string // invoked if the recording failed due to communication issues between Tropo and the AnsweringMachine
	RecorderEndpoint			string // URI to record the messages
	TranscriptsReceiver			string // email of the transcriptions receiver
}


func NewAnsweringMachine() *AnsweringMachine {
	app := AnsweringMachine{tropo.VOICE_AUDREY, "/", "/answer", "/timeout", "/error", "http://answeringmachine.localtunnel.me/recordings", "mailto:steve.sfartz@gmail.com"}
	return &app
}

func (app *AnsweringMachine) RegisterHandlers() {
	http.HandleFunc(app.WelcomeMessageRoute, app.welcomeHandler)
	http.HandleFunc(app.SuccessRoute, app.recordingSuccessHandler)
	http.HandleFunc(app.IncompleteRoute, app.recordingIncompleteHandler)
	http.HandleFunc(app.CommunicationErrorRoute, app.recordingErrorHandler)
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

	// please leave a message, start recording
	compo := tropoHandler.NewComposer()
	compo.AddCommand(&tropo.SayCommand{Message:"Bienvenue chez Jeanne, Olivia, Stève et Valérie. Bonne année 2016 ! Après le bip c'est à vous...", Voice:tropo.VOICE_AUDREY})
	choices := tropo.RecordChoices{Terminator:"#"}
	transcript := tropo.RecordTranscription{ID:session.CallID, URL:app.TranscriptsReceiver}

	compo.AddCommand(&tropo.RecordCommand{Bargein:true, Attempts:3, Beep:true, Choices:&choices, MaxSilence:3, Timeout:10, MaxTime:60, Name:"recording", URL:"https://recorder.localtunnel.me/recordings", AsyncUpload:true, Transcription:&transcript})
	compo.AddCommand(&tropo.OnCommand{Event:"continue", Next:"/answer", Required:true})
	compo.AddCommand(&tropo.OnCommand{Event:"incomplete", Next:"/timeout", Required:true})
	compo.AddCommand(&tropo.OnCommand{Event:"error", Next:"/error", Required:true})

	tropoHandler.ExecuteComposer(compo)
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


