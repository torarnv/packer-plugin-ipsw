# Copyright (c) Tor Arne Vestbø
# SPDX-License-Identifier: MPL-2.0

integration {
  name = "IPSW"
  description = "IPSW plugins for HashiCorp Packer"
  identifier = "packer/torarnv/ipsw"
  docs {
    process_docs = true
  }
  license {
    type = "MPL-2.0"
    url = "https://github.com/hashicorp/integration-template/blob/main/LICENSE.md"
  }
  component {
    type = "data-source"
    name = "IPSW"
    slug = "ipsw"
  }
}
