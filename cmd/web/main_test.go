package main

import "testing"

func TestRun(t *testing.T) {
	_, err := Run()
	if err != nil {
		t.Error("failed run")
	}
}
