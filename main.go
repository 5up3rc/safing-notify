// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the GPL license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	_ "github.com/Safing/safing-notify/meta"
)

// var _ = runtime.GOMAXPROCS(16)

// func init() {
// 	runtime.GOMAXPROCS(4)
// }

func main() {

	// catch interrupt for clean shutdown
	signalCh := make(chan os.Signal)
	signal.Notify(
		signalCh,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	<-signalCh
	log.Println("program was interrupted, shutting down.")
	// TODO: do a clean shutdown
	os.Exit(0)

}

func StartUI() {

	var path string

	// check if we find the binary in the directory of the notify executable
	expath, err := os.Executable()
	if err == nil {
		possiblePath := filepath.Join(filepath.Dir(expath), "safing-ui")
		if fileExists(possiblePath) {
			path = possiblePath
		}
	}

	// check if we find the binary in the current working directory
	if path != "" {
		possiblePath, err := filepath.Abs("./safing-ui")
		if err == nil && fileExists(possiblePath) {
			path = possiblePath
		}
	}

	if path != "" {
		log.Printf("starting Safing UI from %s", path)
		cmd := exec.Command(path)
		err := cmd.Start()
		if err != nil {
			log.Printf("failed to start Safing UI: %s", err)
			Notify(ErrorLevel, "failed to start Safing UI application.")
		}
	} else {
		Notify(ErrorLevel, "could not find Safing UI application.")
	}

}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
