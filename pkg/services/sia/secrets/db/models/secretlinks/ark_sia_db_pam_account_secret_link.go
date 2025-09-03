package secretlinks

// ArkSIADBPAMAccountSecretLink represents the link to a PAM account secret in the Ark SIA DB.
type ArkSIADBPAMAccountSecretLink struct {
	ArkSIADBSecretLink
	Safe        string `json:"safe,omitempty" mapstructure:"safe" desc:"Safe of the account"`
	AccountName string `json:"account_name,omitempty" mapstructure:"account_name" desc:"Account name"`
}

// GetLinkSecretType returns the secret type of the secret link.
func (s *ArkSIADBPAMAccountSecretLink) GetLinkSecretType() string {
	return "cyberark_pam"
}
