module github.com/peltzi/terraform-provider-netbox

go 1.13

require (
	github.com/go-openapi/runtime v0.19.14
	github.com/go-openapi/strfmt v0.19.5
	github.com/hashicorp/terraform-plugin-sdk v1.6.0
	github.com/netbox-community/go-netbox v0.0.0
	github.com/sirupsen/logrus v1.4.2
)

replace github.com/netbox-community/go-netbox v0.0.0 => github.com/peltzi/go-netbox v0.0.0-20200425085534-1915fbcdfd46
