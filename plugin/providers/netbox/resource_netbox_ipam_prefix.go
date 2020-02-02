package netbox

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/h0x91b-wix/go-netbox/netbox/client/ipam"
	"github.com/h0x91b-wix/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// resourceNetboxIpamPrefix is the core Terraform resource structure for the netbox_ipam_Prefix_domain resource.
func resourceNetboxIpamPrefix() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamPrefixCreate,
		Read:   resourceNetboxIpamPrefixRead,
		Update: resourceNetboxIpamPrefixUpdate,
		Delete: resourceNetboxIpamPrefixDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"prefix": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"prefix_id": &schema.Schema{
				Type:     schema.TypeInt,
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
			"is_pool": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
		},
	}
}

// resourceNetboxIpamPrefixCreate creates a new Prefix in Netbox.
func resourceNetboxIpamPrefixCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	prefix := d.Get("prefix").(string)
	description := d.Get("description").(string)
	vrfID := int64(d.Get("vrf_id").(int))
	isPool := d.Get("is_pool").(bool)
	//status := d.Get("status").(string)
	tenantID := int64(d.Get("tenant_id").(int))

	var parm = ipam.NewIpamPrefixesCreateParams().WithData(
		&models.WritablePrefix{
			Prefix:      &prefix,
			Description: description,
			IsPool:      isPool,
			Tags:        []string{},
			Vrf:         &vrfID,
			Tenant:      &tenantID,
		},
	)

	log.Debugf("Executing IpamPrefixesCreate against Netbox: %v", parm)

	out, err := netboxClient.Ipam.IpamPrefixesCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamPrefixesCreate: %v", err)

		return err
	}

	// TODO Probably a better way to parse this ID
	d.SetId(fmt.Sprintf("ipam/prefix/%d", out.Payload.ID))
	d.Set("prefix_id", out.Payload.ID)

	log.Debugf("Done Executing IpamPrefixesCreate: %v", out)

	return nil
}

// resourceNetboxIpamPrefixUpdate applies updates to a Prefix by ID when deltas are detected by Terraform.
func resourceNetboxIpamPrefixUpdate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("prefix_id").(int))

	prefix := d.Get("prefix").(string)
	description := d.Get("description").(string)
	vrfID := int64(d.Get("vrf_id").(int))
	isPool := d.Get("is_pool").(bool)
	//status := d.Get("status").(string)
	tenantID := int64(d.Get("tenant_id").(int))

	var parm = ipam.NewIpamPrefixesUpdateParams().
		WithID(id).
		WithData(
			&models.WritablePrefix{
				Prefix:      &prefix,
				Description: description,
				IsPool:      isPool,
				Tags:        []string{},
				Vrf:         &vrfID,
				Tenant:      &tenantID,
			},
		)

	log.Debugf("Executing IpamPrefixesUpdate against Netbox: %v", parm)

	out, err := netboxClient.Ipam.IpamPrefixesUpdate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamPrefixesUpdate: %v", err)

		return err
	}

	log.Debugf("Done Executing IpamPrefixesUpdate: %v", out)

	return nil
}

// resourceNetboxIpamPrefixRead reads an existing Prefix by ID.
func resourceNetboxIpamPrefixRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("prefix_id").(int))

	var readParams = ipam.NewIpamPrefixesReadParams().WithID(id)

	readResult, err := netboxClient.Ipam.IpamPrefixesRead(readParams, nil)

	if err != nil {
		log.Debugf("Error fetching Prefix ID # %d from Netbox = %v", id, err)
		return err
	}

	var vrfID int64
	if readResult.Payload.Vrf != nil {
		vrfID = readResult.Payload.Vrf.ID
	}

	d.Set("prefix", readResult.Payload.Prefix)
	d.Set("description", readResult.Payload.Description)
	d.Set("vrf_id", vrfID)
	d.Set("is_pool", readResult.Payload.IsPool)

	var tenantID int64
	if readResult.Payload.Tenant != nil {
		tenantID = readResult.Payload.Tenant.ID
	}
	d.Set("tenant_id", tenantID)

	return nil
}

// resourceNetboxIpamPrefixDelete deletes an existing Prefix by ID.
func resourceNetboxIpamPrefixDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debugf("Deleting Prefix: %v\n", d)

	id := int64(d.Get("prefix_id").(int))

	var deleteParameters = ipam.NewIpamPrefixesDeleteParams().WithID(id)

	c := meta.(*ProviderNetboxClient).client

	out, err := c.Ipam.IpamPrefixesDelete(deleteParameters, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamPrefixesDelete: %v", err)
	}

	log.Debugf("Done Executing IpamPrefixesDelete: %v", out)

	return nil
}
