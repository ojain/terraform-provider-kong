package kong

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/kong/go-kong/kong"
)

func TestAccKongGlobalPlugin(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateGlobalPluginConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.hmac_auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "name", "hmac-auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "enabled", "true"),
				),
			},
			{
				Config: testUpdateGlobalPluginConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.hmac_auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "name", "hmac-auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccKongGlobalPluginDisabled(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateGlobalPluginConfigDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.hmac_auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "name", "hmac-auth"),
					resource.TestCheckResourceAttr("kong_plugin.hmac_auth", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccKongPluginForASpecificConsumer(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreatePluginForASpecificConsumerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.rate_limit"),
					testAccCheckKongConsumerExists("kong_consumer.plugin_consumer"),
					testAccCheckForChildIDCorrect("kong_consumer.plugin_consumer", "kong_plugin.rate_limit", "consumer_id"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "name", "rate-limiting"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "enabled", "true"),
				),
			},
			{
				Config: testUpdatePluginForASpecificConsumerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.rate_limit"),
					testAccCheckKongConsumerExists("kong_consumer.plugin_consumer"),
					testAccCheckForChildIDCorrect("kong_consumer.plugin_consumer", "kong_plugin.rate_limit", "consumer_id"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "name", "rate-limiting"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccKongPluginForASpecificServiceID(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreatePluginForASpecificServiceIdConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.rate_limit"),
					testAccCheckKongServiceExists("kong_service.service"),
					testAccCheckForChildIDCorrect("kong_service.service", "kong_plugin.rate_limit", "service_id"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "name", "rate-limiting"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "tags.#", "2"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "tags.0", "foo"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "tags.1", "bar"),
				),
			},
			{
				Config: testUpdatePluginForASpecificServiceIdConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.rate_limit"),
					testAccCheckKongServiceExists("kong_service.service"),
					testAccCheckForChildIDCorrect("kong_service.service", "kong_plugin.rate_limit", "service_id"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "name", "rate-limiting"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "tags.#", "1"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "tags.0", "foo"),
				),
			},
		},
	})
}

func TestAccKongPluginForASpecificServiceName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config:             testCreatePluginForASpecificServiceNameConfig,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.rate_limit"),
					testAccCheckKongServiceExists("kong_service.service"),
					testAccCheckForChildIDCorrect("kong_service.service", "kong_plugin.rate_limit", "service_id"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "name", "rate-limiting"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "tags.#", "2"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "tags.0", "foo"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "tags.1", "bar"),
				),
			},
		},
	})
}

func TestAccKongPluginForASpecificRoute(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreatePluginForASpecificRouteConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.rate_limit"),
					testAccCheckKongServiceExists("kong_service.service"),
					testAccCheckKongRouteExists("kong_route.route"),
					testAccCheckForChildIDCorrect("kong_route.route", "kong_plugin.rate_limit", "route_id"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "name", "rate-limiting"),
				),
			},
			{
				Config: testUpdatePluginForASpecificRouteConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKongPluginExists("kong_plugin.rate_limit"),
					testAccCheckKongServiceExists("kong_service.service"),
					testAccCheckKongRouteExists("kong_route.route"),
					testAccCheckForChildIDCorrect("kong_route.route", "kong_plugin.rate_limit", "route_id"),
					resource.TestCheckResourceAttr("kong_plugin.rate_limit", "name", "rate-limiting"),
				),
			},
		},
	})
}

func TestAccKongPluginImportConfigJson(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKongPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateGlobalPluginConfig,
			},
			{
				ResourceName:      "kong_plugin.hmac_auth",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccCheckKongPluginDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(*config).adminClient.Plugins

	plugins := getResourcesByType("kong_plugin", state)

	if len(plugins) != 1 {
		return fmt.Errorf("expecting only 1 plugin resource found %v", len(plugins))
	}

	response, err := client.Get(context.Background(), kong.String(plugins[0].Primary.ID))

	if !kong.IsNotFoundErr(err) && err != nil {
		return fmt.Errorf("error calling get plugin by id: %v", err)
	}

	if response != nil {
		return fmt.Errorf("plugin %s still exists, %+v", plugins[0].Primary.ID, response)
	}

	return nil
}

