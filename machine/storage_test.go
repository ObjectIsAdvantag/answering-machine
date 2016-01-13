// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

package machine

import (
	"testing"
)

func TestStorage_create (t *testing.T) {
	_, err := NewStorage("tests.db", true)
	if err != nil {
		t.Failed()
	}

}

