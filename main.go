package main

import (
	"github.com/h0x91b-wix/terraform-provider-netbox-wix/plugin/providers/netbox"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Info("Loading terraform-provider-netbox plugin")

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: netbox.Provider,
	})
}
