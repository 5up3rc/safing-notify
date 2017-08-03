// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the GPL license that can be found in the LICENSE file.

package meta

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var (
	version      = "0.0.0"
	commit       string
	buildOptions string
	buildUser    string
	buildHost    string
	buildDate    string
	buildSource  string
)

func init() {
	if !strings.HasSuffix(os.Args[0], ".test") {
		if commit == "" ||
			buildUser == "" ||
			buildHost == "" ||
			buildDate == "" ||
			buildSource == "" {
			fmt.Fprintln(os.Stderr, "FATAL ERROR: please build using the supplied build script.\n$ ./build")
			os.Exit(1)
		}
	}
}

func Version() string {
	if strings.HasPrefix(commit, fmt.Sprintf("v%s-0-", version)) {
		return version
	} else {
		return version + "*"
	}
}

func FullVersion() string {
	s := ""
	if strings.HasPrefix(commit, fmt.Sprintf("v%s-0-", version)) {
		s += fmt.Sprintf("Safing Notify\nversion %s\n", version)
	} else {
		s += fmt.Sprintf("Safing Notify\ndevelopment build, built on top version %s\n", version)
	}
	s += fmt.Sprintf("\ncommit %s\n", commit)
	s += fmt.Sprintf("built with %s (%s) %s/%s\n", runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
	s += fmt.Sprintf("  using options %s\n", strings.Replace(buildOptions, "§", " ", -1))
	s += fmt.Sprintf("  by %s@%s\n", buildUser, buildHost)
	s += fmt.Sprintf("  on %s\n", buildDate)
	s += fmt.Sprintf("\nLicensed under the GPL license.\nThe source code is available here: %s", buildSource)
	return s
}
