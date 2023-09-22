package sdkv2provider

import (
	"context"
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
	deadLetterQueue := d.Get("dead_letter_queue").(string)

	settings := cloudflare.QueueConsumerSettings{}
	if batchSize, ok := d.GetOk("settings.batch_size"); ok {
		settings.BatchSize = batchSize.(int)
	}
	if maxWaitTime, ok := d.GetOk("settings.max_wait_time"); ok {
		settings.MaxWaitTime = maxWaitTime.(int)
	}
	if maxRetries, ok := d.GetOk("settings.max_retries"); ok {
		settings.MaxRetires = maxRetries.(int)
	}

	_, err := client.UpdateQueueConsumer(ctx, cloudflare.AccountIdentifier(accountID), cloudflare.UpdateQueueConsumerParams{
		QueueName: queueName,
		Consumer: cloudflare.QueueConsumer{
			Name:            scriptName,
			ScriptName:      scriptName,
			QueueName:       queueName,
			Settings:        settings,
			DeadLetterQueue: deadLetterQueue,
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
	accountID := d.Get(consts.AccountIDSchemaKey).(string)
	scriptName := d.Get("script_name").(string)
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
			if err := d.Set("settings.batch_size", queueConsumer.Settings.BatchSize); err != nil {
				return diag.FromErr(fmt.Errorf("failed to set batch size queue settings attribute: %w", err))
			}
			if err := d.Set("settings.max_retries", queueConsumer.Settings.MaxRetires); err != nil {
				return diag.FromErr(fmt.Errorf("failed to set max retries queue settings attribute: %w", err))
			}
			if err := d.Set("settings.max_wait_time", queueConsumer.Settings.MaxWaitTime); err != nil {
				return diag.FromErr(fmt.Errorf("failed to set max wait time queue settings attribute: %w", err))
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
	accountID := d.Get(consts.AccountIDSchemaKey).(string)
	scriptName := d.Get("script_name").(string)
	queueName := d.Get("queue_name").(string)

	err := client.DeleteQueueConsumer(ctx, cloudflare.AccountIdentifier(accountID), cloudflare.DeleteQueueConsumerParams{
		QueueName:    queueName,
		ConsumerName: scriptName,
	})
	if err != nil {
		// If the resource is already deleted, we should not return without an error
		// according to the terraform spec
		var notFoundError *cloudflare.NotFoundError
		if errors.As(err, &notFoundError) {
			return nil
		}

		return diag.FromErr(fmt.Errorf("failed to delete queue consumers: %w", err))
	}

	d.SetId("")

	return nil
}
