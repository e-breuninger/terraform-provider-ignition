package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/e-breuninger/terraform-provider-ignition/internal"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: internal.Provider,
	})
}
