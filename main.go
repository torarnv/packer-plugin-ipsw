// Copyright (c) Tor Arne Vestbø
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/hashicorp/packer-plugin-sdk/version"

	"github.com/torarnv/packer-plugin-ipsw/datasource"
)

var (
	Version           = "0.0.10"
	VersionPrerelease = ""
	PluginVersion     = version.InitializePluginVersion(Version, VersionPrerelease)
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterDatasource(plugin.DEFAULT_NAME, new(ipsw.Datasource))
	pps.SetVersion(PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
