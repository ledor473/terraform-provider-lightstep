package lightstep

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// New -
func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("LIGHTSTEP_API_KEY", nil),
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LIGHTSTEP_PROJECT", nil),
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LIGHTSTEP_ORGANIZATION", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"lightstep_stream":        resourceStream(),
			"lightstep_condition":     resourceCondition(),
			"lightstep_workflow_link": resourceWorkflowLink(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"lightstep_stream":        dataSourceStream(),
			"lightstep_condition":     dataSourceCondition(),
			"lightstep_workflow_link": dataSourceWorkflowLink(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apikey := d.Get("api_key").(string)
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if apikey == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "API Key is mandatory",
		})
	}
	if org == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Organization is mandatory",
		})
	}
	if project == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Project is mandatory",
		})
	}

	if diags.HasError() {
		return nil, diags
	}

	c := Config{APIKey: apikey, Organization: org, Project: project}
	return &c, nil
}