func testAccCheckForChildIDCorrect(parentResource string, childResource string, childIDField string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[parentResource]

		if !ok {
			return fmt.Errorf("not found: %s", parentResource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		parentID := rs.Primary.ID

		rs, ok = s.RootModule().Resources[childResource]

		if !ok {
			return fmt.Errorf("not found: %s", parentResource)
		}

		childID, ok := rs.Primary.Attributes[childIDField]

		if !ok {
			return fmt.Errorf("child id field %s not set on %s", childIDField, childResource)
		}

		if parentID != childID {
			return fmt.Errorf("expected %s id of %s to equal %s id field %s of %s", parentResource, parentID, childResource, childIDField, childID)
		}

		return nil
	}
}

func testAccCheckKongPluginExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*config).adminClient.Plugins
		api, err := client.Get(context.Background(), kong.String(rs.Primary.ID))

		if !kong.IsNotFoundErr(err) && err != nil {
			return err
		}

		if api == nil {
			return fmt.Errorf("plugin with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateGlobalPluginConfig = `
resource "kong_plugin" "hmac_auth" {
	name  = "hmac-auth"
	enabled = "true"
	config_json = <<EOT
	{
    	"algorithms": [
    	    "hmac-sha1",
    	    "hmac-sha256",
    	    "hmac-sha384",
    	    "hmac-sha512"
    	],
    	"clock_skew": 300,
    	"enforce_headers": [],
    	"hide_credentials": true,
    	"validate_request_body": false
	}
EOT
}
`

const testCreateGlobalPluginConfigDisabled = `
resource "kong_plugin" "hmac_auth" {
	name  = "hmac-auth"
	enabled = "false"
	config_json = <<EOT
	{
    	"algorithms": [
    	    "hmac-sha1",
    	    "hmac-sha256",
    	    "hmac-sha384",
    	    "hmac-sha512"
    	],
    	"clock_skew": 300,
    	"enforce_headers": [],
    	"hide_credentials": true,
    	"validate_request_body": false
	}
EOT
}
`

const testUpdateGlobalPluginConfig = `
resource "kong_plugin" "hmac_auth" {
	name  = "hmac-auth"
	config_json = <<EOT
	{
    	"algorithms": [
    	    "hmac-sha1",
    	    "hmac-sha256",
    	    "hmac-sha384",
    	    "hmac-sha512"
    	],
    	"clock_skew": 300,
    	"enforce_headers": [],
    	"hide_credentials": false,
    	"validate_request_body": false
	}
EOT
}
`

const testCreatePluginForASpecificConsumerConfig = `
resource "kong_consumer" "plugin_consumer" {
	username  = "PluginUser"
	custom_id = "567"
}

resource "kong_plugin" "rate_limit" {
	name        = "rate-limiting"
	consumer_id = "${kong_consumer.plugin_consumer.id}"
	config_json = <<EOT
	{
		"second": 5,
		"hour" : 1000
	}
EOT
}
`

const testUpdatePluginForASpecificConsumerConfig = `
resource "kong_consumer" "plugin_consumer" {
	username  = "PluginUser"
	custom_id = "567"
}

resource "kong_plugin" "rate_limit" {
	name        = "rate-limiting"
	consumer_id = "${kong_consumer.plugin_consumer.id}"
	config_json = <<EOT
	{
		"second": 10,
		"hour" : 2000
	}
EOT
}
`

const testCreatePluginForASpecificServiceIdConfig = `
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test.org"
}

resource "kong_plugin" "rate_limit" {
	name        = "rate-limiting"
	service_id = "${kong_service.service.id}"
	tags       = ["foo", "bar"]
	config_json = <<EOT
	{
		"second": 10,
		"hour" : 2000
	}
	
EOT
}
`

const testUpdatePluginForASpecificServiceIdConfig = `
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test.org"
}

resource "kong_plugin" "rate_limit" {
	name        = "rate-limiting"
	service_id = "${kong_service.service.id}"
	tags       = ["foo"]
	config_json = <<EOT
	{
		"second": 11,
		"hour" : 4000
	}
EOT
}`

const testCreatePluginForASpecificServiceNameConfig = `
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test2.org"
}

resource "kong_plugin" "rate_limit" {
	name        = "rate-limiting"
	service_name = "${kong_service.service.name}"
	tags       = ["foo", "bar"]
	config_json = <<EOT
	{
		"second": 10,
		"hour" : 2000
	}
	
EOT
}`

const testCreatePluginForASpecificRouteConfig = `
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test.org"
}

resource "kong_route" "route" {
	protocols 		= [ "http", "https" ]
	methods 		= [ "GET", "POST" ]
	hosts 			= [ "example2.com" ]
	paths 			= [ "/test" ]
	strip_path 		= false
	preserve_host 	= true
	service_id 		= "${kong_service.service.id}"
}

resource "kong_plugin" "rate_limit" {
	name        = "rate-limiting"
	route_id = "${kong_route.route.id}"
	config_json = <<EOT
	{
		"second": 12,
		"hour" : 3000
	}
EOT
}
`

const testUpdatePluginForASpecificRouteConfig = `
resource "kong_service" "service" {
	name     = "test"
	protocol = "http"
	host     = "test.org"
}

resource "kong_route" "route" {
	protocols 		= [ "http", "https" ]
	methods 		= [ "GET", "POST" ]
	hosts 			= [ "example2.com" ]
	paths 			= [ "/test" ]
	strip_path 		= false
	preserve_host 	= true
	service_id 		= "${kong_service.service.id}"
}


resource "kong_plugin" "rate_limit" {
	name        = "rate-limiting"
	route_id = "${kong_route.route.id}"
	config_json = <<EOT
	{
		"second": 14,
		"hour" : 4000
	}
EOT
}
`

const testImportPluginForASpecificApiConfig = `
resource "kong_api" "api" {
	name 	= "TestApi"
  	hosts   = [ "example.com" ]
	uris 	= [ "/example" ]
	methods = [ "GET", "POST" ]
	upstream_url = "http://localhost:4140"
	strip_uri = false
	preserve_host = false
	retries = 3
	upstream_connect_timeout = 60000
	upstream_send_timeout = 30000
	upstream_read_timeout = 10000
	https_only = false
	http_if_terminated = false
}

resource "kong_plugin" "basic_auth" {
	name   = "basic-auth"
	api_id = "${kong_api.api.id}"
	config = {
		hide_credentials = "false"
	}
}
`
const testImportPluginForJson = `
resource "kong_api" "api" {
	name 	= "TestApi"
  	hosts   = [ "example.com" ]
	uris 	= [ "/example" ]
	methods = [ "GET", "POST" ]
	upstream_url = "http://localhost:4140"
	strip_uri = false
	preserve_host = false
	retries = 3
	upstream_connect_timeout = 60000
	upstream_send_timeout = 30000
	upstream_read_timeout = 10000
	https_only = false
	http_if_terminated = false
}

resource "kong_plugin" "basic_auth_json" {
	name   = "basic-auth"
	api_id = "${kong_api.api.id}"
	config_json = <<EOT
{
	"hide_credentials": true,
	"anonymous": ""
}
EOT
}
`
