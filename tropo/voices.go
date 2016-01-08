// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

package tropo


import (
	"encoding/json"
)


type Voice struct {
	code 	string `json:"voice"`
	Lang	locale  `json:"-"`
	Gender	gender    `json:"-"`
	Env	 	environment `json:"-"`
}

func (voice *Voice) MarshalJSON() ([]byte, error) {
	return json.Marshal(voice.code)
}

func (voice *Voice) Name() string {
	return voice.code
}

type environment int
const (
	DEV environment = iota
	PROD
)

type locale string
const (
	en_US locale = "en_US"
	fr_CA locale = "fr_CA"
	fr_FR locale = "fr_FR"
)

type gender int
const (
	FEMALE gender = iota
	MALE
)


// Extracted from https://www.tropo.com/docs/webapi/international-features/speaking-multiple-languages (January 2016)
// TODO extend the list of voices
var voices = make(map[string]*Voice)
var VOICE_ALLISON = registerVoice("Allison", en_US, FEMALE, DEV)
var VOICE_AVA = registerVoice("Ava", en_US, FEMALE, DEV)
var VOICE_SAMANTHA = registerVoice("Samantha", en_US, FEMALE, DEV)
var VOICE_SUSAN = registerVoice("Susan", en_US, FEMALE, DEV)
var VOICE_VERONICA = registerVoice("Veronica", en_US, FEMALE, DEV)
var VOICE_VANESSA = registerVoice("Vanessa", en_US, FEMALE, DEV)
var VOICE_DEFAULT = VOICE_VANESSA // Tropo's default
var VOICE_TOM = registerVoice("Tom", en_US, MALE, DEV)
var VOICE_VICTOR = registerVoice("Victor", en_US, MALE, DEV)

var VOICE_AUDREY = registerVoice("Audrey", fr_FR, FEMALE, DEV)
var VOICE_AURELIE = registerVoice("Aurelie", fr_FR, FEMALE, DEV)
var VOICE_THOMAS = registerVoice("Thomas", fr_FR, MALE, DEV)

var VOICE_AMELIE = registerVoice("Amelie", fr_CA, FEMALE, DEV)
var VOICE_CHANTAL = registerVoice("Chantal", fr_CA, FEMALE, DEV)
var VOICE_NICOLAS = registerVoice("Nicolas", fr_CA, MALE, DEV)

func registerVoice(code string, lang locale, g gender, env environment) *Voice {
	voice := &Voice{code, lang, g, env}
	voices[code] = voice
	return voice
}

func GetVoice(code string) *Voice {
	if (code == "") {
		return VOICE_DEFAULT
	}

	return voices[code]
}



