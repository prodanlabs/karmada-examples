package dns

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unsafe"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"

	"github.com/prodanlabs/karmada-examples/pkg/util"
)

const uri = "/apis/cluster.karmada.io/v1alpha1/clusters/%s/proxy/api/v1/namespaces/%s/pods"

type domainName struct {
	ip       string
	hostname string
}

func (c *Controller) record(namespace, serviceName, labelSelector string, dn *[]domainName) error {
	clusterList, err := c.karmadaClient.ClusterV1alpha1().Clusters().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for i := range clusterList.Items {
		var podList corev1.PodList
		data, err := c.karmadaClient.ClusterV1alpha1().
			RESTClient().
			Get().
			RequestURI(fmt.Sprintf(uri, clusterList.Items[i].Name, namespace)).
			Param("labelSelector", labelSelector).
			//Timeout(60 * time.Second).
			DoRaw(context.TODO())
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &podList); err != nil {
			return err
		}

		for i := range podList.Items {
			*dn = append(*dn, domainName{
				hostname: fmt.Sprintf("%s.%s.%s.svc.cluster.local", podList.Items[i].Name, serviceName, namespace),
				ip:       podList.Items[i].Status.PodIP,
			})
		}
	}

	return nil
}

func (c *Controller) filter() ([]corev1.Service, error) {
	var compliantService []corev1.Service

	services, err := c.Clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for i := range services.Items {
		if services.Items[i].Spec.Type != corev1.ServiceTypeClusterIP && services.Items[i].Spec.ClusterIP != corev1.ClusterIPNone {
			break
		}

		annotations := services.Items[i].GetAnnotations()
		if v, ok := annotations["service.karmada.io/global"]; ok && v == "true" {
			compliantService = append(compliantService, services.Items[i])
		}
	}
	klog.V(6).Infof("compliantService: %v", compliantService)
	return compliantService, nil
}

func (c *Controller) aggregation() ([]domainName, error) {
	var dn []domainName

	services, err := c.filter()
	if err != nil {
		return nil, err
	}

	for i := range services {
		selector := services[i].Spec.Selector
		if err := c.record(services[i].Namespace, services[i].Name, util.MapToString(selector), &dn); err != nil {
			return nil, err
		}
	}

	if len(dn) == 0 {
		return nil, nil
	}

	return dn, nil
}

// lockState Get the state of the lock
func (c *Controller) lockState() error {
	state := (*uint32)(unsafe.Pointer(c.mu))
	if *state > 0 {
		return fmt.Errorf("locked")
	}

	return nil
}

func (c *Controller) addOrUpdateConfig() error {
	if err := c.lockState(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	dn, err := c.aggregation()
	if err != nil {
		return err
	}

	configMap, err := c.Clientset.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get(context.TODO(), "coredns", metav1.GetOptions{})
	if err != nil {
		return err
	}
	corefile := NewCorefile(configMap.Data["Corefile"])

	var updateCorefile []byte
	for i := range dn {
		updateCorefile = corefile.AddOrUpdate(dn[i].ip, dn[i].hostname)
	}

	if len(updateCorefile) == 0 {
		klog.V(6).Info("the Corefile config of ConfigMaps is empty")
		return nil
	}

	// update CoreDNS config
	configMap.Data["Corefile"] = strings.ReplaceAll(string(updateCorefile), "\t", "    ")
	klog.V(6).Infof("The new configuration of A after the update:\n", configMap.Data["Corefile"])
	if _, err = c.Clientset.CoreV1().ConfigMaps(metav1.NamespaceSystem).Update(context.TODO(), configMap, metav1.UpdateOptions{}); err != nil {
		return err
	}

	klog.Info("Corefile update complete.")
	return nil
}

func (c *Controller) deleteConfig(serviceName, namespace string) error {
	if err := c.lockState(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	configMap, err := c.Clientset.CoreV1().ConfigMaps(metav1.NamespaceSystem).Get(context.TODO(), "coredns", metav1.GetOptions{})
	if err != nil {
		return err
	}
	corefile := NewCorefile(configMap.Data["Corefile"])

	d := corefile.Delete(fmt.Sprintf("%s.%s", serviceName, namespace))

	// update CoreDNS config
	configMap.Data["Corefile"] = strings.ReplaceAll(string(d), "\t", "    ")
	klog.V(6).Infof("Corefile new configuration after deletion:\n", configMap.Data["Corefile"])
	if _, err = c.Clientset.CoreV1().ConfigMaps(metav1.NamespaceSystem).Update(context.TODO(), configMap, metav1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

func (c *Controller) Worker(interval time.Duration) {
	for {
		if err := c.addOrUpdateConfig(); err != nil {
			klog.Error(err)
		}

		time.Sleep(interval)
	}
}
