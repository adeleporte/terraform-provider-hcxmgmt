package main

import (
	"context"

	"github.com/adeleporte/terraform-provider-hcxmgmt/hcxmgmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRoleMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleMappingCreate,
		ReadContext:   resourceRoleMappingRead,
		UpdateContext: resourceRoleMappingUpdate,
		DeleteContext: resourceRoleMappingDelete,

		Schema: map[string]*schema.Schema{
			"admin": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_group": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "vsphere.local\\Administrators",
						},
					},
				},
			},
			"enterprise": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_group": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			"sso": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRoleMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return resourceRoleMappingUpdate(ctx, d, m)
}

func resourceRoleMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func resourceRoleMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*hcxmgmt.Client)

	admin := d.Get("admin").([]interface{})
	enterprise := d.Get("enterprise").([]interface{})

	admin_groups := []string{}
	for _, j := range admin {
		tmp := j.(map[string]interface{})
		admin_groups = append(admin_groups, tmp["user_group"].(string))
	}

	enterprise_groups := []string{}
	for _, j := range enterprise {
		tmp := j.(map[string]interface{})
		enterprise_groups = append(enterprise_groups, tmp["user_group"].(string))
	}

	body := []hcxmgmt.RoleMapping{
		hcxmgmt.RoleMapping{
			Role:       "System Administrator",
			UserGroups: admin_groups,
		},
		hcxmgmt.RoleMapping{
			Role:       "Enterprise Administrator",
			UserGroups: enterprise_groups,
		},
	}
	/*
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(body)
		return diag.Errorf("%s", buf)
	*/
	_, err := hcxmgmt.PutRoleMapping(client, body)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("role_mapping")

	return resourceRoleMappingRead(ctx, d, m)
}

func resourceRoleMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*hcxmgmt.Client)
	body := []hcxmgmt.RoleMapping{
		hcxmgmt.RoleMapping{
			Role:       "System Administrator",
			UserGroups: []string{},
		},
		hcxmgmt.RoleMapping{
			Role:       "Enterprise Administrator",
			UserGroups: []string{},
		},
	}

	_, err := hcxmgmt.PutRoleMapping(client, body)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
