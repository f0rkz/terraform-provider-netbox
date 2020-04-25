package netbox

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
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

// resourceNetboxIpamIpAddressCreate creates a new IP Address in Netbox.
func resourceNetboxIpamIPAddressCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	address := d.Get("address").(string)
	vrfID := int64(d.Get("vrf_id").(int))
	tenantID := int64(d.Get("tenant_id").(int))
	status := d.Get("status").(string)
	role := d.Get("role").(string)
	description := d.Get("description").(string)
	natInsideID := int64(d.Get("nat_inside_ip_address_id").(int))
	natOutsideID := int64(d.Get("nat_outside_ip_address_id").(int))
	interfaceID := int64(d.Get("interface_id").(int))

	var parm = ipam.NewIpamIPAddressesCreateParams().WithData(
		&models.WritableIPAddress{
			Address:     &address,
			Description: description,
			Vrf:         &vrfID,
			Tenant:      nilFromInt64Ptr(&tenantID),
			Status:      status,
			Role:        role,
			NatInside:   nilFromInt64Ptr(&natInsideID),
			NatOutside:  &natOutsideID,
			Interface:   nilFromInt64Ptr(&interfaceID),
			Tags:        []string{},
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
	status := d.Get("status").(string)
	role := d.Get("role").(string)
	description := d.Get("description").(string)
	natInsideID := int64(d.Get("nat_inside_ip_address_id").(int))
	natOutsideID := int64(d.Get("nat_outside_ip_address_id").(int))
	interfaceID := int64(d.Get("interface_id").(int))

	var parm = ipam.NewIpamIPAddressesUpdateParams().
		WithID(id).
		WithData(
			&models.WritableIPAddress{
				Address:     &address,
				Description: description,
				Vrf:         &vrfID,
				Tenant:      nilFromInt64Ptr(&tenantID),
				Status:      status,
				Role:        role,
				NatInside:   nilFromInt64Ptr(&natInsideID),
				NatOutside:  &natOutsideID,
				Interface:   nilFromInt64Ptr(&interfaceID),
				Tags:        []string{},
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
		status = strings.ToLower(*readResult.Payload.Status.Label)
	}
	d.Set("status", status)

	var role string
	if readResult.Payload.Role != nil {
		role = strings.ToLower(*readResult.Payload.Role.Label)
	}
	d.Set("role", role)

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

	var interfaceID int64
	if readResult.Payload.Interface != nil {
		interfaceID = readResult.Payload.Interface.ID
	}
	d.Set("interface_id", interfaceID)

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
