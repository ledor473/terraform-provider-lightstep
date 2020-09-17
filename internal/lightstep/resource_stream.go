package lightstep

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ledor473/lightstep-api-go/pkg/v0.2/client/streams"
	"github.com/ledor473/lightstep-api-go/pkg/v0.2/models"
)

func resourceStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStreamCreate,
		ReadContext:   resourceStreamRead,
		UpdateContext: resourceStreamUpdate,
		DeleteContext: resourceStreamDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceStreamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating LightStep Stream")

	p := &streams.PostStreamParams{Organization: config.Organization, Project: config.Project}
	p.WithData(&models.CreateOrUpdateBody{
		Data: &models.CreateOrUpdateRequest{
			Type: stringPointer("stream"),
			Attributes: &models.StreamRequestAttributes{
				Name:  stringPointer(d.Get("name").(string)),
				Query: d.Get("query").(string),
			},
		},
	})

	resp, err := api.Streams.PostStream(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Payload.Data.ID)

	return nil
}

func resourceStreamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading LightStep Stream")

	p := &streams.GetStreamParams{Organization: config.Organization, Project: config.Project, StreamID: d.Id()}

	resp, err := api.Streams.GetStream(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", resp.Payload.Data.Attributes.Name)
	d.Set("query", resp.Payload.Data.Attributes.Query)

	return nil
}

func resourceStreamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating LightStep Stream")

	p := &streams.PatchStreamParams{Organization: config.Organization, Project: config.Project, StreamID: d.Id()}
	p.WithData(&models.CreateOrUpdateBody{
		Data: &models.CreateOrUpdateRequest{
			ID:   d.Id(),
			Type: stringPointer("stream"),
			Attributes: &models.StreamRequestAttributes{
				Name: stringPointer(d.Get("name").(string)),
			},
		},
	})

	resp, err := api.Streams.PatchStream(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Payload.Data.ID)

	return resourceStreamRead(ctx, d, meta)
}

func resourceStreamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting LightStep Stream")

	p := &streams.DeleteStreamParams{Organization: config.Organization, Project: config.Project, StreamID: d.Id()}

	_, err = api.Streams.DeleteStream(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
