package sso

// Possible token types for Ark SIA SSO
const (
	Password          string = "password"
	ClientCertificate string = "client_certificate"
	OracleWallet      string = "oracle_wallet"
	RDPFile           string = "rdp_file"
)

// ArkSIASSOGetTokenInfo is a struct that represents the request for getting token information from the Ark SIA SSO service.
type ArkSIASSOGetTokenInfo struct {
	TokenType string `json:"token_type" validate:"required" mapstructure:"token_type" flag:"token-type" desc:"Which token type to get the info for [DPA-K8S, DPA-DB, DPA-RDP, DPA-SSH]" choices:"password,client_certificate,oracle_wallet,rdp_file"`
	Service   string `json:"service" validate:"required" mapstructure:"service" flag:"service" desc:"Which service to get the token info for [password, client_certificate, oracle_wallet, rdp_file]" choice:"DPA-DB,DPA-K8S,DPA-RDP,DPA-SSH"`
}

// ArkSIASSOTokenInfo is a struct that represents the response from the Ark SIA SSO service for token information.
type ArkSIASSOTokenInfo struct {
	Metadata map[string]interface{} `json:"metadata" validate:"required" mapstructure:"metadata"`
}
