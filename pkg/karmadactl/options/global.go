package options

import "github.com/spf13/pflag"

type GlobalOptions struct {
	Namespace  string
	Kubeconfig string
}

func NewGlobalOptions() *GlobalOptions {
	return &GlobalOptions{}
}

func (o *GlobalOptions) AddFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&o.Namespace, "namespace", "n", "default", "kubernetes name")
	flags.StringVar(&o.Kubeconfig, "kubeconfig", "", "--kubeconfig absolute path to the kubeconfig file")
}
