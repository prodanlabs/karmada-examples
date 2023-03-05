package dns

import (
	"testing"

	"k8s.io/klog/v2"
)

func TestNewCorefile(t *testing.T) {
	config := `.:53 {
        errors
        health {
                lameduck 5s
        }
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
                pods insecure
                fallthrough in-addr.arpa ip6.arpa
                ttl 30
        }

        hosts {
                0.0.0.0 nginx-0.nginx-headless.default.svc.cluster.local
                2.2.2.2 nginx-2.nginx-headless.default.svc.cluster.local
                3.3.3.3 nginx-3.nginx-headless.default.svc.cluster.local
                fallthrough
        }
        prometheus :9153
        forward . /etc/resolv.conf {
                max_concurrent 1000
        }
        cache 30
        loop
        reload
        loadbalance
}`
	c := NewCorefile(config)

	if !c.isExists("nginx-1.nginx-headless.default.svc.cluster.local") {
		klog.Errorln("nginx-1.nginx-headless.default.svc.cluster.local not Exists")
	}
	add := c.AddOrUpdate("1.1.1.1", "nginx-1.nginx-headless.default.svc.cluster.local")
	klog.Infof("Add: \n %s", add)

	update := c.AddOrUpdate("22.22.22.22", "nginx-2.nginx-headless.default.svc.cluster.local")
	klog.Infof("Update: \n %s", update)

	del := c.Delete("nginx-3.nginx-headless.default.svc.cluster.local")
	klog.Infof("Delete: \n %s", del)
}
