package sdkv2provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cloudflare/cloudflare-go"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/consts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudflareQueueConsumer() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceCloudflareQueueConsumerSchema(),
		CreateContext: resourceCloudflareQueueConsumerUpdate,
		ReadContext:   resourceCloudflareQueueConsumerRead,
		UpdateContext: resourceCloudflareQueueConsumerUpdate,
		DeleteContext: resourceCloudflareQueueConsumerDelete,
		Description: heredoc.Doc(fmt.Sprintf(`
			A consumer Worker receives messages from your queue. When the consumer 
			Worker receives your queue’s messages, it can write them to another source, 
			such as a logging console or storage objects.
		`)),
	}
}

func resourceCloudflareQueueConsumerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)
	accountID := d.Get(consts.AccountIDSchemaKey).(string)
	scriptName := d.Get("script_name").(string)
	queueName := d.Get("queue_name").(string)

	settings := d.Get("settings").(map[string]interface{})

	_, err := client.UpdateQueueConsumer(ctx, cloudflare.AccountIdentifier(accountID), cloudflare.UpdateQueueConsumerParams{
		Consumer: cloudflare.QueueConsumer{
			QueueName:  queueName,
			ScriptName: scriptName,
			Settings: cloudflare.QueueConsumerSettings{
				BatchSize:   settings["batch_size"].(int),
				MaxWaitTime: settings["max_wait_time"].(int),
				MaxRetires:  settings["max_retries"].(int),
			},
			DeadLetterQueue: settings["dead_letter_queue"].(string),
		},
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update queue consumer: %w", err))
	}

	d.SetId(stringChecksum(scriptName))

	return nil
}

func resourceCloudflareQueueConsumerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)
	scriptName := d.Get("script_name").(string)
	accountID := d.Get(consts.AccountIDSchemaKey).(string)
	queueName := d.Get("queue_name").(string)

	params := cloudflare.ListQueueConsumersParams{
		QueueName: queueName,
	}

	queueConsumers, _, err := client.ListQueueConsumers(ctx, cloudflare.AccountIdentifier(accountID), params)
	if err != nil {
		var notFoundError *cloudflare.NotFoundError
		if errors.As(err, &notFoundError) {
			d.SetId("")
			return nil
		}

		return diag.FromErr(fmt.Errorf("failed to read queue consumers: %w", err))
	}

	for _, queueConsumer := range queueConsumers {
		if queueConsumer.ScriptName == scriptName {
			if settingsJson, err := json.Marshal(queueConsumer.Settings); err != nil {
				return diag.FromErr(fmt.Errorf("failed to marshal settings attribute: %w", err))
			} else if err := d.Set("settings", string(settingsJson)); err != nil {
				return diag.FromErr(fmt.Errorf("failed to set settings attribute: %w", err))
			}

			if err := d.Set("dead_letter_queue", queueConsumer.DeadLetterQueue); err != nil {
				return diag.FromErr(fmt.Errorf("failed to set dead letter queue attribute: %w", err))
			}
		}
	}

	return nil
}

func resourceCloudflareQueueConsumerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)
	scriptName := d.Get("script_name").(string)
	accountID := d.Get(consts.AccountIDSchemaKey).(string)
	queueName := d.Get("queue_name").(string)

	err := client.DeleteQueueConsumer(ctx, cloudflare.AccountIdentifier(accountID), cloudflare.DeleteQueueConsumerParams{
		QueueName:    queueName,
		ConsumerName: scriptName,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
