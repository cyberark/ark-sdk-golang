package k8s

// ArkSIAK8SGenerateKubeconfig is a struct that represents the request for generating a kubeconfig file from the Ark SIA K8S service.
type ArkSIAK8SGenerateKubeconfig struct {
	Folder string `json:"folder" mapstructure:"folder" flag:"folder" desc:"Output folder to write the kubeconfig to" default:"~/.kube"`
}
