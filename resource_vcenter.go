package main

import (
	"context"
	"time"

	b64 "encoding/base64"

	"github.com/adeleporte/terraform-provider-hcxmgmt/hcxmgmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcevCenter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcevCenterCreate,
		ReadContext:   resourcevCenterRead,
		UpdateContext: resourcevCenterUpdate,
		DeleteContext: resourcevCenterDelete,

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcevCenterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcxmgmt.Client)

	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	body := hcxmgmt.InsertvCenterBody{
		Data: hcxmgmt.InsertvCenterData{
			Items: []hcxmgmt.InsertvCenterDataItem{
				hcxmgmt.InsertvCenterDataItem{
					Config: hcxmgmt.InsertvCenterDataItemConfig{
						Username: username,
						Password: b64.StdEncoding.EncodeToString([]byte(password)),
						URL:      url,
					},
				},
			},
		},
	}

	res, err := hcxmgmt.InsertvCenter(client, body)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.InsertvCenterData.Items[0].Config.UUID)

	// Restart App Deamon
	hcxmgmt.AppEngineStop(client)

	// Wait for App Deamon to be stopped
	for {
		jr, err := hcxmgmt.GetAppEngineStatus(client)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.Result == "STOPPED" {
			break
		}
		time.Sleep(5 * time.Second)
	}

	hcxmgmt.AppEngineStart(client)

	// Wait for App Deamon to be stopped
	for {
		jr, err := hcxmgmt.GetAppEngineStatus(client)
		if err != nil {
			return diag.FromErr(err)
		}

		if jr.Result == "RUNNING" {
			break
		}
		time.Sleep(5 * time.Second)
	}

	return resourcevCenterRead(ctx, d, m)
}

func resourcevCenterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func resourcevCenterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourcevCenterRead(ctx, d, m)
}

func resourcevCenterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcxmgmt.Client)

	hcxmgmt.DeletevCenter(client, d.Id())

	return diags
}
