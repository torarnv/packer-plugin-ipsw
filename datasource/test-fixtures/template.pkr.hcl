# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "appledb_test_path" {
  type = string
}

data "ipsw" "bigsur" {
  os = "macOS"
  version = "^13"
  device = "VirtualMac2,1"

  offline = true
  appledb_local_path = var.appledb_test_path
}

data "ipsw" "bigsur-beta" {
  os = "macOS"
  version = "^13-0"
  device = "VirtualMac2,1"

  offline = true
  appledb_local_path = var.appledb_test_path
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
      "echo '${data.ipsw.bigsur.os} ${data.ipsw.bigsur.version}'",
      "echo '${data.ipsw.bigsur.version_components.major}' | grep 13",
      "echo '${data.ipsw.bigsur.version_components.minor}' | grep 4",
      "echo '${data.ipsw.bigsur.version_components.patch}' | grep 1",
      "echo '${data.ipsw.bigsur.version_components.prerelease}' | grep '^$'",
      "echo '${data.ipsw.bigsur.version_components.metadata}' | grep 22F82",
      "echo '${data.ipsw.bigsur.released}' | grep 2023-06-21",
      "echo '${data.ipsw.bigsur.beta}' | grep false",
      "echo '${data.ipsw.bigsur.url}' | grep 22F82",
    ]
  }

  provisioner "shell-local" {
    inline = [
      "echo '${data.ipsw.bigsur-beta.os} ${data.ipsw.bigsur-beta.version}'",
      "echo '${data.ipsw.bigsur-beta.version_components.major}' | grep 13",
      "echo '${data.ipsw.bigsur-beta.version_components.minor}' | grep 5",
      "echo '${data.ipsw.bigsur-beta.version_components.patch}' | grep 0",
      "echo '${data.ipsw.bigsur-beta.version_components.prerelease}' | grep beta.5",
      "echo '${data.ipsw.bigsur-beta.version_components.metadata}' | grep 22G5072a",
      "echo '${data.ipsw.bigsur-beta.released}' | grep 2023-07-10",
      "echo '${data.ipsw.bigsur-beta.beta}' | grep true",
      "echo '${data.ipsw.bigsur-beta.url}' | grep 22G5072a",
    ]
  }
}
