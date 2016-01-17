// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

// Entrypoint for google app engine
package gae


import (

	"github.com/ObjectIsAdvantag/answering-machine/machine"
)

func init() {
	env, messages := machine.ReadConfiguration("env.json", "messages-en.json")

	service := machine.NewAnsweringMachine(env, messages)
	service.RegisterHandlers()
}