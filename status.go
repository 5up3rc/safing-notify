// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the GPL license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/Safing/safing-core/formats/dsd"
	"github.com/tevino/abool"
)

const (
	SecurityLevelUnknown  int8 = -1
	SecurityLevelOff      int8 = 0
	SecurityLevelDynamic  int8 = 1
	SecurityLevelSecure   int8 = 2
	SecurityLevelFortress int8 = 3
)

var (
	currentStatus     SystemStatus
	currentStatusLock sync.Mutex
	initialStatus     = abool.NewBool(true)
)

// SystemStatus saves basic information about the current system status.
type SystemStatus struct {
	CurrentSecurityLevel  int8
	SelectedSecurityLevel int8
	ThreatLevel           int8   `json:",omitempty" bson:",omitempty"`
	ThreatReason          string `json:",omitempty" bson:",omitempty"`
}

func (status *SystemStatus) FmtCurrentLevel() string {
	return fmtLevel(status.CurrentSecurityLevel)
}

func (status *SystemStatus) FmtSelectedLevel() string {
	return fmtLevel(status.SelectedSecurityLevel)
}

func fmtLevel(level int8) string {
	switch level {
	case SecurityLevelOff:
		return "Off"
	case SecurityLevelDynamic:
		return "Dynamic"
	case SecurityLevelSecure:
		return "Secure"
	case SecurityLevelFortress:
		return "Fortress"
	}
	return "Unknown"
}

func HandleNewStatus(newStatus SystemStatus) {
	currentStatusLock.Lock()
	defer currentStatusLock.Unlock()

	// change displayed level
	displayLevel(newStatus.CurrentSecurityLevel)

	if !initialStatus.IsSet() {
		// to display a threat warning, the following conditions must be met:
		// (1) ThreatReason must differ from last status
		// (2) SelectedSecurityLevel must be less than or equal to ThreatLevel
		if newStatus.ThreatReason != currentStatus.ThreatReason {
			if newStatus.ThreatReason == "" {

				notify := "All detected attacks ended.\nYou are safe."
				if newStatus.CurrentSecurityLevel == newStatus.SelectedSecurityLevel && currentStatus.CurrentSecurityLevel != currentStatus.SelectedSecurityLevel {
					notify += fmt.Sprintf("\nReturned to selected Level %s", newStatus.FmtCurrentLevel())
				}
				go Notify(InfoLevel, notify)

			} else {

				notify := fmt.Sprintf("Safing %s.", newStatus.ThreatReason)
				if newStatus.ThreatLevel >= SecurityLevelFortress {
					notify += "As Safing cannot handle this threat, please call for help."
					go Notify(CriticalLevel, notify)
				} else {
					if newStatus.CurrentSecurityLevel > currentStatus.CurrentSecurityLevel {
						notify += fmt.Sprintf("\nLevel raised to %s to keep you safe.", newStatus.FmtCurrentLevel())
					} else {
						notify += "\nYou are safe."
					}
					go Notify(WarningLevel, notify)
				}

			}
		}
	} else {
		initialStatus.Set()
	}

	currentStatus = newStatus
}

func GetThreatMessage() string {
	currentStatusLock.Lock()
	defer currentStatusLock.Unlock()
	if currentStatus.ThreatLevel < 1 {
		return "You are safe."
	}
	if currentStatus.ThreatLevel < 3 {
		var recommendedLevel string
		switch currentStatus.ThreatLevel {
		case SecurityLevelDynamic:
			recommendedLevel = "Secure or Fortress"
		case SecurityLevelSecure:
			recommendedLevel = "Fortress"
		}
		return fmt.Sprintf("Safing %s and recommends to stay at Level %s", currentStatus.ThreatReason, recommendedLevel)
	}
	return fmt.Sprintf("Safing %s and recommends you immediately shutdown your computer, call for help and get some coffee for whoever will have to clean up this mess. Now.", currentStatus.ThreatReason)
}

func SetLevel(level int8) {
	currentStatusLock.Lock()
	currentStatus.CurrentSecurityLevel = level
	currentStatus.SelectedSecurityLevel = level
	currentStatusLock.Unlock()
	sendToCore <- PackStatus()
}

func PackStatus() []byte {
	currentStatusLock.Lock()
	defer currentStatusLock.Unlock()
	dumped, err := dsd.Dump(currentStatus, dsd.AUTO)
	if err != nil {
		log.Printf("failed to pack current status: %s", err)
		return []byte{}
	}
	return append([]byte("update|/Me/SystemStatus:status|"), *dumped...)
}
