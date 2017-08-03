// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the GPL license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/Safing/safing-core/formats/dsd"
)

const (
	backOffTimer = 1 * time.Second
)

var (
	sendToCore = make(chan []byte, 10)

	sessionKey string
	lastError  string
)

func init() {
	go connector()
}

func connector() {
	log.Println("connecting to Safing Core")
	for {
		connectionHandler()
		displayLevel(SecurityLevelUnknown)
	}
}

func connectionHandler() {

	wsConn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:18/api/v1", nil)
	if err != nil {
		if err.Error() != lastError {
			log.Printf("error connecting to safing core: %s\n", err)
		}
		lastError = err.Error()
		time.Sleep(backOffTimer)
		return
	}
	defer wsConn.Close()

	failed := make(chan interface{}, 0)
	readC := make(chan []byte, 10)
	writeC := make(chan []byte, 10)

	go reader(wsConn, readC, failed)
	go writer(wsConn, writeC, failed)

	// clean up when handler dies
	defer func() {
		select {
		case <-failed:
		default:
			close(failed)
		}
		close(writeC)
		log.Println("lost connection to Safing Core, reconnecting...")
	}()

	// initialize connection

	created := true
	if sessionKey != "" {
		writeC <- []byte("resume|" + sessionKey)
	} else {
		writeC <- []byte("start")
	}

	message := string(<-readC)
	if !strings.HasPrefix(message, "session|") {
		log.Printf("could not start session: %s\n", message)
		sessionKey = ""
		return
	}
	sessionKey = strings.SplitN(message, "|", 2)[1]
	if created {
		log.Printf("session started.\n")
	} else {
		log.Printf("session resumed.\n")
	}
	defer log.Printf("session ended.\n")

	// subscriptions
	writeC <- []byte("subscribe|/Me/SystemStatus:status")

	// start handling
	for {
		select {
		case <-failed:
			return
		case msg := <-readC:
			log.Printf("got message: %s\n", string(msg))
			handleMessage(msg)
		case msg := <-sendToCore:
			log.Printf("sending message: %s\n", string(msg))
			writeC <- msg
		}
	}

}

func reader(wsConn *websocket.Conn, readC chan []byte, failed chan interface{}) {
	for {
		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("read error: %s", err)
			} else {
				log.Printf("socket closed by server.")
			}
			close(failed)
			return
		}
		readC <- msg
	}
}

func writer(wsConn *websocket.Conn, writeC chan []byte, failed chan interface{}) {
	for {
		select {
		case <-failed:
			return
		case msg := <-writeC:
			err := wsConn.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				log.Printf("write error: %s\n", err)
				close(failed)
				return
			}
		}
	}
}

func handleMessage(msg []byte) {
	splitted := bytes.SplitN(msg, []byte("|"), 2)
	action, data := string(splitted[0]), splitted[1]
	switch action {
	case "current", "created", "updated":
		splitted := bytes.SplitN(data, []byte("|"), 2)
		key, object := string(splitted[0]), splitted[1]
		switch key {
		case "/Me/SystemStatus:status":
			var newStatus SystemStatus
			_, err := dsd.Load(&object, &newStatus)
			if err != nil {
				// TODO: handle error better
				log.Printf("failed to parse %s: %s", key, err)
				return
			}
			HandleNewStatus(newStatus)
		}
	case "notify":
		splitted := strings.SplitN(string(data), "|", 2)
		levelString, message := splitted[0], splitted[1]
		level, err := strconv.ParseInt(levelString, 10, 8)
		if err != nil {
			// TODO: handle error better
			log.Printf("failed to parse %s: %s", levelString, err)
			return
		}
		Notify(level, message)
	}
}
