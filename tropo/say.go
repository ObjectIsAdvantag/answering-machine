package tropo

import (
	"fmt"
	"encoding/json"
)


type sayCommandWrapper struct {
	SayCommand `json:"say"`
}

type SayCommand struct {
	Message string `json:"value"`
	Voice 	*TropoVoice `json:"voice,omitempty"`
}

// Commands interface
func (cmd *SayCommand) MarshalJSON() ([]byte, error) {

	wrapper := sayCommandWrapper{*cmd}

	b, err := json.Marshal(wrapper)
	if err != nil {
		fmt.Println("cannot encode SayCommand: ", err)
		return nil, err
	}

	return b, nil
}





