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

	"github.com/paked/configure"
)


const version = "v0.3"

func main() {
	// Read arguments (prevail)
	var showVersion bool
	var port, name, properties string
	flag.StringVar(&port, "port", "8080", "ip port of the server, defaults to 8080")
	flag.StringVar(&name, "name", "Answering Machine", "name of the service, defaults to Answering Machine")
	flag.StringVar(&properties, "conf", "env-tropofs.json", "answering machine configuration filename")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.Parse()

	// Read configuration (env variables then properties, then default values)
	env, routes, messages := readConfiguration(properties)

	if showVersion {
		glog.Infof("SmartProxy version %s\n", version)
		return
	}

	if _, err := strconv.Atoi(port); err != nil {
		glog.Errorf("Invalid port: %s (%s)\n", port, err)
	}

	service := machine.NewAnsweringMachine(env, routes, messages)

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


func readConfiguration(properties string) (env *machine.EnvConfiguration, routes *machine.HandlerRoutes, messages *machine.I18nMessages) {

	var env machine.EnvConfiguration
	var messages machine.I18nMessages
	var routes machine.HandlerRoutes

	conf := configure.New()
	conf.Use(configure.NewEnvironment())
	if properties != "" {
		conf.Use(configure.NewJSONFromFile(properties))
	}

	checkerPhoneNumber := conf.String("GOLAM_CHECKER_NUMBER", "", "the checker phone number to automate new messages check")
	checkerName := conf.String("GOLAM_CHECKER_NAME", "", "to enhance the welcome message of the new messages checker")
	recorderEndpoint := conf.String("GOLAM_RECORDER_ENDPOINT", "", "to receive the recordings")
	recorderUsername := conf.String("GOLAM_RECORDER_USERNAME", "", "credentials to the recorder endpoint")
	recorderPassword := conf.String("GOLAM_RECORDER_PASSWORD", "", "credentials to the recorder endpoint")
	audioEndpoint := conf.String("GOLAM_AUDIO_ENDPOINT", "", "audio files server")
	transcriptsEmail := conf.String("GOLAM_TRANSCRIPTS_EMAIL", "", "to receive transcripts via email")

	conf.Parse()




		env.recorderEndpoint
	env.recorderUsername
	env.recorderPassword
	env.audioServerEndpoint
	env.transcriptsReceiver
	env.checkerPhoneNumber
	env.checkerFirstName
	env.dbFilename


		routes.welcomeMessageRoute
	routes.recordingSuccessRoute
	routes.recordingIncompleteRoute
	routes.recordingFailedRoute
	routes.adminRoute


 	return &env, &routes, &messages
}

