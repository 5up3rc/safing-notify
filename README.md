[![Travis](https://img.shields.io/travis/Safing/safing-notify.svg?style=flat-square)](https://travis-ci.org/Safing/safing-notify)
[![Coveralls](https://img.shields.io/coveralls/Safing/safing-notify.svg?branch=master&style=flat-square)](https://coveralls.io/github/Safing/safing-notify?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Safing/safing-notify?style=flat-square)](https://goreportcard.com/report/github.com/Safing/safing-notify)

# Safing Notify

Safing Notify is a small utility that sits in the system tray, shows the current status of Safing, notifies the user about things going on, and can be used to switch between modes quickly.

For more information about Safing, please check out [Safing Core](https://github.com/Safing/safing-core).

## Download

We recommend to download a packaged version of all components [here](https://github.com/Safing/safing-installer/releases).  
You can also just [download Safing Notify](https://github.com/Safing/safing-notify/releases).

## Developing

Currently Safing is only supported on Linux.

1. Install Go 1.8+ ([https://golang.org/dl/](https://golang.org/dl/))
2. Check additional requirements (see below)
3. Build:

        ./build

## Additional Requirements

#### Linux

- go dependencies
  - Some dependencies cannot be vendored, because the build tool gets wound up:

        go get github.com/mattn/go-gtk/glib
        go get github.com/mattn/go-gtk/gtk

- gtk (FIXME)
