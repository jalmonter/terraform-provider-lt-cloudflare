package sdkv2provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/consts"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccCheckCloudflareQueueConsumerSource(resourceID, accountID string) string {
	return fmt.Sprintf(`
resource "cloudflare_queue" "%[1]s" {
	account_id = "%[2]s"
	name = "%[1]s"
}
resource "cloudflare_worker_script" "%[1]s" {
	account_id = "%[2]s"
	name = "%[1]s"
	content = "addEventListener('queue', (batch, event) => {});"
}
resource "cloudflare_queue_consumer" "%[1]s" {
	account_id 	= "%[2]s"
	queue_name 	= "%[1]s"
	script_name = cloudflare_worker_script.%[1]s.name

	settings {
		batch_size 		= 5
		max_retries 	= 1
		max_wait_time = 10
	}
}`, resourceID, accountID)
}

func TestAccCloudflareQueueConsumer_Basic(t *testing.T) {
	skipForDefaultAccount(t, "Pending investigation into automating the setup and teardown.")

	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")

	rnd := generateRandomResourceName()
	queueResourceName := "cloudflare_queue." + rnd
	workerResourceName := "cloudflare_worker_script." + rnd
	queueConsumerResourceName := "cloudflare_queue_consumer." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCloudflareQueueConsumerSource(rnd, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(queueResourceName, consts.AccountIDSchemaKey, accountID),
					resource.TestCheckResourceAttr(queueResourceName, "name", rnd),

					resource.TestCheckResourceAttr(workerResourceName, consts.AccountIDSchemaKey, accountID),
					resource.TestCheckResourceAttr(workerResourceName, "name", rnd),

					resource.TestCheckResourceAttr(queueConsumerResourceName, consts.AccountIDSchemaKey, accountID),
					resource.TestCheckResourceAttr(queueConsumerResourceName, "queue_name", rnd),
					resource.TestCheckResourceAttr(queueConsumerResourceName, "script_name", rnd),
					resource.TestCheckResourceAttr(queueConsumerResourceName, "settings.0.batch_size", "5"),
					resource.TestCheckResourceAttr(queueConsumerResourceName, "settings.0.max_retries", "1"),
					resource.TestCheckResourceAttr(queueConsumerResourceName, "settings.0.max_wait_time", "10"),
				),
			},
		},
	})
}
