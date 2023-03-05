package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/prodanlabs/karmada-examples/cmd/custom-controller-manager/app/options"
	"github.com/prodanlabs/karmada-examples/pkg/controllers/deployment"
	"github.com/prodanlabs/karmada-examples/pkg/controllers/dns"
	"github.com/prodanlabs/karmada-examples/pkg/util"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"net"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"strconv"
	"time"
)

const (
	CheckEndpointHealthz = "healthz"
	CheckEndpointReadyz  = "readyz"
)

func NewCustomControllerManagerCommand(ctx context.Context) *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:  "karmada-custom-controller-manager",
		Long: `karmada custom controller manager`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			return Run(ctx, opts)
		},
	}

	klog.InitFlags(flag.CommandLine)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	opts.AddFlags(cmd.Flags())
	return cmd
}

func Run(ctx context.Context, opts *options.Options) error {
	config, err := ctrl.GetConfig()
	if err != nil {
		return err
	}
	util.SetupKubeConfig(config)

	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme:                     runtime.NewScheme(),
		SyncPeriod:                 &opts.ResyncPeriod.Duration,
		LeaderElection:             opts.LeaderElection.LeaderElect,
		LeaderElectionID:           opts.LeaderElection.ResourceName,
		LeaderElectionNamespace:    opts.LeaderElection.ResourceNamespace,
		LeaderElectionResourceLock: opts.LeaderElection.ResourceLock,
		HealthProbeBindAddress:     net.JoinHostPort(opts.BindAddress, strconv.Itoa(opts.SecurePort)),
		MetricsBindAddress:         opts.MetricsBindAddress,
	})
	if err != nil {
		return fmt.Errorf("new controller manager failed: %v", err)
	}

	if err := mgr.AddHealthzCheck(CheckEndpointHealthz, healthz.Ping); err != nil {
		return fmt.Errorf("failed to add %q health check endpoint: %v", CheckEndpointHealthz, err)
	}
	if err := mgr.AddReadyzCheck(CheckEndpointReadyz, healthz.Ping); err != nil {
		return fmt.Errorf("failed to add %q health check endpoint: %v", CheckEndpointReadyz, err)
	}

	if err := deployment.AddToManager(mgr); err != nil {
		return err
	}

	dnsController := dns.NewController(mgr)
	if err := dnsController.AddToManager(mgr); err != nil {
		return err
	}
	go dnsController.Worker(5 * time.Second)

	if err := mgr.Start(ctx); err != nil {
		return fmt.Errorf("controller manager exit: %v", err)
	}

	return nil
}
