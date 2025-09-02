package models

// ArkSecHubSetConfiguration represnets the response when updating configuraiton settings.
type ArkSecHubSetConfiguration struct {
	SyncSettings ArkSecHubSyncSettings `json:"sync_settings" mapstructure:"sync_settings" desc:"Sync Settings for Secrets Hub" flag:"sync-settings" validate:"required"`
}
