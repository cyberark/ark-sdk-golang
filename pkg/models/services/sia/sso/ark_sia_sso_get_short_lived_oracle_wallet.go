package sso

// Possible wallet types for Ark SIA SSO
const (
	PEM string = "PEM"
	SSO string = "SSO"
)

// ArkSIASSOGetShortLivedOracleWallet is a struct that represents the request for getting a short-lived Oracle wallet from the Ark SIA SSO service.
type ArkSIASSOGetShortLivedOracleWallet struct {
	AllowCaching bool   `json:"allow_caching" mapstructure:"allow_caching" flag:"allow-caching" desc:"Allow short lived token caching" default:"false"`
	UnzipWallet  bool   `json:"unzip_wallet" mapstructure:"unzip_wallet" flag:"unzip-wallet" desc:"Whether to save zipped or not" default:"true"`
	Folder       string `json:"folder" validate:"required" mapstructure:"folder" flag:"folder" desc:"Output folder to write the wallet to"`
	WalletType   string `json:"wallet_type" mapstructure:"wallet_type" flag:"wallet-type" desc:"Type of wallet to generate, if PEM, no zip will be generated, only an ewallet.pem file" default:"SSO" choices:"PEM,SSO"`
}
