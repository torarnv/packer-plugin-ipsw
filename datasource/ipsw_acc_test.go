// Copyright (c) Tor Arne Vestbø
// SPDX-License-Identifier: MPL-2.0

package ipsw

import (
    _ "embed"
    "fmt"
    "os/exec"
    "testing"

    "github.com/hashicorp/packer-plugin-sdk/acctest"
)

//go:embed test-fixtures/template.pkr.hcl
var testDatasourceHCL2Basic string

func TestAccIpswDatasource(t *testing.T) {
    testCase := &acctest.PluginTestCase{
        Name:     "ipsw_datasource_basic_test",
        Template: testDatasourceHCL2Basic,
        Check: func(buildCommand *exec.Cmd, logfile string) error {
            if buildCommand.ProcessState != nil {
                if buildCommand.ProcessState.ExitCode() != 0 {
                    return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
                }
            }
            return nil
        },
    }
    acctest.TestPlugin(t, testCase)
}
