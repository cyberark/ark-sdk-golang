package actions

import ssomodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/sso/models"

// ActionToSchemaMap is a map that defines the mapping between SSO action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"short-lived-password":           &ssomodels.ArkSIASSOGetShortLivedPassword{},
	"short-lived-client-certificate": &ssomodels.ArkSIASSOGetShortLivedClientCertificate{},
	"short-lived-oracle-wallet":      &ssomodels.ArkSIASSOGetShortLivedOracleWallet{},
	"short-lived-rdp-file":           &ssomodels.ArkSIASSOGetShortLivedRDPFile{},
	"short-lived-token-info":         &ssomodels.ArkSIASSOGetTokenInfo{},
	"short-lived-ssh-key":            &ssomodels.ArkSIASSOGetSSHKey{},
}
