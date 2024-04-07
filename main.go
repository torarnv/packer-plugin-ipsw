// Copyright (c) Tor Arne Vestbø
// SPDX-License-Identifier: MPL-2.0

//go:generate go run git.rootprojects.org/root/go-gitver/v2@v2.0.2 --outfile version.go

package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	packerVersion "github.com/hashicorp/packer-plugin-sdk/version"

	"github.com/torarnv/packer-plugin-ipsw/datasource"
)

var (
	// Generated at build time
	commit  = "0000000"
	version = "0.0.0-pre0+0000000"
	date    = "0000-00-00T00:00:00+0000"
)

func main() {
	pps := plugin.NewSet()

	regex := regexp.MustCompile("(.*)?-pre.*$")
	version = regex.ReplaceAllString(version, "$1-dev")
	pps.SetVersion(packerVersion.InitializePluginVersion(version, ""))

	pps.RegisterDatasource(plugin.DEFAULT_NAME, new(ipsw.Datasource))

	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
