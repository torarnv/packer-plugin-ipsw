# Packer Plugin IPSW

The `IPSW` data-source plugin can be used with HashiCorp [Packer](https://www.packer.io)
to fetch information about [Apple device firmware](https://en.wikipedia.org/wiki/IPSW),
for example for building macOS images with [Tart](https://github.com/cirruslabs/packer-plugin-tart).

## Installation

To install this plugin, copy and paste this code into your Packer configuration .
Then, run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    ipsw = {
      version = ">= 0.1.3"
      source  = "github.com/torarnv/ipsw"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/torarnv/ipsw
```

## Usage

To fetch information about the latest macOS release:

```hcl
data "ipsw" "macos" {
  os = "macOS"
  version = "^14"
  device = "VirtualMac2,1"
}
```

Then use the result via e.g. `${data.ipsw.macos.version}` or `${data.ipsw.macos.url}`.

## Configuration

For more information on how to configure the plugin, please read the
documentation located [here](https://developer.hashicorp.com/packer/plugins/datasources/ipsw).

## Contributing

* If you think you've found a bug in the code or you have a question regarding
  the usage of this software, please reach out by opening an issue in
  this GitHub repository.
* Contributions to this project are welcome: if you want to add a feature or a
  fix a bug, please do so by opening a Pull Request in this GitHub repository.
  In case of feature contribution, we kindly ask you to open an issue to
  discuss it beforehand.
