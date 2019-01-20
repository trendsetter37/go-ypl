package main

import (
	"errors"
	"io"
	"testing"
)

func TestCheck(t *testing.T) {
	eof := io.EOF

	if check(eof) != true {
		t.Fail()
	}
	if check(nil) != false {
		t.Fail()
	}
}

func TestCheckPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Check function did not panic")
		}
	}()
	check(errors.New("Random legit error"))
}
