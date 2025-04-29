package connections

// ArkConnectionCredentials represents the credentials for a connection.
type ArkConnectionCredentials struct {
	User               string `json:"user" mapstructure:"user"`
	Password           string `json:"password" mapstructure:"password"`
	PrivateKeyFilepath string `json:"private_key_filepath" mapstructure:"private_key_filepath"`
	PrivateKeyContents string `json:"private_key_contents" mapstructure:"private_key_contents"`
}
