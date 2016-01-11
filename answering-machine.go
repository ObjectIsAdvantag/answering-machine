// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

package main

import (
	"os"
	"flag"
	"strconv"

	"github.com/golang/glog"

	"github.com/ObjectIsAdvantag/answering-machine/server"
	"github.com/ObjectIsAdvantag/answering-machine/machine"
	"github.com/ObjectIsAdvantag/answering-machine/tropo"

	"github.com/paked/configure"
)


const version = "v0.3"

func main() {
	// Read arguments (prevail)
	var showVersion bool
	var port, name, envConfig, messagesConfig string
	flag.StringVar(&port, "port", "8080", "ip port of the server, defaults to 8080")
	flag.StringVar(&name, "name", "Answering Machine", "name of the service, defaults to Answering Machine")
	flag.StringVar(&envConfig, "env", "env.json", "environment configuration file, defaults to env.json, can be overloaded by env variables")
	flag.StringVar(&messagesConfig, "messages", "messages-en.json", "defaults messages, defaults to messages-en.json")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.Parse()

	// Read configuration (env variables then properties, then default values)
	env, messages := readConfiguration(envConfig, messagesConfig)

	if showVersion {
		glog.Infof("SmartProxy version %s\n", version)
		return
	}

	if _, err := strconv.Atoi(port); err != nil {
		glog.Errorf("Invalid port: %s (%s)\n", port, err)
	}

	service := machine.NewAnsweringMachine(env, messages)

	glog.Infof("Starting %s, version: %s\n", name, version)

	if err := server.Run(port, service, version, name); err != nil {
		glog.Errorf("Service exited with error: %s\n", err)
		glog.Flush()
		os.Exit(255)
		return
	}

	glog.Info("Service exited gracefully\n")
	glog.Flush()
}


func readConfiguration(envProperties string, messagesProperties string) (*machine.EnvConfiguration, *machine.I18nMessages) {

	conf := configure.New()
	conf.Use(configure.NewEnvironment())
	if envProperties != "" {
		conf.Use(configure.NewJSONFromFile(envProperties))
	}
	if messagesProperties != "" {
		conf.Use(configure.NewJSONFromFile(messagesProperties))
	}

	checkerPhoneNumber := conf.String("GOLAM_CHECKER_NUMBER", "", "the checker phone number to automate new messages check")
	checkerName := conf.String("GOLAM_CHECKER_NAME", "", "to enhance the welcome message of the new messages checker")
	recorderEndpoint := conf.String("GOLAM_RECORDER_ENDPOINT", "", "to receive the recordings")
	recorderUsername := conf.String("GOLAM_RECORDER_USERNAME", "", "credentials to the recorder endpoint")
	recorderPassword := conf.String("GOLAM_RECORDER_PASSWORD", "", "credentials to the recorder endpoint")
	audioEndpoint := conf.String("GOLAM_AUDIO_ENDPOINT", "", "audio files server")
	transcriptsEmail := conf.String("GOLAM_TRANSCRIPTS_EMAIL", "", "to receive transcripts via email")

	defaultVoice := conf.String("GOLAM_VOICE", "Vanessa", "defaults to English")
	welcome := conf.String("GOLAM_WELCOME", "Welcome, please leave a message after the beep", "to enhance the welcome message of the new messages checker")
	welcomeAlt := conf.String("GOLAM_WELCOME_ALT", "Sorry we do not take any message currently, please call again later", "alternative message if storage service could not be started")
	checkNoMessage := conf.String("GOLAM_CHECK_NO_MESSAGE", "Hello %s, no new messages. Have a good day !", "")
	checkNewMessage := conf.String("GOLAM_CHECK_NEW_MESSAGES", "Hello %s, you have %d new messages", "")
	recordingOK := conf.String("GOLAM_RECORDING_OK", "Your message is recorded. Have a great day !", "")
	recordingFailed := conf.String("GOLAM_RECORDING_FAILED", "Sorry, we could not record your message. Please try again later", "")

	conf.Parse()


	var env machine.EnvConfiguration
	env.RecorderEndpoint = *recorderEndpoint
	env.RecorderUsername= *recorderUsername
	env.RecorderPassword = *recorderPassword
	env.AudioServerEndpoint= *audioEndpoint
	env.TranscriptsReceiver= *transcriptsEmail
	env.CheckerPhoneNumber= *checkerPhoneNumber
	env.CheckerFirstName = *checkerName
	env.DBfilename= "messages.db"

	var messages machine.I18nMessages
	messages.DefaultVoice = tropo.GetVoice(*defaultVoice)
	messages.WelcomeMessage = *welcome
	messages.WelcomeAltMessage = *welcomeAlt
	messages.CheckNewMessages = *checkNewMessage
	messages.CheckNoMessage = *checkNoMessage
	messages.RecordingOKMessage = *recordingOK
	messages.RecordingFailedMessage = *recordingFailed

 	return &env, &messages
}

