// Copyright 2015, Stève Sfartz
// Licensed under the MIT License
package machine

import (
	"testing"
)

func TestStorage_create (t *testing.T) {
	_, err := CreateStorage("messages.db")
	if err != nil {
		t.Failed()
	}
}

/*
func TestStorage_create (t *testing.T) {
	storage, err := CreateStorage("messages.db")
	if err != nil {
		t.Failed()
	}
	msg := storage.CreateTrace()
	fmt.Printf("trace created successfully with id " + trace.ID)//storage.store(trace)
}
*/

