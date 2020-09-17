package lightstep

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ledor473/lightstep-api-go/pkg/v0.2/client/workflow_links"
)

func dataSourceWorkflowLink() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkflowLinkRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				// TODO: Define DiffSuppressFunc to ignore whitespace changes
				//DiffSuppressFunc: ,
			},
		},
	}
}

func dataSourceWorkflowLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading LightStep Workflow Links")

	searchName := d.Get("name").(string)
	p := &workflow_links.ListWorkflowLinksParams{Organization: config.Organization, Project: config.Project}

	// TODO: Retry all API Calls with backoff
	resp, err := api.WorkflowLinks.ListWorkflowLinks(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool

	for _, s := range resp.GetPayload().Data {
		if searchName != "" && s.Attributes.Name == searchName {
			d.SetId(s.ID)
			d.Set("name", s.Attributes.Name)
			d.Set("url", s.Attributes.URL)
			if s.Attributes.Rules != nil {
				if rules, err := marshalRules(s.Attributes.Rules); err == nil {
					d.Set("rules", string(rules))
				} else {
					return diag.FromErr(err)
				}
			}

			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("Unable to find a Workflow Link with the name: %s", searchName)
	}

	return nil
}
