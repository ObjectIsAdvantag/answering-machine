package machine


import (
	"net/http"

	"github.com/golang/glog"
	"github.com/ObjectIsAdvantag/answering-machine/tropo"
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
	app := AnsweringMachine{"Audrey", "/", "/answer", "/timeout", "/error", "http://answeringmachine.localtunnel.me/recordings", "steve.sfartz@gmail.com"}
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
	// tropo.Say("Bienvenue chez Stève, Valérie, Jeanne et Olivia. Bonne année 2016 ! Laissez votre message.", app.Voice)

	//tropo.SendRaw(`{"tropo":[{"record":{"say":[{"value":"Bienvenue chez Stève, Valérie, Jeanne et Olivia. Bonne année 2016 ! Laissez votre message.","voice":"Audrey"},{"event":"timeout","value":"Désolé, nous n'avons pas entendu votre message. Merci de ré-essayer.","voice":"Audrey"}],"name":"foo","url":"https://recording.localtunnel.me/","transcription":{"id":"1234","url":"mailto:steve.sfartz@gmail.com"},"choices":{"terminator":"#"}}}]}`)
	//tropoHandler.SendRawJSON(`{"tropo":[{"say":{"value":"Bienvenue chez Jeanne, Olivia, Stève et Valérie. Bonne année 2016 ! Après le bip c'est à vous...","voice":"Audrey"}},{"record":{"beep":"true","attempts":3,"bargein":false,"choices":{"terminator":"#"},"maxSilence":5,"maxTime":60,"name":"recording","timeout":10,"url":"https://recorder.localtunnel.me/recordings","asyncUpload":"true","transcription":{"id":"1234","url":"mailto:steve.sfartz@gmail.com"}}},{"on":{"event":"continue","next":"/answer","required":true}},{"on":{"event":"incomplete","next":"/timeout","required":true}},{"on":{"event":"error","next":"/error","required":true}}]}`)

	compo := tropoHandler.NewComposer()
	compo.AddCommand(&tropo.SayCommand{"Bienvenue chez Jeanne, Olivia, Stève et Valérie. Bonne année 2016 ! Après le bip c'est à vous...", "Audrey"})
	compo.AddCommand(&tropo.RecordCommand{Beep:true, MaxSilence:5, Timeout:10, MaxTime:60, Name:"recording", URL:"https://recorder.localtunnel.me/recordings"})
	compo.AddCommand(&tropo.OnCommand{Event:"continue", Next:"/answer", Required:true})
	compo.AddCommand(&tropo.OnCommand{Event:"incomplete", Next:"/timeout", Required:true})
	// TODO check if Required could not be changed to false
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
	tropoHandler.Say("Votre message est bien enregistré. Bonne journée !", app.Voice)
}

func (app *AnsweringMachine) recordingIncompleteHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("RecordingIncompleteHandler")
}

func (app *AnsweringMachine) recordingErrorHandler(w http.ResponseWriter, req *http.Request) {
	glog.V(2).Infof("RecordingErrordHandler")
}





