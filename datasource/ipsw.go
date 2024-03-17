// Copyright (c) Tor Arne Vestbø
// SPDX-License-Identifier: MPL-2.0

//go:generate -command packer-sdc go run github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc
//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput,VersionComponents
//go:generate packer-sdc struct-markdown

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
    // The operating system to filter on, e.g. `macOS`.
    OS                 string `mapstructure:"os" required:"true"`
    // The version to filter on. Semantic version conditions such
    // as `>= 12.2` or `~13.1` are supported. To include beta releases
    // in the search, use the `-0` prerelease constraint, e.g.
    // `^14-0` for the latest macOS 14 beta.
    Version            string `mapstructure:"version" required:"true"`
    // The device identifier to filter on, e.g. `VirtualMac2,1`.
    Device             string `mapstructure:"device" required:"true"`

    // The AppleDB Git URL to use for fetching release information.
    // Defaults to `https://github.com/littlebyteorg/appledb.git`.
    AppleDBGitURL      string `mapstructure:"appledb_git_url"`
    // The AppleDB local file path. Used as Git clone destination,
    // or for a pre-populated offline database.
    AppleDBLocalPath   string `mapstructure:"appledb_local_path"`
    // Set this property to `true` to skip fetching new
    // releases. Requires a valid AppleDB file structure in
    // `appledb_local_path`.
    Offline              bool `mapstructure:"offline"`

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
    } else {
        d.config.versionConstraints = *constraints
    }

    if errs != nil && len(errs.Errors) > 0 {
        return errs
    }
    return nil
}

type DatasourceOutput struct {
    // The operating system of the resulting release.
    OS                string `mapstructure:"os"`
    // The version of the resulting release, in full semantic version format.
    // To use individual components of the version, see `version_components`.
    Version           string `mapstructure:"version"`
    // The build identifier of the release. Also available as `metadata`
    // of the `version_components` field.
    Build             string `mapstructure:"build"`
    // The release date of the release.
    Released          string `mapstructure:"released"`
    // A boolean value reflecting whether the release is a beta release or not.
    Beta              bool   `mapstructure:"beta"`
    // The URL of the IPSW file for the release.
    URL               string `mapstructure:"url"`
    // Individual components of the `version` field.
    VersionComponents *VersionComponents `mapstructure:"version_components"`

    semVer            *semver.Version
}

type VersionComponents struct {
    // The major version of the release.
    Major      uint64 `mapstructure:"major"`
    // The minor version of the release.
    Minor      uint64 `mapstructure:"minor"`
    // The patch version of the release.
    Patch      uint64 `mapstructure:"patch"`
    // The prerelease tag of the release, e.g. `beta.2`.
    Prerelease string `mapstructure:"prerelease"`
    // The metadata of the release, e.g. the build identifier.
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
