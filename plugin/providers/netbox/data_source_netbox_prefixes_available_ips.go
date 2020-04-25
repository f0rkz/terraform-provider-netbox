package netbox

import (
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
)

func dataSourceNetboxPrefixesAvailableIps() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNetboxPrefixesAvailableIpsRead,
		Schema: dataSourceNetboxPrefixesAvailableIpsSchema(),
	}
}

func dataSourceNetboxPrefixesAvailableIpsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"prefix_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"v4_addresses": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"v6_addresses": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func dataSourceNetboxPrefixesAvailableIpsRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("prefix_id").(int))

	var readParams = ipam.NewIpamPrefixesAvailableIpsReadParams().WithID(id)

	log.Debugf("Executing IpamPrefixesAvailableIpsRead againts Netbox")

	result, err := netboxClient.Ipam.IpamPrefixesAvailableIpsRead(readParams, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamPrefixesAvailableIpsRead againts Netbox: %v", err)

		return err
	}

	log.Debugf("Result: %v", result.Payload)

	v4_addresses := make([]string, 0)
	v6_addresses := make([]string, 0)

	for _, address := range *result.Payload {
		if address.Family == 4 {
			v4_addresses = append(v4_addresses, *address.Address)
		}

		if address.Family == 6 {
			v6_addresses = append(v6_addresses, *address.Address)
		}
	}

	data_id := strconv.Itoa(int(d.Get("prefix_id").(int)))
	d.SetId(data_id)
	d.Set("v4_addresses", &v4_addresses)
	d.Set("v6_addresses", &v6_addresses)

	return nil
}
