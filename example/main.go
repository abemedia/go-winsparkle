//go:build windows
// +build windows

package main

import (
	"log"
	"time"

	"github.com/abemedia/go-winsparkle"
	_ "github.com/abemedia/go-winsparkle/dll"
)

func main() {
	log.Println("starting app")

	winsparkle.SetAppcastURL("https://winsparkle.org/example/appcast.xml")
	winsparkle.SetAppDetails("winsparkle.org", "WinSparkle Go Example", "1.0.0")

	if err := winsparkle.SetEdDSAPublicKey("payYa5ap0XtF8HWR4AYBdCIcXWtJZPen7bJqFcqlp7o="); err != nil {
		log.Fatal(err)
	}

	c := make(chan struct{})

	winsparkle.SetShutdownRequestCallback(func() {
		log.Println("installing update")
		close(c)
	})

	winsparkle.SetUpdateCancelledCallback(func() {
		log.Println("cancelled update")
		close(c)
	})

	winsparkle.Init()
	defer winsparkle.Cleanup()

	winsparkle.CheckUpdateWithUI()

	// waits until update is installed or cancelled (10min timeout)
	select {
	case <-c:
	case <-time.After(10 * time.Minute):
	}
	log.Println("shutting down")
}
