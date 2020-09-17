package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/ledor473/terraform-provider-lightstep/internal/lightstep"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: lightstep.New})
}
