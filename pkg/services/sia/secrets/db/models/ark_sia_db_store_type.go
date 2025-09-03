package models

// Possible store types
const (
	Managed string = "managed"
	PAM     string = "pam"
)

// SecretTypeToStoreDict maps secret types to store types.
var SecretTypeToStoreDict = map[string]string{
	UsernamePassword: Managed,
	CyberArkPAM:      PAM,
	IAMUser:          Managed,
	AtlasAccessKeys:  Managed,
}
