# Copyright (c) Tor Arne Vestbø
# SPDX-License-Identifier: MPL-2.0

integration {
  name = "IPSW"
  description = "The IPSW data source provides information about Apple IPSW releases based on the filter options provided in the configuration."
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
