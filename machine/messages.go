package machine

import (
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

//
func LoadMessagesConfiguration(filename string) (*I18nMessages) {

	conf := configure.New()
	glog.V(0).Infof("Loading env preferences from 1. env variables")
	conf.Use(configure.NewEnvironment())
	if filename != "" {
		glog.V(0).Infof("Loading messages from: %s", filename)
		if file, err := os.Open(filename); err != nil {
			glog.Warningf("Could not open Messages definition, switching to default values", filename)
		} else {
			file.Close()
			conf.Use(configure.NewJSONFromFile(filename))
		}
	}

	// messages
	defaultVoice := conf.String("GOLAM_VOICE", "Vanessa", "defaults to English")
	welcome := conf.String("GOLAM_WELCOME", "Welcome, please leave a message after the beep", "to enhance the welcome message of the new messages checker")
	welcomeAlt := conf.String("GOLAM_WELCOME_ALT", "Sorry we do not take any message currently, please call again later", "alternative message if storage service could not be started")
	checkNoMessage := conf.String("GOLAM_CHECK_NO_MESSAGE", "Hello %s, no new messages. Have a good day !", "")
	checkNewMessage := conf.String("GOLAM_CHECK_NEW_MESSAGES", "Hello %s, you have %d new messages", "")
	recordingOK := conf.String("GOLAM_RECORDING_OK", "Your message is recorded. Have a great day !", "")
	recordingFailed := conf.String("GOLAM_RECORDING_FAILED", "Sorry, we could not record your message. Please try again later", "")

	conf.Parse()

	var messages I18nMessages
	messages.DefaultVoice = tropo.GetVoice(*defaultVoice)
	messages.WelcomeMessage = *welcome
	messages.WelcomeAltMessage = *welcomeAlt
	messages.CheckNewMessages = *checkNewMessage
	messages.CheckNoMessage = *checkNoMessage
	messages.RecordingOKMessage = *recordingOK
	messages.RecordingFailedMessage = *recordingFailed

	return &messages
}