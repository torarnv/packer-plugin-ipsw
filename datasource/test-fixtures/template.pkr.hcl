# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

data "ipsw" "ventura" {
  os = "macOS"
  version = "13.4"
  device = "VirtualMac2,1"
}

data "ipsw" "ventura-beta" {
  os = "macOS"
  version = "13.5.0-beta.5"
  device = "VirtualMac2,1"
}

source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = [
    "source.null.basic-example"
  ]

  provisioner "shell-local" {
    inline = [
      "echo '${data.ipsw.ventura.os} ${data.ipsw.ventura.version}'",
      "echo '${data.ipsw.ventura.version_components.major}' | grep 13",
      "echo '${data.ipsw.ventura.version_components.minor}' | grep 4",
      "echo '${data.ipsw.ventura.version_components.patch}' | grep 1",
      "echo '${data.ipsw.ventura.version_components.prerelease}' | grep '^$'",
      "echo '${data.ipsw.ventura.version_components.metadata}' | grep 22F82",
      "echo '${data.ipsw.ventura.released}' | grep 2023-06-21",
      "echo '${data.ipsw.ventura.beta}' | grep false",
      "echo '${data.ipsw.ventura.url}' | grep 22F82",
    ]
  }

  provisioner "shell-local" {
    inline = [
      "echo '${data.ipsw.ventura-beta.os} ${data.ipsw.ventura-beta.version}'",
      "echo '${data.ipsw.ventura-beta.version_components.major}' | grep 13",
      "echo '${data.ipsw.ventura-beta.version_components.minor}' | grep 5",
      "echo '${data.ipsw.ventura-beta.version_components.patch}' | grep 0",
      "echo '${data.ipsw.ventura-beta.version_components.prerelease}' | grep beta.5",
      "echo '${data.ipsw.ventura-beta.version_components.metadata}' | grep 22G5072a",
      "echo '${data.ipsw.ventura-beta.released}' | grep 2023-07-10",
      "echo '${data.ipsw.ventura-beta.beta}' | grep true",
      "echo '${data.ipsw.ventura-beta.url}' | grep 22G5072a",
    ]
  }
}
