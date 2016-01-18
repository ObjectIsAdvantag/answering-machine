// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

// Entrypoint for google app engine
package gae


import (
	"github.com/ObjectIsAdvantag/answering-machine/machine"
)

func init() {

	env := machine.GetDefaultConfigurationBackedWithTropoFS("ObjectIsAdvantag", "XXXXXX", "5048353", "Stève", "33678007899", "steve.sfartz@gmail.com")
	env.DBfilename = ""
	messages := machine.GetDefaultMessages

	service := machine.NewAnsweringMachine(env, messages)
	service.RegisterHandlers()
}