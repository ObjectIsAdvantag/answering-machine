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
)


const version = "v0.5"

func main() {
	// Read arguments (prevail)
	var showVersion bool
	var port, name, envConfig, messagesConfig string
	flag.StringVar(&port, "port", "8080", "ip port of the server, defaults to 8080")
	flag.StringVar(&name, "name", "AnsweringMachine", "name of the service, defaults to AnsweringMachine")
	flag.StringVar(&envConfig, "env", "env.json", "environment configuration file")
	flag.StringVar(&messagesConfig, "messages", "messages-en.json", "defaults messages, defaults to messages-en.json")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.Parse()

	// Read configuration (env variables then properties, then default values)
	env, messages := machine.ReadConfiguration(envConfig, messagesConfig)

	if showVersion {
		glog.Infof("%s version %s\n", name, version)
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


