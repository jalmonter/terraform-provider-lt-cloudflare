---
page_title: "cloudflare_teams_list Resource - Cloudflare"
subcategory: ""
description: |-
  Provides a Cloudflare Teams List resource. Teams lists are
  referenced when creating secure web gateway policies or device
  posture rules.
---

# cloudflare_teams_list (Resource)

Provides a Cloudflare Teams List resource. Teams lists are
referenced when creating secure web gateway policies or device
posture rules.

## Example Usage

```terraform
resource "cloudflare_teams_list" "example" {
  account_id  = "f037e56e89293a057740de681ac9abbe"
  name        = "Corporate devices"
  type        = "SERIAL"
  description = "Serial numbers for all corporate devices."
  items       = ["8GE8721REF", "5RE8543EGG", "1YE2880LNP"]
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_id` (String) The account identifier to target for the resource.
- `name` (String) Name of the teams list.
- `type` (String) The teams list type. Available values: `IP`, `SERIAL`, `URL`, `DOMAIN`, `EMAIL`.

### Optional

- `description` (String) The description of the teams list.
- `items` (Set of String) The items of the teams list.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
$ terraform import cloudflare_teams_list.example <account_id>/<teams_list_id>
```
