package namespace

import (
	"context"
	"encoding/json"
	"fmt"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/prodanlabs/karmada-examples/pkg/util"
)

// ValidatingAdmission validates ClusterOverridePolicy object when creating/updating/deleting.
type ValidatingAdmission struct {
	client.Client
	Clientset *kubernetes.Clientset
	Config    *rest.Config
	decoder   *admission.Decoder
}

func NewValidatingAdmission(mgr manager.Manager) *ValidatingAdmission {
	clientset, err := util.NewClientSet(mgr.GetConfig())
	if err != nil {
		klog.Fatal(err)
	}
	return &ValidatingAdmission{
		Client:    mgr.GetClient(),
		Clientset: clientset,
		Config:    mgr.GetConfig(),
	}
}

// Check if our ValidatingAdmission implements necessary interface
var _ admission.Handler = &ValidatingAdmission{}

// Handle implements admission.Handler interface.
// It yields a response to an AdmissionRequest.
func (v *ValidatingAdmission) Handle(ctx context.Context, req admission.Request) admission.Response {
	klog.Infof("%s namespace %q", req.Operation, req.Name)
	clusters, err := v.ClusterList(ctx)
	if err != nil {
		return admission.Denied(err.Error())
	}

	for _, c := range clusters {
		if err := v.PodList(c, req.Name); err != nil {
			return admission.Denied(err.Error())
		}
	}

	return admission.Allowed("")
}

func (v *ValidatingAdmission) ClusterList(ctx context.Context) ([]string, error) {
	clusterList := &clusterv1alpha1.ClusterList{}
	if err := v.Client.List(ctx, clusterList); err != nil {
		return nil, fmt.Errorf("failed to list clusters, error: %v", err)
	}

	var clusters []string
	for _, c := range clusterList.Items {
		// TODO Pull 模式可以部署 apiserver-network-proxy（ANP）来访问。这里跳过.
		// https://github.com/karmada-io/karmada/blob/master/docs/userguide/aggregated-api-endpoint.md
		if c.Spec.SyncMode == clusterv1alpha1.Pull {
			continue
		}
		clusters = append(clusters, c.Name)
	}
	return clusters, nil
}

func (v *ValidatingAdmission) PodList(cluster, namespace string) error {
	uil := fmt.Sprintf("%s/apis/cluster.karmada.io/v1alpha1/clusters/%s/proxy/api/v1/namespaces/%s/pods", v.Config.Host, cluster, namespace)
	request := v.Clientset.RESTClient().Get().RequestURI(uil)
	result := request.Do(context.Background())
	if err := result.Error(); err != nil {
		return fmt.Errorf("result err: %v", err)
	}
	pods, err := result.Raw()
	if err != nil {
		return fmt.Errorf("raw err: %v", err)
	}
	podList := &corev1.PodList{}

	if err := json.Unmarshal(pods, podList); err != nil {
		return fmt.Errorf("unmarshal err: %v", err)
	}

	if len(podList.Items) == 0 {
		return nil
	}

	var activityPods []string
	for _, p := range podList.Items {
		fmt.Println(p.Name)
		activityPods = append(activityPods, p.Name)
	}

	return fmt.Errorf("the workload of cluster %q namespace %q is not empty, pods: %s", cluster, namespace, activityPods)
}
