package tropo

import (
	"fmt"
	"encoding/json"
)


type sayWrapper struct {
	SayCommand `json:"say"`
}

type SayCommand struct {
	Value string `json:"value"`
	Voice string `json:"voice"`
}

// Commands interface
func (cmd *SayCommand) MarshalJSON() ([]byte, error) {

	wrapper := sayWrapper{*cmd}

	b, err := json.Marshal(wrapper)
	if err != nil {
		fmt.Println("cannot encode SayCommand: ", err)
		return nil, err
	}

	return b, nil
}





