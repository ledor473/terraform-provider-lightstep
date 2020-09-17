package lightstep

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ledor473/lightstep-api-go/pkg/v0.2/client/conditions"
	"github.com/ledor473/lightstep-api-go/pkg/v0.2/models"
)

func resourceCondition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConditionCreate,
		ReadContext:   resourceConditionRead,
		UpdateContext: resourceConditionUpdate,
		DeleteContext: resourceConditionDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// TODO: Feels like expression should be required, but that's not what the API Doc indicate
			"expression": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"eval_window_ms": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			// TODO: We might need a ForceNew: True here
			"stream_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConditionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating LightStep Condition")

	p := &conditions.PostConditionParams{Organization: config.Organization, Project: config.Project}
	p.WithData(&models.ConditionRequestBody{
		Data: &models.ConditionRequest{
			Type: stringPointer("condition"),
			Attributes: &models.ConditionRequestAttributes{
				Name:               d.Get("name").(string),
				Expression:         d.Get("expression").(string),
				EvaluationWindowMs: int64(d.Get("eval_window_ms").(int)),
			},
			Relationships: &models.ConditionRequestRelationships{
				Stream: &models.RelatedResourceObject{
					ID: stringPointer(d.Get("stream_id").(string)),
				},
			},
		},
	})

	resp, err := api.Conditions.PostCondition(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Payload.Data.ID)

	return nil
}

func resourceConditionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading LightStep Condition")

	p := &conditions.GetConditionParams{Organization: config.Organization, Project: config.Project, ConditionID: d.Id()}

	resp, err := api.Conditions.GetCondition(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", resp.Payload.Data.Attributes.Name)
	d.Set("expression", resp.Payload.Data.Attributes.Expression)
	d.Set("eval_window_ms", resp.Payload.Data.Attributes.Eval_window_ms)

	return nil
}

func resourceConditionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating LightStep Condition")

	p := &conditions.PatchConditionParams{Organization: config.Organization, Project: config.Project, ConditionID: d.Id()}
	p.WithData(&models.ConditionRequestBody{
		Data: &models.ConditionRequest{
			Type: stringPointer("condition"),
			Attributes: &models.ConditionRequestAttributes{
				Name:               d.Get("name").(string),
				Expression:         d.Get("expression").(string),
				EvaluationWindowMs: int64(d.Get("eval_window_ms").(int)),
			},
			Relationships: &models.ConditionRequestRelationships{
				Stream: &models.RelatedResourceObject{
					ID: stringPointer(d.Get("stream_id").(string)),
				},
			},
		},
	})

	resp, err := api.Conditions.PatchCondition(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Payload.Data.ID)

	return resourceConditionRead(ctx, d, meta)
}

func resourceConditionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	api, err := config.Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting LightStep Condition")

	p := &conditions.DeleteConditionParams{Organization: config.Organization, Project: config.Project, ConditionID: d.Id()}

	_, err = api.Conditions.DeleteCondition(ctx, p)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
