package options

import (
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	componentbaseconfig "k8s.io/component-base/config"
)

type Options struct {
	BindAddress        string
	SecurePort         int
	LeaderElection     componentbaseconfig.LeaderElectionConfiguration
	DisableControllers string
	EnableControllers  string
	MetricsBindAddress string
	ResyncPeriod       metav1.Duration
}

// NewOptions builds an empty options.
func NewOptions() *Options {
	return &Options{
		LeaderElection: componentbaseconfig.LeaderElectionConfiguration{
			LeaderElect:  true,
			ResourceLock: resourcelock.LeasesResourceLock,
			ResourceName: "karmada-custom-controllers",
		},
	}
}

// AddFlags adds flags to the specified FlagSet.
func (o *Options) AddFlags(flags *pflag.FlagSet) {
	flags.Lookup("kubeconfig").Usage = "absolute path to the kubeconfig file"
	flags.StringVar(&o.BindAddress, "bind-address", "0.0.0.0", "The IP address on which to listen for the --secure-port port.")
	flags.IntVar(&o.SecurePort, "secure-port", 10258, "The secure port on which to serve HTTPS.")
	flags.BoolVar(&o.LeaderElection.LeaderElect, "leader-elect", true, "Enable leader election for controller manager.")
	flags.StringVarP(&o.LeaderElection.ResourceNamespace, "namespace", "n", "karmada-system", "Kubernetes namespace")
	flags.StringVar(&o.MetricsBindAddress, "metrics-bind-address", ":8080", "The address the metrics bind to.")
	flags.StringVar(&o.EnableControllers, "enable-controllers", "", "enabled controllers.")
	flags.StringVar(&o.DisableControllers, "disable-controllers", "", "disable controllers.")
	flags.DurationVar(&o.ResyncPeriod.Duration, "resync-period", 0, "informers resync period.")
}

func (o *Options) Validate() error {
	return nil
}

/*func (o *Options) Complete(args []string) error {
        return nil
}*/
