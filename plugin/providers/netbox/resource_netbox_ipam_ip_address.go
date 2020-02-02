package netbox

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/h0x91b-wix/go-netbox/netbox/client/ipam"
	"github.com/h0x91b-wix/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// resourceNetboxIpamIpAddress is the core Terraform resource structure for the netbox_ipam_ip_address resource.
func resourceNetboxIpamIPAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamIPAddressCreate,
		Read:   resourceNetboxIpamIPAddressRead,
		Update: resourceNetboxIpamIPAddressUpdate,
		Delete: resourceNetboxIpamIPAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ip_address_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
			// TODO interface - looks hard
			// TODO tags
			// TODO custom_fields
		},
	}
}

// resourceNetboxIpamIpAddressCreate creates a new IP Address in Netbox.
func resourceNetboxIpamIPAddressCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	address := d.Get("address").(string)
	vrfID := int64(d.Get("vrf_id").(int))
	tenantID := int64(d.Get("tenant_id").(int))
	status := string(d.Get("status").(string))
	roleID := string(d.Get("role_id").(string))
	description := d.Get("description").(string)
	natInsideID := int64(d.Get("nat_inside_ip_address_id").(int))
	natOutsideID := int64(d.Get("nat_outside_ip_address_id").(int))

	var parm = ipam.NewIpamIPAddressesCreateParams().WithData(
		&models.WritableIPAddress{
			Address:     &address,
			Description: description,
			Vrf:         &vrfID,
			Tenant:      &tenantID,
			Status:      status,
			Role:        roleID,
			NatInside:   &natInsideID,
			NatOutside:  &natOutsideID,
			// TODO Interface
			Tags: []string{},
		},
	)

	log.Debugf("Executing IpamIPAddressesCreate against Netbox: %v", parm)

	out, err := netboxClient.Ipam.IpamIPAddressesCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamIPAddressesCreate: %v", err)

		return err
	}

	d.SetId(fmt.Sprintf("ipam/ip-address/%d", out.Payload.ID))
	d.Set("ip_address_id", out.Payload.ID)

	log.Debugf("Done Executing IpamIPAddressesCreate: %v", out)

	return nil
}

// resourceNetboxIpamIpAddressUpdate applies updates to a IP Address by ID when deltas are detected by Terraform.
func resourceNetboxIpamIPAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("ip_address_id").(int))

	address := d.Get("address").(string)
	vrfID := int64(d.Get("vrf_id").(int))
	tenantID := int64(d.Get("tenant_id").(int))
	status := string(d.Get("status").(string))
	roleID := string(d.Get("role_id").(string))
	description := d.Get("description").(string)
	natInsideID := int64(d.Get("nat_inside_ip_address_id").(int))
	natOutsideID := int64(d.Get("nat_outside_ip_address_id").(int))

	var parm = ipam.NewIpamIPAddressesUpdateParams().
		WithID(id).
		WithData(
			&models.WritableIPAddress{
				Address:     &address,
				Description: description,
				Vrf:         &vrfID,
				Tenant:      &tenantID,
				Status:      status,
				Role:        roleID,
				NatInside:   &natInsideID,
				NatOutside:  &natOutsideID,
				// TODO Interface
				Tags: []string{},
			},
		)

	log.Debugf("Executing IpamIPAddressesUpdate against Netbox: %v", parm)

	out, err := netboxClient.Ipam.IpamIPAddressesUpdate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamIPAddressesUpdate: %v", err)

		return err
	}

	log.Debugf("Done Executing IpamIPAddressesUpdate: %v", out)

	return nil
}

// resourceNetboxIpamIpAddressRead reads an existing IP Address by ID.
func resourceNetboxIpamIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("ip_address_id").(int))

	var readParams = ipam.NewIpamIPAddressesReadParams().WithID(id)

	readResult, err := netboxClient.Ipam.IpamIPAddressesRead(readParams, nil)

	if err != nil {
		log.Debugf("Error fetching IpAddress ID # %d from Netbox = %v", id, err)
		return err
	}
	d.Set("address", readResult.Payload.Address)

	var vrfID int64
	if readResult.Payload.Vrf != nil {
		vrfID = readResult.Payload.Vrf.ID
	}
	d.Set("vrf_id", vrfID)

	var tenantID int64
	if readResult.Payload.Tenant != nil {
		tenantID = readResult.Payload.Tenant.ID
	}
	d.Set("tenant_id", tenantID)

	var status string
	if readResult.Payload.Status != nil {
		status = *readResult.Payload.Status.Value
	}
	d.Set("status", status)

	var roleID string
	if readResult.Payload.Role != nil {
		roleID = *readResult.Payload.Role.Value
	}
	d.Set("role_id", roleID)

	d.Set("description", readResult.Payload.Description)

	var natInsideID int64
	if readResult.Payload.NatInside != nil {
		natInsideID = readResult.Payload.NatInside.ID
	}
	d.Set("nat_inside_id", natInsideID)

	var natOutsideID int64
	if readResult.Payload.NatOutside != nil {
		natOutsideID = readResult.Payload.NatOutside.ID
	}
	d.Set("nat_outside_id", natOutsideID)

	return nil
}

// resourceNetboxIpamIpAddressDelete deletes an existing IP Address by ID.
func resourceNetboxIpamIPAddressDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debugf("Deleting IpAddress: %v\n", d)

	id := int64(d.Get("ip_address_id").(int))

	var deleteParameters = ipam.NewIpamIPAddressesDeleteParams().WithID(id)

	c := meta.(*ProviderNetboxClient).client

	out, err := c.Ipam.IpamIPAddressesDelete(deleteParameters, nil)

	if err != nil {
		log.Debugf("Failed to execute IpamIpAddresssDelete: %v", err)
	}

	log.Debugf("Done Executing IpamIpAddresssDelete: %v", out)

	return nil
}
