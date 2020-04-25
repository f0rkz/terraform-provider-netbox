package netbox

import (
	"errors"
	//"fmt"
	"log"
	"strconv"
	"strings"

	// "errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/models"
)

func dataSourceNetboxIPAddress() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNetboxIPAddressesRead,
		Schema: bareIPAddressesSchema(),
	}
}

func dataSourceNetboxIPAddressParse(d *schema.ResourceData, obj *models.IPAddress) {
	d.SetId(strconv.FormatInt(obj.ID, 10))
	d.Set("created", obj.Created.String())
	d.Set("description", obj.Description)
	d.Set("status", *obj.Status.Label)
	d.Set("family", *obj.Family.Label)
	d.Set("address", *obj.Address)
	d.Set("last_updated", obj.LastUpdated)

	if obj.Vrf != nil {
		d.Set("vrf", *obj.Vrf.Name)
	}

	if obj.Role != nil {
		d.Set("role", *obj.Role.Label)
	}

	if obj.Tenant != nil {
		d.Set("tenant", *obj.Tenant.Name)
	}

	// interface ?

	log.Printf("Finished parsing results from IpamIPAddressesRead")
}

func dataSourceNetboxIPAddressAttrPrep(in string) (out string) {
	lowerstr := strings.ToLower(in)
	out = strings.Replace(lowerstr, " ", "-", -1)

	return
}

// Read will fetch the data of a resource.
func dataSourceNetboxIPAddressesRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderNetboxClient).client

	// primary key lookup, direct
	if id, idOk := d.GetOk("id"); idOk {
		parm := ipam.NewIpamIPAddressesReadParams()
		parm.SetID(int64(id.(int)))

		out, err := c.Ipam.IpamIPAddressesRead(parm, nil)

		if err != nil {
			log.Printf("error from IpamIPAddressesRead: %v\n", err)
			return err
		}

		dataSourceNetboxIPAddressParse(d, out.Payload)
	} else { // anything else, requires a search
		param := ipam.NewIpamIPAddressesListParams()

		// Add any lookup params

		if query, queryOk := d.GetOk("query"); queryOk {
			query_str := query.(string)
			param.SetQ(&query_str)
		}

		if family, familyOk := d.GetOk("family"); familyOk {
			family_str := family.(string)
			param.SetFamily(&family_str)
		}

		if parent, parentOk := d.GetOk("parent"); parentOk {
			parent_str := parent.(string)
			param.SetParent(&parent_str)
		}

		if tenant, tenantOk := d.GetOk("tenant"); tenantOk {
			tenant_str := dataSourceNetboxIPAddressAttrPrep(tenant.(string))
			param.SetTenant(&tenant_str)
		}

		//if site, siteOk := d.GetOk("site"); siteOk {
		//  site_str := dataSourceNetboxIPAddressAttrPrep(site.(string))
		//  param.SetSite(&site_str)
		//}

		//if role, roleOk := d.GetOk("role"); roleOk {
		//  role_str := dataSourceNetboxIPAddressAttrPrep(role.(string))
		//  param.SetRole(&role_str)
		//}

		// limit to 2
		limit := int64(2)
		param.SetLimit(&limit)

		out, err := c.Ipam.IpamIPAddressesList(param, nil)

		if err != nil {
			log.Printf("error from IpamIPAddressesList: %v\n", err)
			return err
		}

		if *out.Payload.Count == 0 {
			return errors.New("IPAddress not found")
		} else if *out.Payload.Count > 1 {
			return errors.New("More than one prefix matches search terms, please narrow")
		}

		dataSourceNetboxIPAddressParse(d, out.Payload.Results[0])
	}

	return nil
}

func bareIPAddressesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
		},
		"created": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"address": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"family": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"vrf": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"status": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"last_updated": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"query": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"tenant": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"role": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"parent": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}
