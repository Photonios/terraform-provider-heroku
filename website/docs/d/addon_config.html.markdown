---
layout: "heroku"
page_title: "Heroku: heroku_addon_config"
sidebar_current: "docs-heroku-datasource-addon-config-x"
description: |-
  Get the configuration for a Heroku Addon.
---

# Data Source: heroku_addon_config

Use this data source to gte the configuration for a Heroku Addon.

## Example Usage

```hcl
data "heroku_addon_config" "from_another_app" {
  name = "addon-from-another-app"
}

output "heroku_addon_data_basic" {
  value = [
    "config: ${transpose(data.heroku_addon_config.from_another_app.config)}",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The add-on name

## Attributes Reference

The following attributes are exported:

* `config` - The add-on configuration.
