package machine

import (
	"net/http"

	"github.com/golang/glog"
	api "github.com/ObjectIsAdvantag/answering-machine/tropo"
)


type AnsweringMachine struct {
	Voice						string // see https://www.tropo.com/docs/webapi/international-features/speaking-multiple-languages
	WelcomeMessageRoute			string // route to the welcome message
	SuccessRoute				string // invoked after message are recorded
	IncompleteRoute				string // invoked if a timeout occurs
	CommunicationErrorRoute		string // invoked if the recording failed due to communication issues between Tropo and the AnsweringMachine
	RecorderEndpoint			string // URI to record the messages
	TranscriptsReceiver			string // email of the transcriptions receiver
}


func NewAnsweringMachine() *AnsweringMachine {
	app := AnsweringMachine{"Audrey", "/", "/answer", "/timeout", "/error", "http://recorder.localtunnel.me/recordings", "steve.sfartz@gmail.com"}
	return &app
}

func (app *AnsweringMachine) RegisterHandlers() {
	http.HandleFunc(app.WelcomeMessageRoute, app.welcomeHandler)
	http.HandleFunc(app.SuccessRoute, app.recordingSuccessHandler)
	http.HandleFunc(app.IncompleteRoute, app.recordingIncompleteHandler)
	http.HandleFunc(app.CommunicationErrorRoute, app.recordingErrorHandler)
}

func (app *AnsweringMachine) welcomeHandler(w http.ResponseWriter, req *http.Request) {
	tropo := api.NewDriver(w, req)

	glog.V(2).Infof("Incoming call")
	var session *api.Session
	var err error
    if session, err = tropo.ReadSession(); err != nil {
		tropo.ReplyInternalError()
		return
	}

	// check a human issued the call
	if !(session.IsHumanInitiated() && session.IsCall()) {
		glog.V(1).Infof("Unsupported request, a voice call is expected\n")
		tropo.ReplyBadInput()
		return
	}

	caller := session.From.ID
	glog.V(0).Infof(`SessionID "%s", CallID "%s", From "+%s"`, session.ID, session.CallID, caller)

	// echo leave a message

	// tropo.Say("Bienvenue chez Stève, Valérie, Jeanne et Olivia. Bonne année 2016 ! Laissez votre message.", app.Voice)

	// TODO Create higher level library
	//tropo.SendRaw(`{"tropo":[{"record":{"say":[{"value":"Bienvenue chez Stève, Valérie, Jeanne et Olivia. Bonne année 2016 ! Laissez votre message.","voice":"Audrey"},{"event":"timeout","value":"Désolé, nous n'avons pas entendu votre message. Merci de ré-essayer.","voice":"Audrey"}],"name":"foo","url":"https://recording.localtunnel.me/","transcription":{"id":"1234","url":"mailto:steve.sfartz@gmail.com"},"choices":{"terminator":"#"}}}]}`)
	tropo.SendRaw(`{"tropo":[{"say":{"value":"Bienvenue chez Stève, Valérie, Jeanne et Olivia. Bonne année 2016 !","voice":"Audrey"}},{"record":{"attempts":3,"bargein":false,"choices":{"terminator":"#"},"maxSilence":5,"maxTime":60,"name":"recording","say":{"value":"Laissez votre message après le bip","voice":"Audrey"},"timeout":10,"url":"https://recordings.localtunnel.me/recordings","transcription":{"id":"1234","url":"mailto:steve.sfartz@gmail.com"}}},{"on":{"event":"continue","next":"/answer","required":true}},{"on":{"event":"incomplete","next":"/timeout","required":true}},{"on":{"event":"error","next":"/error","required":true}}]}`)
}

func (app *AnsweringMachine) recordingSuccessHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("RecordingSuccessHandler")
}

func (app *AnsweringMachine) recordingIncompleteHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("RecordingIncompleteHandler")
}

func (app *AnsweringMachine) recordingErrorHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("RecordingErrordHandler")
}





