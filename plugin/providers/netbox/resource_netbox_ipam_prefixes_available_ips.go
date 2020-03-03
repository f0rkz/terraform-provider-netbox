package netbox

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxIPAMPrefixesAvailableIps() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIPAMPrefixesAvailableIpsCreate,
		Read:   resourceNetboxIPAMIPAddressRead,
		Update: resourceNetboxIPAMIPAddressUpdate,
		Delete: resourceNetboxIPAMIPAddressDelete,
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
				Type:     schema.TypeInt,
				Optional: true,
			},
			"role_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
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
		},
	}
}

func resourceNetboxIPAMPrefixesAvailableIpsCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	prefix_id := int64(d.Get("prefix_id").(int))
	vrfID := int64(d.Get("vrf_id").(int))
	tenantID := int64(d.Get("tenant_id").(int))
	status := int64(d.Get("status").(int))
	roleID := int64(d.Get("role_id").(int))
	description := d.Get("description").(string)
	natInsideID := int64(d.Get("nat_inside_ip_address_id").(int))
	natOutsideID := int64(d.Get("nat_outside_ip_address_id").(int))

	var parm = ipam.NewIPAMPrefixesAvailableIpsCreateParams().
		WithID(prefix_id).
		WithData(
			&models.WritableAvailableIPAddress{
				Description: description,
				Vrf:         &vrfID,
				Tenant:      nilFromInt64Ptr(&tenantID),
				Status:      status,
				Role:        nilFromInt64Ptr(&roleID),
				NatInside:   nilFromInt64Ptr(&natInsideID),
				NatOutside:  &natOutsideID,
				// TODO Interface
				Tags: []string{},
			},
		)

	log.Debugf("Executing IPAMPrefixesAvailableIpsCreate against Netbox: %v", parm)

	out, err := netboxClient.IPAM.IPAMPrefixesAvailableIpsCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IPAMPrefixesAvailableIpsCreate: %v", err)

		return err
	}

	d.SetId(fmt.Sprintf("ipam/ip-address/%d", out.Payload.ID))
	d.Set("ip_address_id", out.Payload.ID)

	log.Debugf("Done Executing IPAMPrefixesAvailableIpsCreate: %v", out)

	return nil

}
