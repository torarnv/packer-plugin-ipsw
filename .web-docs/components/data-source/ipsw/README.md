Type: `ipsw`

The IPSW data source will fetch and filter information about
[Apple device firmware](https://en.wikipedia.org/wiki/IPSW)
releases from [AppleDB](https://appledb.dev/), providing release
data and IPSW URLs for builders such as [Tart](https://github.com/cirruslabs/packer-plugin-tart).

## Configuration

### Required

<!-- Code generated from the comments of the Config struct in datasource/ipsw.go; DO NOT EDIT MANUALLY -->

- `os` (string) - The operating system to filter on, e.g. `macOS`.

- `version` (string) - The version to filter on. Semantic version conditions such
  as `>= 12.2` or `~13.1` are supported. To include beta releases
  in the search, use the `-0` prerelease constraint, e.g.
  `^14-0` for the latest macOS 14 beta.

- `device` (string) - The device identifier to filter on, e.g. `VirtualMac2,1`.

<!-- End of code generated from the comments of the Config struct in datasource/ipsw.go; -->


## Output

<!-- Code generated from the comments of the DatasourceOutput struct in datasource/ipsw.go; DO NOT EDIT MANUALLY -->

- `os` (string) - The operating system of the resulting release.

- `version` (string) - The version of the resulting release, in full semantic version format.
  To use individual components of the version, see `version_components`.

- `build` (string) - The build identifier of the release. Also available as `metadata`
  of the `version_components` field.

- `released` (string) - The release date of the release.

- `beta` (bool) - A boolean value reflecting whether the release is a beta release or not.

- `url` (string) - The URL of the IPSW file for the release.

- `version_components` (\*VersionComponents) - Individual components of the `version` field.

<!-- End of code generated from the comments of the DatasourceOutput struct in datasource/ipsw.go; -->


### Version components

<!-- Code generated from the comments of the VersionComponents struct in datasource/ipsw.go; DO NOT EDIT MANUALLY -->

- `major` (uint64) - The major version of the release.

- `minor` (uint64) - The minor version of the release.

- `patch` (uint64) - The patch version of the release.

- `prerelease` (string) - The prerelease tag of the release, e.g. `beta.2`.

- `metadata` (string) - The metadata of the release, e.g. the build identifier.

<!-- End of code generated from the comments of the VersionComponents struct in datasource/ipsw.go; -->
