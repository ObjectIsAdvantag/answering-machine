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


const version = "v0.2.1"

func main() {
	// Read arguments (prevail)
	var showVersion bool
	var port, name, properties string
	flag.StringVar(&port, "port", "8080", "ip port of the server, defaults to 8080")
	flag.StringVar(&name, "name", "Answering Machine", "name of the service, defaults to Answering Machine")
	flag.StringVar(&properties, "conf", "config-tropofs.json", "answering machine configuration filename")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.Parse()

	// Read configuration (env variables then properties, then default values)
	conf := configure.New()
	conf.Use(configure.NewEnvironment())
	if properties != "" {
		conf.Use(configure.NewJSONFromFile(properties))
	}
	welcome := conf.String("GOLAM_WELCOME_MESSAGE", "Welcome, Leave a message after the bip.", "your welcome message")
	welcomeAlt := conf.String("GOLAM_WELCOME_ALT_MESSAGE", "Sorry we do not take any message currently, please call again later", "the alt welcome message when recording is not active")
	voiceCode := conf.String("GOLAM_VOICE", "Vanessa", "Machine's default message for Text To Speach")
	checkerPhoneNumber := conf.String("GOLAM_CHECKER_NUMBER", "", "the checker phone number to automate new messages check")
	checkerName := conf.String("GOLAM_CHECKER_NAME", "", "to enhance the welcome message of the new messages checker")
	recorderEndpoint := conf.String("GOLAM_RECORDER_ENDPOINT", "", "to receive the recordings")
	recorderUsername := conf.String("GOLAM_RECORDER_USERNAME", "", "credentials to the recorder endpoint")
	recorderPassword := conf.String("GOLAM_RECORDER_PASSWORD", "", "credentials to the recorder endpoint")
	audioEndpoint := conf.String("GOLAM_AUDIO_ENDPOINT", "", "audio files server")
	transcriptsEmail := conf.String("GOLAM_TRANSCRIPTS_EMAIL", "", "to receive transcripts via email")
	conf.Parse()

	if showVersion {
		glog.Infof("SmartProxy version %s\n", version)
		return
	}

	if _, err := strconv.Atoi(port); err != nil {
		glog.Errorf("Invalid port: %s (%s)\n", port, err)
	}

	service := machine.NewAnsweringMachine(
		*welcome,
		*welcomeAlt,
		tropo.GetVoice(*voiceCode),
		*recorderEndpoint,
		*recorderUsername,
		*recorderPassword,
		*audioEndpoint,
		*transcriptsEmail,
		*checkerPhoneNumber,
		*checkerName)

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
