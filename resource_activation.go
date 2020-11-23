package main

import (
	"context"

	"github.com/adeleporte/terraform-provider-hcxmgmt/hcxmgmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceActivation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActivationCreate,
		ReadContext:   resourceActivationRead,
		UpdateContext: resourceActivationUpdate,
		DeleteContext: resourceActivationDelete,

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://connect.hcx.vmware.com",
			},
			"activationkey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceActivationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcxmgmt.Client)

	url := d.Get("url").(string)
	activationkey := d.Get("activationkey").(string)

	body := hcxmgmt.ActivateBody{
		Data: hcxmgmt.ActivateData{
			Items: []hcxmgmt.ActivateDataItem{
				hcxmgmt.ActivateDataItem{
					Config: hcxmgmt.ActivateDataItemConfig{
						URL:           url,
						ActivationKey: activationkey,
					},
				},
			},
		},
	}

	// First, check if already activated
	res, err := hcxmgmt.GetActivate(client)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(res.Data.Items) == 0 {
		// No activation config found
		_, err := hcxmgmt.PostActivate(client, body)

		if err != nil {
			return diag.FromErr(err)
		}

		return resourceActivationRead(ctx, d, m)
	}

	d.SetId(res.Data.Items[0].Config.UUID)

	return resourceActivationRead(ctx, d, m)
}

func resourceActivationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcxmgmt.Client)

	res, err := hcxmgmt.GetActivate(client)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.Data.Items[0].Config.UUID)

	return diags
}

func resourceActivationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceActivationRead(ctx, d, m)
}

func resourceActivationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	/*
		client := m.(*hcxmgmt.Client)

		url := d.Get("url").(string)
		activationkey := d.Get("activationkey").(string)

		body := hcxmgmt.ActivateBody{
			Data: hcxmgmt.ActivateData{
				Items: []hcxmgmt.ActivateDataItem{
					hcxmgmt.ActivateDataItem{
						Config: hcxmgmt.ActivateDataItemConfig{
							URL:           url,
							ActivationKey: activationkey,
							UUID:          d.Id(),
						},
					},
				},
			},
		}

		_, err := hcxmgmt.DeleteActivate(client, body)
		if err != nil {
			return diag.FromErr(err)
		}
	*/
	return diags
}
