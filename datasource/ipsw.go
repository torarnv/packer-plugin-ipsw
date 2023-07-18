// Copyright (c) Tor Arne Vestbø
// SPDX-License-Identifier: MPL-2.0

//go:generate -command packer-sdc go run github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc
//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput,VersionComponents

package ipsw

import (
    "fmt"
    "sort"
    "time"

    "github.com/hashicorp/hcl/v2/hcldec"
    "github.com/hashicorp/packer-plugin-sdk/hcl2helper"
    "github.com/hashicorp/packer-plugin-sdk/packer"
    "github.com/hashicorp/packer-plugin-sdk/template/config"
    "github.com/zclconf/go-cty/cty"

    "github.com/Masterminds/semver/v3"
)

type Datasource struct {
    config Config
}

type Config struct {
    OS                 string `mapstructure:"os" required:"true"`
    Version            string `mapstructure:"version" required:"true"`
    Device             string `mapstructure:"device" required:"true"`

    Offline              bool `mapstructure:"offline"`
    AppleDBGitURL      string `mapstructure:"appledb_git_url"`
    AppleDBLocalPath   string `mapstructure:"appledb_local_path"`

    versionConstraints semver.Constraints
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
    return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
    err := config.Decode(&d.config, nil, raws...)
    if err != nil {
        return err
    }

    var errs *packer.MultiError

    if d.config.OS == "" {
        errs = packer.MultiErrorAppend(errs, fmt.Errorf("an 'os' must be provided"))
    }

    if d.config.Version == "" {
        errs = packer.MultiErrorAppend(errs, fmt.Errorf("a 'version' must be provided"))
    }

    if d.config.Device == "" {
        errs = packer.MultiErrorAppend(errs, fmt.Errorf("a 'device' must be provided"))
    }

    constraints, err := semver.NewConstraint(d.config.Version)
    if err != nil {
        errs = packer.MultiErrorAppend(errs, fmt.Errorf("Could not parse version constraint '%s'",
            d.config.Version))
    }
    d.config.versionConstraints = *constraints

    if errs != nil && len(errs.Errors) > 0 {
        return errs
    }
    return nil
}

type DatasourceOutput struct {
    OS                string `mapstructure:"os"`
    Version           string `mapstructure:"version"`
    Build             string `mapstructure:"build"`
    Released          string `mapstructure:"released"`
    Beta              bool   `mapstructure:"beta"`
    URL               string `mapstructure:"url"`
    VersionComponents *VersionComponents `mapstructure:"version_components"`

    semVer            *semver.Version
}

type VersionComponents struct {
    Major      uint64 `mapstructure:"major"`
    Minor      uint64 `mapstructure:"minor"`
    Patch      uint64 `mapstructure:"patch"`
    Prerelease string `mapstructure:"prerelease"`
    Metadata   string `mapstructure:"metadata"`
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
    return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
    results, err := QueryAppleDB(d.config)
    if err != nil {
        return cty.NullVal(cty.EmptyObject), err
    }

    if len(results) == 0 {
        return cty.NullVal(cty.EmptyObject), fmt.Errorf("No matching IPSWs for the given filters")
    }

    sort.Sort(results)
    var mostRecent = results[len(results)-1]

    mostRecent.Version = mostRecent.semVer.String()
    mostRecent.VersionComponents = &VersionComponents{
        Major:      mostRecent.semVer.Major(),
        Minor:      mostRecent.semVer.Minor(),
        Patch:      mostRecent.semVer.Patch(),
        Prerelease: mostRecent.semVer.Prerelease(),
        Metadata:   mostRecent.semVer.Metadata(),
    }

    return hcl2helper.HCL2ValueFromConfig(mostRecent, d.OutputSpec()), nil
}


type DatasourceOutputs []DatasourceOutput

func (o DatasourceOutputs) Len() int {
    return len(o)
}

func (o DatasourceOutputs) Swap(i, j int) {
    o[i], o[j] = o[j], o[i]
}

func (o DatasourceOutputs) Less(i, j int) bool {
    v1 := o[i].semVer
    v2 := o[j].semVer
    if v1.Equal(v2) {
        const dateFormat = "2006-01-02"
        d1, _ := time.Parse(dateFormat, o[i].Released)
        d2, _ := time.Parse(dateFormat, o[j].Released)
        return d1.Before(d2)
    } else {
        return v1.LessThan(v2)
    }
}
