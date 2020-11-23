package main

import (
	"context"

	"github.com/adeleporte/terraform-provider-hcxmgmt/hcxmgmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSSO() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSOCreate,
		ReadContext:   resourceSSORead,
		UpdateContext: resourceSSOUpdate,
		DeleteContext: resourceSSODelete,

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vcenter": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSSOCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcxmgmt.Client)

	url := d.Get("url").(string)

	body := hcxmgmt.InsertSSOBody{
		Data: hcxmgmt.InsertSSOData{
			Items: []hcxmgmt.InsertSSODataItem{
				hcxmgmt.InsertSSODataItem{
					Config: hcxmgmt.InsertSSODataItemConfig{
						LookupServiceUrl: url,
						ProviderType:     "PSC",
					},
				},
			},
		},
	}

	// First, check if SSO config is already present
	res, err := hcxmgmt.GetSSO(client)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(res.InsertSSOData.Items) == 0 {
		// No SSO config found
		res, err := hcxmgmt.InsertSSO(client, body)

		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(res.InsertSSOData.Items[0].Config.UUID)
		return resourceSSORead(ctx, d, m)
	}

	// Update existing SSO
	d.SetId(res.InsertSSOData.Items[0].Config.UUID)
	return resourceSSOUpdate(ctx, d, m)

}

func resourceSSORead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func resourceSSOUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcxmgmt.Client)

	url := d.Get("url").(string)

	body := hcxmgmt.InsertSSOBody{
		Data: hcxmgmt.InsertSSOData{
			Items: []hcxmgmt.InsertSSODataItem{
				hcxmgmt.InsertSSODataItem{
					Config: hcxmgmt.InsertSSODataItemConfig{
						LookupServiceUrl: url,
						UUID:             d.Id(),
						ProviderType:     "PSC",
					},
				},
			},
		},
	}

	_, err := hcxmgmt.UpdateSSO(client, body)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSSORead(ctx, d, m)
}

func resourceSSODelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcxmgmt.Client)

	hcxmgmt.DeleteSSO(client, d.Id())

	return diags
}
