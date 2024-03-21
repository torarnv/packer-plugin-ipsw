// Copyright (c) Tor Arne Vestbø
// SPDX-License-Identifier: MPL-2.0

package ipsw

import (
    "context"
    "encoding/json"
    "log"
    "os"
    "os/signal"
    "path/filepath"
    "strings"

    "github.com/go-git/go-git/v5"

    "golang.org/x/exp/slices"

    "github.com/Masterminds/semver/v3"
)

var (
    AppleDbGitUrl    = "https://github.com/littlebyteorg/appledb.git"
    AppleDbLocalPath = filepath.Join(os.TempDir(), "appledb")
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
    var appleDbPath = config.AppleDBLocalPath
    if appleDbPath == "" {
        appleDbPath = AppleDbLocalPath
    }

    if !config.Offline {
        var appleDbGitURL = config.AppleDBGitURL
        if appleDbGitURL == "" {
            appleDbGitURL = AppleDbGitUrl
        }

        log.Printf("Fetching latest AppleDB information from %s into %s",
            appleDbGitURL, appleDbPath)


        stop := make(chan os.Signal, 1)
        signal.Notify(stop, os.Interrupt)

        ctx, cancel := context.WithCancel(context.TODO())
        defer cancel()

        go func() {
            <-stop
            cancel()
        }()

        if _, err := os.Stat(appleDbPath); os.IsNotExist(err) {
            if _, err := git.PlainCloneContext(ctx, appleDbPath, false, &git.CloneOptions{
                URL:           appleDbGitURL,
                SingleBranch:  true,
                ReferenceName: "refs/heads/main",
                Progress:      log.Writer(),
            }); err != nil {
                return nil, err
            }
        } else {
            r, err := git.PlainOpen(appleDbPath)
            if err != nil {
                return nil, err
            }

            w, err := r.Worktree()
            if err != nil {
                return nil, err
            }

            if err = w.Pull(&git.PullOptions{
                Progress: log.Writer(),
            }); err != nil && err != git.NoErrAlreadyUpToDate {
                return nil, err
            }
        }
    }

    log.Printf("Ingesting AppleDb from %s", appleDbPath)
    if _, err := os.Stat(appleDbPath); os.IsNotExist(err) {
        return nil, err
    }

    var outputs []DatasourceOutput
    if err := filepath.Walk(filepath.Join(appleDbPath, "osFiles"),
        func(path string, f os.FileInfo, err error) error {
            if f.IsDir() {
                return nil
            }

            dat, err := os.ReadFile(path)
            if err != nil {
                return err
            }

            var osFile OsFile
            if err := json.Unmarshal(dat, &osFile); err != nil {
                log.Printf("Skipping %s (%s)", path, err)
                return nil
            }

            if osFile.OS != config.OS {
                return nil
            }
            if len(osFile.Sources) == 0 {
                return nil
            }

            version := strings.Replace(osFile.Version, " ", "-", 1)
            version = strings.Replace(version, " ", ".", -1)
            semVer, err := semver.NewVersion(version)
            if err != nil {
                log.Printf("Skipping un-parsable version %s", osFile.Version)
                return nil
            }

            if osFile.Beta && semVer.Prerelease() == "" {
                *semVer, _ = semVer.SetPrerelease("beta")
            }
            if !config.versionConstraints.Check(semVer) {
                return nil
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
                return nil
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

            return nil
        }); err != nil {
        return nil, err
    }

    return outputs, nil
}
