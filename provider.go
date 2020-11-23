package main

import (
	"context"

	"github.com/adeleporte/terraform-provider-hcxmgmt/hcxmgmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hcx": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HCX_URL", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HCX_MGMT_USER", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HCX_MGMT_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hcxmgmt_vcenter":     resourcevCenter(),
			"hcxmgmt_sso":         resourceSSO(),
			"hcxmgmt_activation":  resourceActivation(),
			"hcxmgmt_rolemapping": resourceRoleMapping(),
			"hcxmgmt_location":    resourceLocation(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	hcxurl := d.Get("hcx").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	if hcxurl != "" {
		c, err := hcxmgmt.NewClient(&hcxurl, &username, &password)
		//c := &http.Client{Timeout: 10 * time.Second}

		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, diags
	}

	return nil, diag.Errorf("Missing credentials")
}
