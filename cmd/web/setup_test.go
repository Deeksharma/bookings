package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// this part will be executed before you run the test
	// use this to make variables that are needed in tests using interfaces
	os.Exit(m.Run()) // exit after running the tests
}

type myHandler struct {} // implement all the functions that are implemented by handlers

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){

}
