package lightstep

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ledor473/lightstep-api-go/pkg/v0.2/client/workflow_links"
	"github.com/ledor473/lightstep-api-go/pkg/v0.2/models"
)

func resourceWorkflowLink() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowLinkCreate,
		ReadContext:   resourceWorkflowLinkRead,
		UpdateContext: resourceWorkflowLinkUpdate,
		DeleteContext: resourceWorkflowLinkDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeString,
				Optional: true,
				// TODO: Define DiffSuppressFunc to ignore whitespace changes
				//DiffSuppressFunc: ,
			},
		},
	}
}

func resourceWorkflowLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating LightStep WorkflowLink")

	p := &workflow_links.CreateWorkflowLinkParams{Organization: config.Organization, Project: config.Project}
	rules, err := unmarshalRules(d.Get("rules").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	p.WithData(&models.ExternalLinkRequestBody{
		Data: &models.ExternalLinkRequest{
			Type: stringPointer("condition"),
			Attributes: &models.ExternalLinkAttributes{
				Name:  models.Name(d.Get("name").(string)),
				URL:   models.URL(d.Get("url").(string)),
				Rules: rules,
			},
		},
	})

	resp, err := api.WorkflowLinks.CreateWorkflowLink(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Payload.Data.ID)

	return nil
}

func resourceWorkflowLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading LightStep WorkflowLink")

	p := &workflow_links.GetWorkflowLinkParams{Organization: config.Organization, Project: config.Project, LinkID: d.Id()}

	resp, err := api.WorkflowLinks.GetWorkflowLink(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", resp.Payload.Data.Attributes.Name)
	d.Set("url", resp.Payload.Data.Attributes.URL)
	rules, err := marshalRules(resp.Payload.Data.Attributes.Rules)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("rules", rules)

	return nil
}

func resourceWorkflowLinkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating LightStep WorkflowLink")

	p := &workflow_links.PatchWorkflowLinkParams{Organization: config.Organization, Project: config.Project, LinkID: d.Id()}
	rules, err := unmarshalRules(d.Get("rules").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	p.WithData(&models.ExternalLinkRequestBody{
		Data: &models.ExternalLinkRequest{
			// TODO: API Docs mentions this field as required... unclear why
			Type: stringPointer("condition"),
			Attributes: &models.ExternalLinkAttributes{
				// TODO: Remove type alias on string if possible
				Name:  models.Name(d.Get("name").(string)),
				URL:   models.URL(d.Get("url").(string)),
				Rules: rules,
			},
		},
	})

	resp, err := api.WorkflowLinks.PatchWorkflowLink(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Payload.Data.ID)

	return resourceWorkflowLinkRead(ctx, d, meta)
}

func resourceWorkflowLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting LightStep WorkflowLink")

	p := &workflow_links.DeleteWorkflowLinkParams{Organization: config.Organization, Project: config.Project, LinkID: d.Id()}

	_, err = api.WorkflowLinks.DeleteWorkflowLink(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
