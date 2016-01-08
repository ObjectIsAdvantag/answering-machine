// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

package tropo

import (
	"fmt"
	"encoding/json"
)


type onCommandWrapper struct {
	OnCommand `json:"on"`
}

type OnCommand struct {
	Event string `json:"event"`
	Next string `json:"next"`
	Required bool `json:"required,omitempty"`
}

// Commands interface
func (cmd *OnCommand) MarshalJSON() ([]byte, error) {

	wrapper := onCommandWrapper{*cmd}

	b, err := json.Marshal(wrapper)
	if err != nil {
		fmt.Println("cannot encode RecordCommand: ", err)
		return nil, err
	}

	return b, nil
}
