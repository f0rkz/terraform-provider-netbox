package netbox

import (
	"fmt"
	//"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/netbox-community/go-netbox/netbox/client/virtualization"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxVirtualizationVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxVirtualizationVirtualMachineCreate,
		Read:   resourceNetboxVirtualizationVirtualMachineRead,
		Update: resourceNetboxVirtualizationVirtualMachineUpdate,
		Delete: resourceNetboxVirtualizationVirtualMachineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"virtual_machine_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"comments": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_gb": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory_mb": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vcpus": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"site": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_ip4_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"role_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"offline",
					"active",
					"planned",
					"staged",
					"failed",
					"decommissioning",
				}, true),
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceNetboxVirtualizationVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	clusterID := int64(d.Get("cluster_id").(int))
	diskGB := int64(d.Get("disk_gb").(int))
	memoryMB := int64(d.Get("memory_mb").(int))
	vcpus := int64(d.Get("vcpus").(int))
	name := d.Get("name").(string)
	primaryIp4ID := int64(d.Get("primary_ip4_id").(int))
	roleID := int64(d.Get("role_id").(int))
	tenantID := int64(d.Get("tenant_id").(int))

	var parm = virtualization.NewVirtualizationVirtualMachinesCreateParams().WithData(
		&models.WritableVirtualMachineWithConfigContext{
			Cluster:    &clusterID,
			Comments:   d.Get("comments").(string),
			Disk:       nilFromInt64Ptr(&diskGB),
			Memory:     nilFromInt64Ptr(&memoryMB),
			Vcpus:      nilFromInt64Ptr(&vcpus),
			Name:       &name,
			PrimaryIp4: nilFromInt64Ptr(&primaryIp4ID),
			Role:       nilFromInt64Ptr(&roleID),
			Status:     d.Get("status").(string),
			Tenant:     nilFromInt64Ptr(&tenantID),
			Tags:       []string{},
		},
	)

	log.Debugf("Executing VirtualizationVirtualMachinesCreate againts Netbox: %v", parm)

	out, err := netboxClient.Virtualization.VirtualizationVirtualMachinesCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute VirtualizationVirtualMachinesCreate: %v", err)

		return err
	}

	d.SetId(fmt.Sprintf("virtualization/virtual-machines/%d", out.Payload.ID))
	d.Set("virtual_machine_id", out.Payload.ID)

	log.Debugf("Done Executing VirtualizationVirtualMachinesCreate: %v", out)

	return nil
}

func resourceNetboxVirtualizationVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("virtual_machine_id").(int))
	//id_string := strconv.FormatInt(id, 10)

	var parm = virtualization.NewVirtualizationVirtualMachinesReadParams().WithID(id)

	result, err := netboxClient.Virtualization.VirtualizationVirtualMachinesRead(parm, nil)

	if err != nil {
		log.Debugf("Error fetching Virtualization VirtualMachine ID # %d from Netbox = %v", id, err)
		return err
	}

	var clusterID int64
	if result.Payload.Cluster != nil {
		clusterID = result.Payload.Cluster.ID
	}
	d.Set("cluster_id", clusterID)

	d.Set("comments", result.Payload.Comments)

	d.Set("disk_gb", result.Payload.Disk)
	d.Set("memory_mb", result.Payload.Memory)
	d.Set("vcpus", result.Payload.Vcpus)
	d.Set("name", result.Payload.Name)
	d.Set("site", result.Payload.Site)

	var primaryIp4ID int64
	if result.Payload.PrimaryIp4 != nil {
		primaryIp4ID = result.Payload.PrimaryIp4.ID
	}
	d.Set("primary_ip4_id", primaryIp4ID)

	var siteID int64
	if result.Payload.Role != nil {
		siteID = result.Payload.Role.ID
	}
	d.Set("site_id", siteID)

	d.Set("status", result.Payload.Status)

	var tenantID int64
	if result.Payload.Tenant != nil {
		tenantID = result.Payload.Tenant.ID
	}
	d.Set("tenant_id", tenantID)

	return nil
}

func resourceNetboxVirtualizationVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("virtual_machine_id").(int))

	clusterID := int64(d.Get("cluster_id").(int))
	diskGB := int64(d.Get("disk_gb").(int))
	memoryMB := int64(d.Get("memory_mb").(int))
	vcpus := int64(d.Get("vcpus").(int))
	name := d.Get("name").(string)
	primaryIp4ID := int64(d.Get("primary_ip4_id").(int))
	roleID := int64(d.Get("role_id").(int))
	tenantID := int64(d.Get("tenant_id").(int))

	var parm = virtualization.NewVirtualizationVirtualMachinesUpdateParams().
		WithID(id).
		WithData(
			&models.WritableVirtualMachineWithConfigContext{
				Cluster:    &clusterID,
				Comments:   d.Get("comments").(string),
				Disk:       nilFromInt64Ptr(&diskGB),
				Memory:     nilFromInt64Ptr(&memoryMB),
				Vcpus:      nilFromInt64Ptr(&vcpus),
				Name:       &name,
				PrimaryIp4: nilFromInt64Ptr(&primaryIp4ID),
				Role:       nilFromInt64Ptr(&roleID),
				Status:     d.Get("status").(string),
				Tenant:     nilFromInt64Ptr(&tenantID),
				Tags:       []string{},
			},
		)

	log.Debugf("Executing VirtualizationVirtualMachinesUpdate againts Netbox: %v", parm)

	out, err := netboxClient.Virtualization.VirtualizationVirtualMachinesUpdate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute VirtualizationVirtualMachinesUpdate: %v", err)

		return err
	}

	log.Debugf("Done Executing VirtualizationVirtualMachinesUpdate: %v", out)

	return nil
}

func resourceNetboxVirtualizationVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client
	log.Debugf("Deleting Virtualization VirtualMachine: %v\n", d)

	id := int64(d.Get("virtual_machine_id").(int))

	var parm = virtualization.NewVirtualizationVirtualMachinesDeleteParams().WithID(id)

	out, err := netboxClient.Virtualization.VirtualizationVirtualMachinesDelete(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute VirtualizationVirtualMachinesDelete: %v", err)
	}

	log.Debugf("Done Executing VirtualizationVirtualMachinesDelete: %v", out)

	return nil
}
