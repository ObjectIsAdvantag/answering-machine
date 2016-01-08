// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

package tropo

import (
	"fmt"
	"encoding/json"
)



type recordCommandWrapper struct {
	RecordCommand `json:"record"`
}

type RecordCommand struct {
	Beep          bool `json:"beep,omitempty"`
	Attempts      int `json:"attempts,omitempty"`
	Bargein       bool `json:"bargein,omitempty"`
	Choices 	  *RecordChoices `json:"choices,omitempty"`
	MaxSilence    int `json:"maxSilence,omitempty"`
	MaxTime       int `json:"maxTime,omitempty"`
	Name          string `json:"name,omitempty"`
	Timeout       int `json:"timeout,omitempty"`
	URL           string `json:"url,omitempty"`
	AsyncUpload   bool `json:"asyncUpload,omitempty"`
	Transcription *RecordTranscription `json:"transcription,omitempty"`
}

type RecordChoices struct {
	Terminator string `json:"terminator,omitempty"`
}

type RecordTranscription struct {
	ID  string `json:"id,omitempty"`
	URL string `json:"url,omitempty"`
}



// Commands interface
func (cmd *RecordCommand) MarshalJSON() ([]byte, error) {

	wrapper := recordCommandWrapper{*cmd}

	b, err := json.Marshal(wrapper)
	if err != nil {
		fmt.Println("cannot encode RecordCommand: ", err)
		return nil, err
	}

	return b, nil
}
