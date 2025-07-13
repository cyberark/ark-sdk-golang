package scans

// ArkSecHubScanIDs represents a list of scan IDs returned when triggering a scan.
type ArkSecHubScanIDs struct {
	ScanIDs []string `json:"scan_ids" mapstructure:"scan_ids" flag:"scan-ids" desc:"List of scan IDs" validate:"required,dive,required"`
}

// ArkSecHubTriggerScans represents the request structure for triggering scans in the Ark Secrets Hub.
type ArkSecHubTriggerScans struct {
	ID              string   `json:"id" mapstructure:"id" flag:"id" desc:"The ID of the scan, defaulted to default" default:"default"`
	Type            string   `json:"type" mapstructure:"type" flag:"type" desc:"The type of the scan (example: secret-store), defaulted to secret-store" default:"secret-store"`
	SecretStoresIds []string `json:"secret_stores_ids" mapstructure:"secret_stores_ids" flag:"secret-stores-ids" desc:"The stores to sync (pattern: store-{uuid-Format})"`
}

// ArkSecHubScanMap represents the request structure for mapping scans in the Ark Secrets Hub.
type ArkSecHubScanMap struct {
	Scope ArkSecHubSecretStoreIds `json:"scope" mapstructure:"scope" desc:"The scope of the secret store ids to scan"`
}

// ArkSecHubSecretStoreIds represents the structure for specifying secret store IDs in the Ark Secrets Hub.
type ArkSecHubSecretStoreIds struct {
	SecretStoresIds []string `json:"secret_stores_ids" mapstructure:"secret_stores_ids" flag:"secret-stores-ids" desc:"The stores to sync (pattern: store-{uuid-Format})"`
}
