package dns

import (
	"context"
	karmadaclientset "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sync"

	"github.com/prodanlabs/karmada-examples/pkg/util"
)

const (
	ControllerName = "dns-controller"
)

var _ reconcile.Reconciler = &Controller{}

// Controller reconciles a ContainerSet object
type Controller struct {
	client.Client
	recorder      record.EventRecorder
	Clientset     *kubernetes.Clientset
	karmadaClient karmadaclientset.Interface
	mu            *sync.Mutex
}

// Reconcile  The function does not differentiate between create, update or deletion events.
// Instead it simply reads the state of the cluster at the time it is called.
func (c *Controller) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	if err := c.deleteConfig(request.Name, request.Namespace); err != nil {
		klog.Errorf("Failed to remove obsolete DNS resolution. error: %v", err)
		return reconcile.Result{Requeue: true}, nil
	}

	klog.Info("Delete obsolete DNS resolution successfully.")

	return reconcile.Result{}, nil
}

func (c *Controller) SetupWithManager(mgr manager.Manager) error {
	predicate := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
	}

	return ctrl.NewControllerManagedBy(mgr).For(&corev1.Service{}).WithEventFilter(predicate).Complete(c)
}

// AddToManager create controller and register to controller manager
func (c *Controller) AddToManager(mgr manager.Manager) error {
	// Setup Scheme for k8s appv1 resources
	if err := corev1.AddToScheme(mgr.GetScheme()); err != nil {
		return err
	}

	return c.SetupWithManager(mgr)
}

// NewController returns a new Controller
func NewController(mgr manager.Manager) *Controller {
	c, err := util.NewClientSet(mgr.GetConfig())
	if err != nil {
		klog.Fatal(err)
	}

	return &Controller{
		Client:        mgr.GetClient(),
		recorder:      mgr.GetEventRecorderFor(ControllerName),
		karmadaClient: karmadaclientset.NewForConfigOrDie(mgr.GetConfig()),
		Clientset:     c,
		mu:            new(sync.Mutex),
	}
}
