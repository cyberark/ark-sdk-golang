---
title: Schemas
description: Schemas
---

# Schemas

Ark SDK is entirely based on schemas constructed from standard Golang structs, along with tagging of "json" and "mapstructure" for serialization.

All `exec` actions in the Ark SDK receive a model parsed from the CLI or from the SDK in code and, some of them, return a model or set of models.

## Example

Any request can be called with a defined model, for example:

```go
secret, err := siaAPI.SecretsVM().AddSecret(
    &vmsecretsmodels.ArkSIAVMAddSecret{
        SecretType:          "ProvisionerUser",
        ProvisionerUsername: "CoolUser",
        ProvisionerPassword: "CoolPassword",
    },
)
```

The above example creates a VM secret service and calls `AddSecret()` to add a new VM secret. and the relevant add secret schema is passed. finally, a result schema for a secret is returned:

```go
// ArkSIAVMSecret represents a secret in the Ark SIA VM.
type ArkSIAVMSecret struct {
	SecretID      string                 `json:"secret_id" mapstructure:"secret_id" flag:"secret-id" desc:"ID of the secret"`
	TenantID      string                 `json:"tenant_id,omitempty" mapstructure:"tenant_id,omitempty" flag:"tenant-id" desc:"Tenant ID of the secret"`
	Secret        ArkSIAVMSecretData     `json:"secret,omitempty" mapstructure:"secret,omitempty" flag:"secret" desc:"Secret itself"`
	SecretType    string                 `json:"secret_type" mapstructure:"secret_type" flag:"secret-type" desc:"Type of the secret" choices:"ProvisionerUser,PCloudAccount"`
	SecretDetails map[string]interface{} `json:"secret_details" mapstructure:"secret_details" flag:"secret-details" desc:"Secret extra details"`
	IsActive      bool                   `json:"is_active" mapstructure:"is_active" flag:"is-active" desc:"Whether this secret is active or not and can be retrieved or modified"`
	IsRotatable   bool                   `json:"is_rotatable" mapstructure:"is_rotatable" flag:"is-rotatable" desc:"Whether this secret can be rotated"`
	CreationTime  string                 `json:"creation_time" mapstructure:"creation_time" flag:"creation-time" desc:"Creation time of the secret"`
	LastModified  string                 `json:"last_modified" mapstructure:"last_modified" flag:"last-modified" desc:"Last time the secret was modified"`
	SecretName    string                 `json:"secret_name,omitempty" mapstructure:"secret_name,omitempty" flag:"secret-name" desc:"A friendly name label"`
}
```

All models can be found [here](https://github.com/cyberark/ark-sdk-golang/tree/main/ark_sdk_golang/pkg/models) and are separated to folders based on topic, from auth to services
