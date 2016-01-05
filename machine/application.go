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
	app := AnsweringMachine{"Audrey", "/", "/success", "/incomplete", "/error", "http://recorder.localtunnel.me/recordings", "steve.sfartz@gmail.com"}
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

	// check a human is calling
	if session.UserType != "HUMAN" || session.From.Channel != "VOICE" {
		glog.V(1).Infof("Unsupported incoming request: %s\n", req.Method)
		tropo.ReplyBadInput()
		return
	}

	// echo leave a message
	number := session.From.ID
	glog.V(0).Infof(`SessionID "%s", CallID "%s", From "+%s"`, session.ID, session.CallID, number)
	tropo.Say("Bienvenue chez Stève, Valérie, Jeanne et Olivia. Bonne année 2016 ! Laissez votre message.", app.Voice)
//	fmt.Fprintf(w, `{"tropo":[{"record":{"say":[{"value":"Bienvenue chez Stève, Valérie, Jeanne et Olivia. Bonne année 2016 ! Laissez votre message.","voice":"Audrey"},{"event":"timeout","value":"Désolé, nous n'avons pas entendu votre message. Merci de ré-essayer.","voice":"Audrey"}],"name":"foo","url":"https://recording.localtunnel.me/","transcription":{"id":"1234","url":"mailto:steve.sfartz@gmail.com"},"choices":{"terminator":"#"}}}]}`)

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





