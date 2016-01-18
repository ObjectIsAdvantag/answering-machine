package machine


import (
	"os"

	"github.com/golang/glog"
	"github.com/paked/configure"
)


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

func LoadEnvConfiguration(filename string) (*EnvConfiguration) {

	conf := configure.New()
	glog.V(0).Infof("Loading configuration from 1. shell variables")
	conf.Use(configure.NewEnvironment())
	if filename != "" {
		glog.V(0).Infof("Loading configuration from 2. json file: %s", filename)
		if file, err := os.Open(filename); err != nil {
			// TODO What if we wanted to pass all configuration via env variables, not possible as of today
			glog.Fatalf("Could not open Environment configuration file: %s", filename)
		} else {
			file.Close()
			conf.Use(configure.NewJSONFromFile(filename))
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

	return &env
}