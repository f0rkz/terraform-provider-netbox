package netbox

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxIpamPrefixesAvailableIps() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamPrefixesAvailableIpsCreate,
		Read:   resourceNetboxIpamIPAddressRead,
		Update: resourceNetboxIpamIPAddressUpdate,
		Delete: resourceNetboxIpamIPAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"prefix_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"ip_address_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vrf_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"active",
					"reserved",
					"deprecated",
					"dhcp",
				}, true),
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"loopback",
					"secondary",
					"anycast",
					"vip",
					"vrrp",
					"hsrp",
					"glbp",
					"carp",
				}, true),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"nat_inside_ip_address_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"nat_outside_ip_address_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"interface_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceNetboxIpamPrefixesAvailableIpsCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	prefix_id := int64(d.Get("prefix_id").(int))
	vrfID := int64(d.Get("vrf_id").(int))
	tenantID := int64(d.Get("tenant_id").(int))
	status := d.Get("status").(string)
	role := d.Get("role").(string)
	description := d.Get("description").(string)
	natInsideID := int64(d.Get("nat_inside_ip_address_id").(int))
	natOutsideID := int64(d.Get("nat_outside_ip_address_id").(int))
	interfaceID := int64(d.Get("interface_id").(int))

	var parm = ipam.NewIpamPrefixesAvailableIpsCreateParams().
		WithID(prefix_id).
		WithData(
			&models.WritableAvailableIPAddress{
				Description: description,
				Vrf:         &vrfID,
				Tenant:      nilFromInt64Ptr(&tenantID),
				Status:      status,
				Role:        role,
				NatInside:   nilFromInt64Ptr(&natInsideID),
				NatOutside:  &natOutsideID,
				Interface:   nilFromInt64Ptr(&interfaceID),
				// TODO Interface
				Tags: []string{},
			},
		)

	log.Debugf("Executing IpamPrefixesAvailableIpsCreate against Netbox: %v", parm)

	out, err := netboxClient.Ipam.IpamPrefixesAvailableIpsCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamPrefixesAvailableIpsCreate: %v", err)

		return err
	}

	d.SetId(fmt.Sprintf("ipam/ip-address/%d", out.Payload.ID))
	d.Set("ip_address_id", out.Payload.ID)
	d.Set("address", out.Payload.Address)
	d.Set("status", out.Payload.Status)

	log.Debugf("Done Executing IpamPrefixesAvailableIpsCreate: %v", out)

	return nil

}
