---
page_title: "cloudflare_waiting_room_settings Resource - Cloudflare"
subcategory: ""
description: |-
  Configure zone-wide settings for Cloudflare waiting rooms.
---

# cloudflare_waiting_room_settings (Resource)

Configure zone-wide settings for Cloudflare waiting rooms.

## Example Usage

```terraform
resource "cloudflare_waiting_room_settings" "example" {
  zone_id                      = "0da42c8d2132a9ddaf714f9e7c920711"
  search_engine_crawler_bypass = true
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `zone_id` (String) The zone identifier to target for the resource. **Modifying this attribute will force creation of a new resource.**

### Optional

- `search_engine_crawler_bypass` (Boolean) Whether to allow verified search engine crawlers to bypass all waiting rooms on this zone. Defaults to `false`.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
$ terraform import cloudflare_waiting_room_settings.example <zone_id>
```
