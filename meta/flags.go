// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the GPL license that can be found in the LICENSE file.

package meta

import (
	"flag"
	"fmt"
	"os"
)

var (
	showVersion = flag.Bool("v", false, "show version and exit")
)

func init() {
	flag.Parse()

	if *showVersion {
		fmt.Println(FullVersion())
		os.Exit(0)
	}
}
