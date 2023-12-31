---
page_title: "cloudflare_logpush_ownership_challenge Resource - Cloudflare"
subcategory: ""
description: |-
  Provides a resource which manages Cloudflare Logpush ownership
  challenges to use in a Logpush Job. On it's own, doesn't do much
  however this resource should be used in conjunction to create
  Logpush jobs.
---

# cloudflare_logpush_ownership_challenge (Resource)

Provides a resource which manages Cloudflare Logpush ownership
challenges to use in a Logpush Job. On it's own, doesn't do much
however this resource should be used in conjunction to create
Logpush jobs.

## Example Usage

```terraform
resource "cloudflare_logpush_ownership_challenge" "example" {
  zone_id          = "0da42c8d2132a9ddaf714f9e7c920711"
  destination_conf = "s3://my-bucket-path?region=us-west-2"
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `destination_conf` (String) Uniquely identifies a resource (such as an s3 bucket) where data will be pushed. Additional configuration parameters supported by the destination may be included. See [Logpush destination documentation](https://developers.cloudflare.com/logs/logpush/logpush-configuration-api/understanding-logpush-api/#destination). **Modifying this attribute will force creation of a new resource.**

### Optional

- `account_id` (String) The account identifier to target for the resource. Must provide only one of `account_id`, `zone_id`.
- `zone_id` (String) The zone identifier to target for the resource. Must provide only one of `account_id`, `zone_id`.

### Read-Only

- `id` (String) The ID of this resource.
- `ownership_challenge_filename` (String) The filename of the ownership challenge which	contains the contents required for Logpush Job creation.


