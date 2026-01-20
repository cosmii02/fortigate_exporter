package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestUserFirewall(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/user/firewall", "testdata/user-firewall.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeUserFirewall, c, r) {
		t.Errorf("probeUserFirewall() returned non-success")
	}

	em := `
	# HELP fortigate_user_firewall_info Info on authenticated firewall users
	# TYPE fortigate_user_firewall_info gauge
	fortigate_user_firewall_info{auth_server="LDAP1",group="Employees",src_ip="10.0.0.5",user="alice",vdom="root"} 1
	fortigate_user_firewall_info{auth_server="",group="Contractors",src_ip="10.0.0.6",user="bob",vdom="root"} 1
	fortigate_user_firewall_info{auth_server="",group="Employees",src_ip="10.0.0.7",user="dave",vdom="root"} 1
	fortigate_user_firewall_info{auth_server="",group="Employees",src_ip="10.0.1.7",user="carol",vdom="vdom2"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
