package netbox

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var descriptions map[string]string

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:         providerSchema(),
		DataSourcesMap: providerDataSourcesMap(),
		ResourcesMap:   providerResources(),
		ConfigureFunc:  providerConfigure,
	}
}

// List of supported configuration fields for your provider.
func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"app_id": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("NETBOX_APP_ID", nil),
			Description: "API key used to access Netbox, generated under Admin -> Users -> Tokens and assigned to a user",
		},
		"endpoint": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("NETBOX_ENDPOINT_ADDR", nil),
			Description: "Endpoint of your Netbox instance",
		},
		/*
			"timeout": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
		*/
	}
}

// List of supported resources and their configuration fields.
func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		// IPAM
		"netbox_ipam_ip_address":             resourceNetboxIPAMIPAddress(),
		"netbox_ipam_prefixes_available_ips": resourceNetboxIPAMPrefixesAvailableIps(),
	}
}

// List of supported data sources and their configuration fields.
func providerDataSourcesMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"netbox_ip_address":             dataSourceNetboxIPAddress(),
		"netbox_prefixes_available_ips": dataSourceNetboxPrefixesAvailableIps(),
	}
}

// This is the function used to fetch the configuration params given
// to our provider which we will use to initialise a dummy client that
// interacts with the API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AppID:    d.Get("app_id").(string),
		Endpoint: d.Get("endpoint").(string),
		//Timeout:  d.Get("timeout").(string),
	}
	return config.Client()
}
