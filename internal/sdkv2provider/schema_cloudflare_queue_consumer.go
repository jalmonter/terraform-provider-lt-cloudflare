package sdkv2provider

import (
	"github.com/cloudflare/terraform-provider-cloudflare/internal/consts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var queueConsumerSettingsResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"batch_size": {
			Type:        schema.TypeInt,
			Description: "Number of messages in the queue to be batched for delivery. Specify a value from 1-100.",
			Optional:    true,
		},
		"max_wait_time": {
			Type:        schema.TypeInt,
			Description: "Maximum time the queue will wait before it sends messages in a batch. Specify a value from 0-60.",
			Optional:    true,
		},
		"max_retries": {
			Type:        schema.TypeInt,
			Description: "Maximum number of times to retry delivery of a message before failing.",
			Optional:    true,
		},
		"dead_letter_queue": {
			Type:        schema.TypeString,
			Description: "Failed messages will redirect to this queue and then can then be viewed, deleted, or resent.",
			Optional:    true,
		},
	},
}

func resourceCloudflareQueueConsumerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		consts.AccountIDSchemaKey: {
			Description: consts.AccountIDSchemaDescription,
			Type:        schema.TypeString,
			Required:    true,
		},
		"queue_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the queue you want to use.",
		},
		"script_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the worker script.",
		},
		"settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     queueConsumerSettingsResource,
		},
	}
}
