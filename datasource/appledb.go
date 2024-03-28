// Copyright (c) Tor Arne Vestbø
// SPDX-License-Identifier: MPL-2.0

package ipsw

import (
    "log"
    "strings"

    "net/http"
    "net/url"
    "io/ioutil"
    "compress/gzip"

    "encoding/json"

    "golang.org/x/exp/slices"

    "github.com/Masterminds/semver/v3"
)

var (
    AppleDbBaseUrl = "https://api.appledb.dev/ios"
)

type OsFileSource struct {
    Type      string   `json:"type"`
    DeviceMap []string `json:"deviceMap"`
    Links     []struct {
        URL       string `json:"url"`
        Preferred bool   `json:"preferred"`
        Active    bool   `json:"active"`
    } `json:"links"`
    Hashes struct {
        Sha2256 string `json:"sha2-256"`
        Sha1    string `json:"sha1"`
    } `json:"hashes"`
    Size int64 `json:"size"`
}

type OsFile struct {
    OS        string         `json:"osStr"`
    Version   string         `json:"version"`
    Build     string         `json:"build"`
    Released  string         `json:"released"`
    Beta      bool           `json:"beta"`
    DeviceMap []string       `json:"deviceMap"`
    Sources   []OsFileSource `json:"sources"`
}

func QueryAppleDB(config Config) (DatasourceOutputs, error) {

    url, err := url.JoinPath(AppleDbBaseUrl, config.OS, "main.json.gz")
    if err != nil {
        return nil, err
    }
    log.Println("Fetching releases from " + url)
    response, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    reader, err := gzip.NewReader(response.Body)
    if err != nil {
        return nil, err
    }
    defer reader.Close()

    rawJson, err := ioutil.ReadAll(reader)
    if err != nil {
        return nil, err
    }

    var osFiles []json.RawMessage
    json.Unmarshal(rawJson, &osFiles)

    var outputs []DatasourceOutput
    for _, osFileRawJson := range osFiles {
        var osFile OsFile
        if err := json.Unmarshal(osFileRawJson, &osFile); err != nil {
            osFileJson, _ := json.Marshal(&osFileRawJson)
            log.Printf("Skipping un-parsable JSON '%s' due to %s",
                osFileJson, err)
            continue
        }
        if osFile.OS != config.OS {
            continue
        }
        if len(osFile.Sources) == 0 {
            continue
        }

        version := strings.Replace(osFile.Version, " ", "-", 1)
        version = strings.Replace(version, " ", ".", -1)
        semVer, err := semver.NewVersion(version)
        if err != nil {
            log.Printf("Skipping un-parsable version '%s'", osFile.Version)
            continue
        }

        if osFile.Beta && semVer.Prerelease() == "" {
            *semVer, _ = semVer.SetPrerelease("beta")
        }
        if !config.versionConstraints.Check(semVer) {
            continue
        }
        if semVer.Metadata() == "" {
            *semVer, _ = semVer.SetMetadata(osFile.Build)
        }

        var url string
        for _, source := range osFile.Sources {
            if source.Type != "ipsw" {
                continue
            }
            if config.Device != "" {
                if !slices.Contains(source.DeviceMap, config.Device) {
                    continue
                }
            }
            for _, link := range source.Links {
                if link.Active && (url == "" || link.Preferred) {
                    url = link.URL
                }
            }
        }

        if url == "" {
            continue
        }

        outputs = append(outputs, DatasourceOutput{
            OS:       osFile.OS,
            Version:  osFile.Version,
            Build:    osFile.Build,
            Released: osFile.Released,
            Beta:     osFile.Beta,
            URL:      url,
            semVer:   semVer,
        })

    }

    return outputs, nil
}
