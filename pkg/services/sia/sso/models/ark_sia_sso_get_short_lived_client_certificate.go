package models

// Possible output formats for Ark SIA SSO
const (
	File       string = "file"
	SingleFile string = "single_file"
	Raw        string = "raw"
	Base64     string = "base64"
)

// ArkSIASSOGetShortLivedClientCertificate is a struct that represents the request for getting a short-lived client certificate from the Ark SIA SSO service.
type ArkSIASSOGetShortLivedClientCertificate struct {
	AllowCaching bool   `json:"allow_caching" mapstructure:"allow_caching" flag:"allow-caching" desc:"Allow short lived token caching" default:"false"`
	Folder       string `json:"folder" validate:"required" mapstructure:"folder" flag:"folder" desc:"Output folder to write the key / certificate to. Required if format is File"`
	OutputFormat string `json:"output_format" mapstructure:"output_format" flag:"output-format" desc:"The output format of the key / ' 'certificate. i.e. File, Raw, Base64" default:"file" choices:"file,single_file,raw,base64"`
	Service      string `json:"service" validate:"required" mapstructure:"service" flag:"service" desc:"Which service to generate the short lived certificate for - DPA-DB, DPA-K8S" choices:"DPA-DB,DPA-K8S"`
}
