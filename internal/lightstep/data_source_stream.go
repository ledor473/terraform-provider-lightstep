package lightstep

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ledor473/lightstep-api-go/pkg/v0.2/client/streams"
)

func dataSourceStream() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStreamRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "query"},
			},
			"query": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "query"},
			},
		},
	}
}

func dataSourceStreamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Reading LightStep Streams")

	searchName := d.Get("name").(string)
	searchQuery := d.Get("query").(string)
	p := &streams.ListStreamsParams{Organization: config.Organization, Project: config.Project}

	// TODO: Retry all API Calls with backoff
	resp, err := api.Streams.ListStreams(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool

	for _, s := range resp.GetPayload().Data {
		if (searchName != "" && s.Attributes.Name == searchName) ||
			(searchQuery != "" && s.Attributes.Query == searchQuery) {
			d.SetId(s.ID)
			d.Set("name", s.Attributes.Name)
			d.Set("query", s.Attributes.Query)
			found = true
			break
		}
	}

	if !found && searchName != "" {
		return diag.Errorf("Unable to find a Stream with the name: %s", searchName)
	} else if !found && searchQuery != "" {
		return diag.Errorf("Unable to find a Stream with the query: %s", searchQuery)
	}

	return nil
}
