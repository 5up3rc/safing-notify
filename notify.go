// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the GPL license that can be found in the LICENSE file.

package main

import (
	"fmt"

	notify "github.com/TheCreeper/go-notify"

	"github.com/Safing/safing-notify/icons"
)

const (
	TraceLevel    int64 = 1
	DebugLevel    int64 = 2
	InfoLevel     int64 = 3
	WarningLevel  int64 = 4
	ErrorLevel    int64 = 5
	CriticalLevel int64 = 6
)

type NotificationIcon struct {
	Width         int
	Height        int
	Rowstride     int
	HasAlpha      bool
	BitsPerSample int
	NChannels     int
	Pixels        []byte
}

func Notify(level int64, message string) {
	n := notify.NewNotification("Safing", message)
	n.Hints = make(map[string]interface{})
	n.Hints[notify.HintUrgency] = min(0, int(level)-4)

	currentStatusLock.Lock()
	chosenIcon := icons.DataByLevel[currentStatus.CurrentSecurityLevel]
	currentStatusLock.Unlock()

	nChannels := 3
	if chosenIcon.HasAlpha {
		nChannels = 4
	}
	notifIcon := NotificationIcon{
		chosenIcon.Width,
		chosenIcon.Height,
		chosenIcon.RowStride,
		chosenIcon.HasAlpha,
		chosenIcon.BitsPerSample,
		nChannels,
		chosenIcon.Data,
	}
	n.Hints["icon_data"] = notifIcon

	n.Show()
}

func LogNotify(level int64, message string) {
	switch level {
	case WarningLevel:
		Notify(WarningLevel, fmt.Sprintf("encountered a warning: %s", message))
	case ErrorLevel:
		Notify(ErrorLevel, fmt.Sprintf("encountered an error: %s", message))
	case CriticalLevel:
		Notify(CriticalLevel, fmt.Sprintf("encountered a critical error: %s", message))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
