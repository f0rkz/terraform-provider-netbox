package netbox

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/netbox-community/go-netbox/netbox/client/virtualization"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func resourceNetboxVirtualizationCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxVirtualizationClusterCreate,
		Read:   resourceNetboxVirtualizationClusterRead,
		Update: resourceNetboxVirtualizationClusterUpdate,
		Delete: resourceNetboxVirtualizationClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"comments": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"site_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"type_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceNetboxVirtualizationClusterCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	comments := d.Get("comments").(string)
	groupID := int64(d.Get("group_id").(int))
	name := d.Get("name").(string)
	siteID := int64(d.Get("site_id").(int))
	typeID := int64(d.Get("type_id").(int))

	var parm = virtualization.NewVirtualizationClustersCreateParams().WithData(
		&models.WritableCluster{
			Comments: comments,
			Group:    nilFromInt64Ptr(&groupID),
			Name:     &name,
			Site:     nilFromInt64Ptr(&siteID),
			Type:     nilFromInt64Ptr(&typeID),
			Tags:     []string{},
		},
	)

	log.Debugf("Executing VirtualizationClustersCreate againts Netbox: %v", parm)

	out, err := netboxClient.Virtualization.VirtualizationClustersCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute VirtualizationClustersCreate: %v", err)

		return err
	}

	d.SetId(fmt.Sprintf("virtualization/clusters/%d", out.Payload.ID))
	d.Set("cluster_id", out.Payload.ID)

	log.Debugf("Done Executing VirtualizationClustersCreate: %v", out)

	return nil
}

func resourceNetboxVirtualizationClusterRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("cluster_id").(int))

	var parm = virtualization.NewVirtualizationClustersReadParams().WithID(id)

	result, err := netboxClient.Virtualization.VirtualizationClustersRead(parm, nil)

	if err != nil {
		log.Debugf("Error fetching Virtualization Cluster ID # %d from Netbox = %v", id, err)
		return err
	}

	d.Set("name", result.Payload.Name)
	d.Set("comments", result.Payload.Comments)

	var groupID int64
	if result.Payload.Group != nil {
		groupID = result.Payload.Group.ID
	}
	d.Set("group_id", groupID)

	var siteID int64
	if result.Payload.Site != nil {
		siteID = result.Payload.Site.ID
	}
	d.Set("site_id", siteID)

	var typeID int64
	if result.Payload.Type != nil {
		typeID = result.Payload.Type.ID
	}
	d.Set("type_id", typeID)

	return nil
}

func resourceNetboxVirtualizationClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("cluster_id").(int))

	comments := d.Get("comments").(string)
	groupID := int64(d.Get("group_id").(int))
	name := d.Get("name").(string)
	siteID := int64(d.Get("site_id").(int))
	typeID := int64(d.Get("type_id").(int))

	var parm = virtualization.NewVirtualizationClustersUpdateParams().
		WithID(id).
		WithData(
			&models.WritableCluster{
				Comments: comments,
				Group:    nilFromInt64Ptr(&groupID),
				Name:     &name,
				Site:     nilFromInt64Ptr(&siteID),
				Type:     nilFromInt64Ptr(&typeID),
				Tags:     []string{},
			},
		)

	log.Debugf("Executing VirtualizationClustersUpdate againts Netbox: %v", parm)

	out, err := netboxClient.Virtualization.VirtualizationClustersUpdate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute VirtualizationClustersUpdate: %v", err)

		return err
	}

	log.Debugf("Done Executing VirtualizationClustersUpdate: %v", out)

	return nil
}

func resourceNetboxVirtualizationClusterDelete(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client
	log.Debugf("Deleting Virtualization Cluster: %v\n", d)

	id := int64(d.Get("cluster_id").(int))

	var parm = virtualization.NewVirtualizationClustersDeleteParams().WithID(id)

	out, err := netboxClient.Virtualization.VirtualizationClustersDelete(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute VirtualizationClustersDelete: %v", err)
	}

	log.Debugf("Done Executing VirtualizationClustersDelete: %v", out)

	return nil
}
