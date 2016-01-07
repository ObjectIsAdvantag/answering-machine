package tropo

import (
	"fmt"
	"encoding/json"
)



type recordWrapper struct {
	RecordCommand `json:"record"`
}

type RecordCommand struct {
	Beep          string `json:"beep,omitempty"`
	Attempts      int `json:"attempts,omitempty"`
	Bargein       bool `json:"bargein,omitempty"`
	Choices*       struct {
					  Terminator string `json:"terminator,omitempty"`
				  } `json:"choices,omitempty"`
	Maxsilence    int `json:"maxSilence,omitempty"`
	Maxtime       int `json:"maxTime,omitempty"`
	Name          string `json:"name,omitempty"`
	Timeout       int `json:"timeout,omitempty"`
	URL           string `json:"url,omitempty"`
	Asyncupload   string `json:"asyncUpload,omitempty"`
	Transcription* struct {
					  ID  string `json:"id,omitempty"`
					  URL string `json:"url,omitempty"`
				  } `json:"transcription,omitempty"`
}

// Commands interface
func (cmd *RecordCommand) MarshalJSON() ([]byte, error) {

	wrapper := recordWrapper{*cmd}

	b, err := json.Marshal(wrapper)
	if err != nil {
		fmt.Println("cannot encode RecordCommand: ", err)
		return nil, err
	}

	return b, nil
}
