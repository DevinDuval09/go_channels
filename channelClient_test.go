package main

import (
	"testing"
	"time"
)

func TestClientAndServer(t *testing.T) {
	go runTwoResponseServer()
	time.Sleep(10000 * time.Millisecond)
	main()
	want := true
	actual := false
	if actual != want {
		t.Fatal("Test")
	}

}
