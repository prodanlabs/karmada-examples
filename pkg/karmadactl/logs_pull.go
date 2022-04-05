package karmadactl

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	grpccredentials "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	konnectivity "sigs.k8s.io/apiserver-network-proxy/konnectivity-client/pkg/client"

	"github.com/prodanlabs/karmada-examples/pkg/karmadactl/options"
	"github.com/prodanlabs/karmada-examples/pkg/util"
)

type LogsPullOptions struct {
	GlobalOptions   *options.GlobalOptions
	ProxyCACert     string
	ProxyCert       string
	ProxyKey        string
	CertDir         string
	ProxyServerHost string
	ProxyServerPort string
	PodName         string
	Follow          bool
	TailLines       int64
}

func (o *LogsPullOptions) AddAddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.ProxyCACert, "proxy-ca", "ca.crt", "anp frontend ca cert. path:certs/frontend/issued/ca.crt")
	flags.StringVar(&o.ProxyCert, "proxy-cert", "proxy-client.crt", "anp frontend proxy client cert. path:certs/frontend/issued/proxy-client.crt")
	flags.StringVar(&o.ProxyKey, "proxy-key", "proxy-client.key", "anp frontend proxy client key. path:certs/frontend/private/proxy-client.key")
	flags.StringVar(&o.CertDir, "cert-dir", "", "anp cert dir")
	flags.StringVar(&o.ProxyServerHost, "proxy-server-host", "127.0.0.1", "anp proxy server host")
	flags.StringVar(&o.ProxyServerPort, "proxy-server-port", "8090", "anp proxy server port")
	flags.BoolVarP(&o.Follow, "follow", "f", false, "follow logs")
	flags.Int64Var(&o.TailLines, "tail", -1, "follow logs")
}

func (o *LogsPullOptions) Complete(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("pod name is required")
	}
	o.PodName = args[0]

	if o.CertDir != "" && !strings.HasSuffix(o.CertDir, "/") {
		o.CertDir = fmt.Sprintf("%s/", o.CertDir)
	}
	o.ProxyCACert, o.ProxyCert, o.ProxyKey = o.CertDir+o.ProxyCACert, o.CertDir+o.ProxyCert, o.CertDir+o.ProxyKey

	return nil
}

// Validate checks the set of flags provided by the user
func (o *LogsPullOptions) Validate() error {
	if o.GlobalOptions.Kubeconfig == "" {
		return fmt.Errorf("absolute path to the kubeconfig file")
	}

	return nil
}

// CreateTunnel Create Grpc Tunnel
func (o *LogsPullOptions) CreateTunnel() (konnectivity.Tunnel, error) {
	tlsCfg, err := util.GetClientTLSConfig(o.ProxyCACert, o.ProxyCert, o.ProxyKey, o.ProxyServerHost, nil)
	if err != nil {
		return nil, err
	}

	return konnectivity.CreateSingleUseGrpcTunnel(
		context.TODO(),
		net.JoinHostPort(o.ProxyServerHost, o.ProxyServerPort),
		grpc.WithTransportCredentials(grpccredentials.NewTLS(tlsCfg)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: time.Second * 5,
		}),
	)
}

// GetPodLogs  get pod logs
func (o *LogsPullOptions) GetPodLogs(client *kubernetes.Clientset) error {
	podLogOpts := corev1.PodLogOptions{
		Follow: o.Follow,
	}
	if o.TailLines >= 0 {
		podLogOpts.TailLines = &o.TailLines
	}

	reqLogs := client.CoreV1().Pods(o.GlobalOptions.Namespace).GetLogs(o.PodName, &podLogOpts)
	podLogs, err := reqLogs.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer podLogs.Close()

	r := bufio.NewReader(podLogs)
	for {
		b, err := r.ReadBytes('\n')
		if err != nil {
			return err
		}
		fmt.Print(string(b))
	}
}

func (o *LogsPullOptions) Run() error {
	dialerTunnel, err := o.CreateTunnel()
	if err != nil {
		return err
	}

	cfg, err := util.RestConfig(false, o.GlobalOptions.Kubeconfig)
	if err != nil {
		return err
	}
	cfg.Dial = dialerTunnel.DialContext

	clientSet, err := util.NewClientSet(cfg)
	if err != nil {
		return err
	}

	if err := o.GetPodLogs(clientSet); err != nil {
		return err
	}

	return nil
}

func NewLogsPull(parentCommand string, opts *options.GlobalOptions) *cobra.Command {
	o := &LogsPullOptions{GlobalOptions: opts}
	cmd := &cobra.Command{
		Use:     "logs",
		Short:   "get pod logs of a member cluster in pull mode",
		Example: fmt.Sprintf("%s logs <POD_NAME>", parentCommand),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}

			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}
	o.AddAddFlags(cmd.Flags())
	return cmd
}
