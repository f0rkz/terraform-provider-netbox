package netbox

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

// we need to convert some int64 pointers to nil in case Terraform SDK passed
// value is 0, this due to https://github.com/hashicorp/terraform-plugin-sdk/issues/90
func nilFromInt64Ptr(i *int64) *int64 {
	if *i == int64(0) {
		return nil
	}

	return i
}

// resourceNetboxIPAMIpAddress is the core Terraform resource structure for the netbox_ipam_ip_address resource.
func resourceNetboxIPAMIPAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIPAMIPAddressCreate,
		Read:   resourceNetboxIPAMIPAddressRead,
		Update: resourceNetboxIPAMIPAddressUpdate,
		Delete: resourceNetboxIPAMIPAddressDelete,
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

// resourceNetboxIPAMIpAddressCreate creates a new IP Address in Netbox.
func resourceNetboxIPAMIPAddressCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	address := d.Get("address").(string)
	vrfID := int64(d.Get("vrf_id").(int))
	tenantID := int64(d.Get("tenant_id").(int))
	status := int64(d.Get("status").(int))
	roleID := int64(d.Get("role_id").(int))
	description := d.Get("description").(string)
	natInsideID := int64(d.Get("nat_inside_ip_address_id").(int))
	natOutsideID := int64(d.Get("nat_outside_ip_address_id").(int))

	var parm = ipam.NewIPAMIPAddressesCreateParams().WithData(
		&models.WritableIPAddress{
			Address:     &address,
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

	log.Debugf("Executing IPAMIPAddressesCreate against Netbox: %v", parm)

	out, err := netboxClient.IPAM.IPAMIPAddressesCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IPAMIPAddressesCreate: %v", err)

		return err
	}

	d.SetId(fmt.Sprintf("ipam/ip-address/%d", out.Payload.ID))
	d.Set("ip_address_id", out.Payload.ID)

	log.Debugf("Done Executing IPAMIPAddressesCreate: %v", out)

	return nil
}

// resourceNetboxIPAMIpAddressUpdate applies updates to a IP Address by ID when deltas are detected by Terraform.
func resourceNetboxIPAMIPAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("ip_address_id").(int))

	address := d.Get("address").(string)
	vrfID := int64(d.Get("vrf_id").(int))
	tenantID := int64(d.Get("tenant_id").(int))
	status := int64(d.Get("status").(int))
	roleID := int64(d.Get("role_id").(int))
	description := d.Get("description").(string)
	natInsideID := int64(d.Get("nat_inside_ip_address_id").(int))
	natOutsideID := int64(d.Get("nat_outside_ip_address_id").(int))

	var parm = ipam.NewIPAMIPAddressesUpdateParams().
		WithID(id).
		WithData(
			&models.WritableIPAddress{
				Address:     &address,
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

	log.Debugf("Executing IPAMIPAddressesUpdate against Netbox: %v", parm)

	out, err := netboxClient.IPAM.IPAMIPAddressesUpdate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IPAMIPAddressesUpdate: %v", err)

		return err
	}

	log.Debugf("Done Executing IPAMIPAddressesUpdate: %v", out)

	return nil
}

// resourceNetboxIPAMIpAddressRead reads an existing IP Address by ID.
func resourceNetboxIPAMIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("ip_address_id").(int))

	var readParams = ipam.NewIPAMIPAddressesReadParams().WithID(id)

	readResult, err := netboxClient.IPAM.IPAMIPAddressesRead(readParams, nil)

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

	var status int64
	if readResult.Payload.Status != nil {
		status = *readResult.Payload.Status.Value
	}
	d.Set("status", status)

	var roleID int64
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

// resourceNetboxIPAMIpAddressDelete deletes an existing IP Address by ID.
func resourceNetboxIPAMIPAddressDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debugf("Deleting IpAddress: %v\n", d)

	id := int64(d.Get("ip_address_id").(int))

	var deleteParameters = ipam.NewIPAMIPAddressesDeleteParams().WithID(id)

	c := meta.(*ProviderNetboxClient).client

	out, err := c.IPAM.IPAMIPAddressesDelete(deleteParameters, nil)

	if err != nil {
		log.Debugf("Failed to execute IPAMIpAddresssDelete: %v", err)
	}

	log.Debugf("Done Executing IPAMIpAddresssDelete: %v", out)

	return nil
}
