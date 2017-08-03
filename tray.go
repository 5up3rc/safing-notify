// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the GPL license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/tevino/abool"

	"github.com/Safing/safing-notify/icons"
)

var (
	systrayIcon *gtk.StatusIcon
	gtkLock     sync.Mutex

	rotation = abool.NewBool(false)
)

func init() {
	go tray()
}

func tray() {
	gtkLock.Lock()

	// log.Println("initializing GTK")
	gtk.Init(nil)
	glib.SetApplicationName("Safing Notify")

	systrayIcon = gtk.NewStatusIcon()
	systrayIcon.SetTitle("Safing Notify")
	systrayIcon.SetTooltipMarkup("Safing Notify displays status information about Safing")
	systrayIcon.Connect("popup-menu", func(cbx *glib.CallbackContext) {
		displayMenu(cbx)
	})
	systrayIcon.Connect("activate", func(cbx *glib.CallbackContext) {
		StartUI()
	})

	displayLevel(SecurityLevelUnknown)
	time.Sleep(10 * time.Millisecond)

	gtkLock.Unlock()
	// log.Println("finished GTK init, entering GTK main...")

	gtk.Main()
}

func displayLevel(level int8) {
	if level < 0 {
		go rotateIcons()
		return
	}
	if level < 4 {
		log.Printf("displaying level: %d", level)
		gtkLock.Lock()
		defer gtkLock.Unlock()
		rotation.UnSet()
		systrayIcon.SetFromPixbuf(icons.ByLevel[level])
	}
}

func rotateIcons() {
	if !rotation.SetToIf(false, true) {
		return
	}
	time.Sleep(1 * time.Second)
	for i := int8(0); ; i = (i + 1) % 4 {
		gtkLock.Lock()
		if !rotation.IsSet() {
			gtkLock.Unlock()
			break
		}
		systrayIcon.SetFromPixbuf(icons.ByLevel[i])
		gtkLock.Unlock()
		time.Sleep(500 * time.Millisecond)
	}
}

func displayMenu(cbx *glib.CallbackContext) {
	gtkLock.Lock()
	defer gtkLock.Unlock()

	menu := gtk.NewMenu()

	menuItemOpen := gtk.NewMenuItemWithLabel("Open Safing UI")
	menuItemOpen.Connect("activate", func() {
		StartUI()
	})
	menu.Append(menuItemOpen)

	menu.Append(gtk.NewSeparatorMenuItem())

	// lock for activating correct entry
	currentStatusLock.Lock()

	levelList := &glib.SList{}

	menuItemLevelDynamic := gtk.NewRadioMenuItemWithLabel(levelList, "Level Dynamic")
	menuItemLevelDynamic.Connect("activate", func() {
		SetLevel(SecurityLevelDynamic)
	})
	if currentStatus.CurrentSecurityLevel == SecurityLevelDynamic {
		menuItemLevelDynamic.Select()
	}
	menu.Append(menuItemLevelDynamic)

	menuItemLevelSecure := gtk.NewRadioMenuItemWithLabel(levelList, "Level Secure")
	menuItemLevelSecure.Connect("activate", func() {
		SetLevel(SecurityLevelSecure)
	})
	if currentStatus.CurrentSecurityLevel == SecurityLevelSecure {
		menuItemLevelSecure.Select()
	}
	menu.Append(menuItemLevelSecure)

	menuItemLevelFortress := gtk.NewRadioMenuItemWithLabel(levelList, "Level Fortress")
	menuItemLevelFortress.Connect("activate", func() {
		SetLevel(SecurityLevelFortress)
	})
	if currentStatus.CurrentSecurityLevel == SecurityLevelFortress {
		menuItemLevelFortress.Select()
	}
	menu.Append(menuItemLevelFortress)

	currentStatusLock.Unlock()

	menu.Append(gtk.NewSeparatorMenuItem())

	menuItemStatus := gtk.NewMenuItem()
	menuItemStatusLabel := gtk.NewLabel(GetThreatMessage())
	menuItemStatusLabel.SetSizeRequest(250, -1)
	menuItemStatusLabel.SetLineWrap(true)
	menuItemStatus.Add(menuItemStatusLabel)
	menuItemStatus.SetSensitive(false)
	menu.Append(menuItemStatus)

	// create return to level button
	currentStatusLock.Lock()
	if currentStatus.CurrentSecurityLevel != currentStatus.SelectedSecurityLevel {
		menuItemReturn := gtk.NewMenuItemWithLabel(fmt.Sprintf("Ignore and return to Level %s", currentStatus.FmtSelectedLevel()))
		menuItemReturn.Connect("activate", func() {
			currentStatusLock.Lock()
			setToLevel := currentStatus.SelectedSecurityLevel
			currentStatusLock.Unlock()
			SetLevel(setToLevel)
		})
		menu.Append(menuItemReturn)
	}
	currentStatusLock.Unlock()

	menu.ShowAll()
	menu.Popup(nil, nil, gtk.StatusIconPositionMenu, systrayIcon, uint(cbx.Args(0)), uint32(cbx.Args(1)))
}
