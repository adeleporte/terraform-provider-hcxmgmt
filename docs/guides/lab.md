---
page_title: "HCX Lab - Setup of the HCX Connector appliance"
---

This TF file automates the configuration of a HCX lab.
If you run this .tf just after OVA deployment, it manages creation/update/deletion of :
* Activation
* vCenter association
* SSO Configuration
* vSphere Role Mappings
* Location setup


## Usage example

```hcl

terraform {
  required_providers {
    hcxmgmt = {
      versions = ["0.1"]
      source = "vcn.cloud/edu/hcxmgmt"
    }
  }
}

provider hcxmgmt {
    hcx         = "https://hcx-connector-01a"
    username    = "admin"
    password    = "VMware1!"
}


resource "hcxmgmt_vcenter" "vcenter" {
    url         = "https://vcsa-01a.corp.local"
    username    = "administrator@vsphere.local"
    password    = "VMware1!"

    depends_on  = [hcxmgmt_activation.activation]
}

resource "hcxmgmt_sso" "sso" {
    vcenter     = hcxmgmt_vcenter.vcenter.id
    url         = "https://vcsa-01a.corp.local"
}


resource "hcxmgmt_activation" "activation" {
    activationkey = "*****-*****-*****-*****-*****"
}


resource "hcxmgmt_rolemapping" "rolemapping" {
    sso = hcxmgmt_sso.sso.id

    admin {
      user_group = "vsphere.local\\Administrators"
    }

    admin {
      user_group = "corp.local\\Administrators"
    }

    enterprise {
      user_group = "corp.local\\Administrators"
    }
}

resource "hcxmgmt_location" "location" {
    city        = "Paris"
    country     = "France"
    province    = "Ile-de-France"
    latitude    = 48.86669293
    longitude   = 2.333335326
}


```
