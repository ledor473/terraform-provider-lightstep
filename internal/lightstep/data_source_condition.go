package lightstep

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ledor473/lightstep-api-go/pkg/v0.2/client/conditions"
)

func dataSourceCondition() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConditionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stream_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eval_window_ms": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceConditionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading LightStep Conditions")

	searchName := d.Get("name").(string)
	p := &conditions.ListConditionsParams{Organization: config.Organization, Project: config.Project}

	// TODO: Retry all API Calls with backoff
	resp, err := api.Conditions.ListConditions(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool

	for _, s := range resp.GetPayload().Data {
		if searchName != "" && s.Attributes.Name == searchName {
			d.SetId(s.ID)
			d.Set("name", s.Attributes.Name)
			d.Set("stream_id", s.Relationships.Stream.Data.ID)
			d.Set("eval_window_ms", s.Attributes.Eval_window_ms)
			d.Set("expression", s.Attributes.Expression)
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("Unable to find a Condition with the name: %s", searchName)
	}

	return nil
}
