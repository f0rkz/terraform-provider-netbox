package netbox

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/netbox-community/go-netbox/netbox/client/virtualization"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxVirtualizationInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxVirtualizationInterfaceCreate,
		Read:   resourceNetboxVirtualizationInterfaceRead,
		Update: resourceNetboxVirtualizationInterfaceUpdate,
		Delete: resourceNetboxVirtualizationInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"virtual_machine_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"interface_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceNetboxVirtualizationInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	name := d.Get("name").(string)
	virtual_machine_id := int64(d.Get("virtual_machine_id").(int))

	var parm = virtualization.NewVirtualizationInterfacesCreateParams().WithData(
		&models.WritableVirtualMachineInterface{
			Name:           &name,
			VirtualMachine: &virtual_machine_id,
			TaggedVlans:    []int64{},
			Tags:           []string{},
		},
	)

	log.Debugf("Executin VirtualizationInterfacesCreate againts Netbox: %v", parm)

	out, err := netboxClient.Virtualization.VirtualizationInterfacesCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute VirtualizationInterfacesCreate: %v", err)

		return err
	}

	d.SetId(fmt.Sprintf("dcim/interfaces/%d", out.Payload.ID))
	d.Set("interface_id", out.Payload.ID)

	return nil
}

func resourceNetboxVirtualizationInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("interface_id").(int))

	var parm = virtualization.NewVirtualizationInterfacesReadParams().WithID(id)

	result, err := netboxClient.Virtualization.VirtualizationInterfacesRead(parm, nil)

	if err != nil {
		log.Debugf("Error fetching Virtualization Interface ID # %d from Netbox = %v", id, err)
		return err
	}

	var virtualMachineID int64
	if result.Payload.VirtualMachine != nil {
		virtualMachineID = result.Payload.VirtualMachine.ID
	}
	d.Set("virtual_machine_id", virtualMachineID)

	d.Set("name", result.Payload.Name)
	d.Set("interface_id", result.Payload.ID)

	return nil
}

func resourceNetboxVirtualizationInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("interface_id").(int))

	name := d.Get("name").(string)
	virtual_machine_id := int64(d.Get("virtual_machine_id").(int))

	var parm = virtualization.NewVirtualizationInterfacesUpdateParams().
		WithID(id).
		WithData(
			&models.WritableVirtualMachineInterface{
				Name:           &name,
				VirtualMachine: &virtual_machine_id,
				TaggedVlans:    []int64{},
				Tags:           []string{},
			},
		)

	log.Debugf("Executing VirtualizationInterfacesUpdate againts Netbox: %v", parm)

	out, err := netboxClient.Virtualization.VirtualizationInterfacesUpdate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute VirtualizationInterfacesUpdate: %v", err)

		return err
	}

	log.Debugf("Done Executing VirtualizationInterfacesUpdate: %v", out)

	return nil
}

func resourceNetboxVirtualizationInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client
	log.Debugf("Deleting Virtualization Interface: %v\n", d)

	id := int64(d.Get("interface_id").(int))

	var parm = virtualization.NewVirtualizationInterfacesDeleteParams().WithID(id)

	out, err := netboxClient.Virtualization.VirtualizationInterfacesDelete(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute VirtualizationInterfacesDelete: %v", err)
	}

	log.Debugf("Done Executing VirtualizationInterfacesDelete: %v", out)

	return nil
}
