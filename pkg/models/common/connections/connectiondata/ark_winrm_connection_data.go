package connectiondata

// ArkWinRMConnectionData represents the connection data for a WinRM connection.
type ArkWinRMConnectionData struct {
	CertificatePath  string `json:"certificate_path" mapstructure:"certificate_path"`
	TrustCertificate bool   `json:"trust_certificate" mapstructure:"trust_certificate"`
}
