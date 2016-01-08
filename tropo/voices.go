// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

package tropo


import (
	"encoding/json"
)


type Voice struct {
	name 	string `json:"voice"`
	Lang	language  `json:"-"`
	Gender	gender    `json:"-"`
	Env	 	environment `json:"-"`
}

func (voice *Voice) MarshalJSON() ([]byte, error) {
	return json.Marshal(voice.name)
}

func (voice *Voice) Name() string {
	return voice.name
}

type environment int
const (
	DEV environment = iota
	PROD
)

type language string
const (
	fr_FR language = "fr_FR"
)

type gender int
const (
	FEMALE gender = iota
	MALE
)


// TODO : extend the list of voices for each env : Dev / Prod
// see https://www.tropo.com/docs/webapi/international-features/speaking-multiple-languages
var VOICE_AUDREY = &Voice{"Audrey", fr_FR, FEMALE, DEV}




