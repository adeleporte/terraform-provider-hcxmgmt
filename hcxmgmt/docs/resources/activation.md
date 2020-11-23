# activation

An activation key is mandatory to use a HCX system



## Example Usage

```hcl
resource "hcxmgmt_activation" "activation" {
    activationkey = "*****-*****-*****-*****-*****"
}

```

## Argument Reference

* `activationkey` - (Required) Activation key.


## Attribute Reference

* `id` - ID of the activation.
