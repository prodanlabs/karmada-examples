package app

import (
	"context"
	"flag"
	"fmt"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/prodanlabs/karmada-examples/cmd/custom-webhook/app/options"
	"github.com/prodanlabs/karmada-examples/pkg/webhook/namespace"
)

const (
	CheckEndpointHealthz = "healthz"
	CheckEndpointReadyz  = "readyz"
)

var aggregatedScheme = runtime.NewScheme()

func init() {
	var _ = corev1.AddToScheme(aggregatedScheme)
	var _ = scheme.AddToScheme(aggregatedScheme)
	var _ = clusterv1alpha1.AddToScheme(aggregatedScheme)
}

// NewWebhookCommand creates a *cobra.Command object with default parameters
func NewWebhookCommand(ctx context.Context) *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:  "karmada-webhook",
		Long: `Start a karmada webhook server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Run(ctx, opts); err != nil {
				return err
			}
			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	klog.InitFlags(flag.CommandLine)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	opts.AddFlags(cmd.Flags())
	return cmd
}

func NewSchema() *runtime.Scheme {
	return aggregatedScheme
}

// Run runs the webhook server with options. This should never exit.
func Run(ctx context.Context, opts *options.Options) error {
	config, err := ctrl.GetConfig()
	if err != nil {
		return err
	}

	hookManager, err := ctrl.NewManager(config, ctrl.Options{
		Scheme: NewSchema(),
		WebhookServer: &webhook.Server{
			Host:          opts.BindAddress,
			Port:          opts.SecurePort,
			CertDir:       opts.CertDir,
			CertName:      opts.CertName,
			KeyName:       opts.KeyName,
			TLSMinVersion: opts.TLSMinVersion,
		},
		LeaderElection:         false,
		MetricsBindAddress:     opts.MetricsBindAddress,
		HealthProbeBindAddress: opts.HealthProbeBindAddress,
	})
	if err != nil {
		klog.Errorf("failed to build webhook server: %v", err)
		return err
	}
	if err := hookManager.AddHealthzCheck(CheckEndpointHealthz, healthz.Ping); err != nil {
		return fmt.Errorf("failed to add %q health check endpoint: %v", CheckEndpointHealthz, err)
	}
	if err := hookManager.AddReadyzCheck(CheckEndpointReadyz, healthz.Ping); err != nil {
		return fmt.Errorf("failed to add %q health check endpoint: %v", CheckEndpointReadyz, err)
	}

	klog.Info("registering webhooks to the webhook server")
	hookServer := hookManager.GetWebhookServer()

	hookServer.Register("/validate-namespace", &webhook.Admission{Handler: namespace.NewValidatingAdmission(hookManager)})

	// blocks until the context is done.
	if err := hookManager.Start(ctx); err != nil {
		klog.Errorf("webhook server exits unexpectedly: %v", err)
		return err
	}

	// never reach here
	return nil
}
