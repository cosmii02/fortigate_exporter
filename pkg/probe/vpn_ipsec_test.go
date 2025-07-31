package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestVPNIPSec(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/vpn/ipsec", "testdata/ipsec.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeVPNIPSec, c, r) {
		t.Errorf("probeVPNIPSec() returned non-success")
	}

	em := `
	# HELP fortigate_ipsec_connections_total Total number of configured IPsec connections (proxy IDs)
	# TYPE fortigate_ipsec_connections_total gauge
	fortigate_ipsec_connections_total{vdom="root"} 2
	# HELP fortigate_ipsec_connections_up Number of active IPsec connections (proxy IDs)
	# TYPE fortigate_ipsec_connections_up gauge
	fortigate_ipsec_connections_up{vdom="root"} 1
	# HELP fortigate_ipsec_tunnel_receive_bytes_total Total number of bytes received over the IPsec tunnel
	# TYPE fortigate_ipsec_tunnel_receive_bytes_total counter
	fortigate_ipsec_tunnel_receive_bytes_total{name="tunnel_1-sub",p2serial="1",parent="tunnel_1",vdom="root"} 1.429824e+07
	fortigate_ipsec_tunnel_receive_bytes_total{name="tunnel_1-sub",p2serial="12",parent="tunnel_1",vdom="root"} 1.429824e+07
	# HELP fortigate_ipsec_tunnel_transmit_bytes_total Total number of bytes transmitted over the IPsec tunnel
	# TYPE fortigate_ipsec_tunnel_transmit_bytes_total counter
	fortigate_ipsec_tunnel_transmit_bytes_total{name="tunnel_1-sub",p2serial="1",parent="tunnel_1",vdom="root"} 1.424856e+07
	fortigate_ipsec_tunnel_transmit_bytes_total{name="tunnel_1-sub",p2serial="12",parent="tunnel_1",vdom="root"} 1.424856e+07
	# HELP fortigate_ipsec_tunnel_up Status of IPsec tunnel (0 - Down, 1 - Up)
	# TYPE fortigate_ipsec_tunnel_up gauge
	fortigate_ipsec_tunnel_up{name="tunnel_1-sub",p2serial="1",parent="tunnel_1",vdom="root"} 1
	fortigate_ipsec_tunnel_up{name="tunnel_1-sub",p2serial="12",parent="tunnel_1",vdom="root"} 0
	# HELP fortigate_ipsec_tunnels_total Total number of configured IPsec tunnels
	# TYPE fortigate_ipsec_tunnels_total gauge
	fortigate_ipsec_tunnels_total{vdom="root"} 1
	# HELP fortigate_ipsec_tunnels_up Number of active IPsec tunnels
	# TYPE fortigate_ipsec_tunnels_up gauge
	fortigate_ipsec_tunnels_up{vdom="root"} 1

	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestVPNIPSecWithCommonP2Names(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/vpn/ipsec", "testdata/ipsec-common-p2.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeVPNIPSec, c, r) {
		t.Errorf("probeVPNIPSec() returned non-success")
	}

	em := `
	# HELP fortigate_ipsec_connections_total Total number of configured IPsec connections (proxy IDs)
	# TYPE fortigate_ipsec_connections_total gauge
	fortigate_ipsec_connections_total{vdom="root"} 6
	# HELP fortigate_ipsec_connections_up Number of active IPsec connections (proxy IDs)
	# TYPE fortigate_ipsec_connections_up gauge
	fortigate_ipsec_connections_up{vdom="root"} 4
	# HELP fortigate_ipsec_tunnel_receive_bytes_total Total number of bytes received over the IPsec tunnel
	# TYPE fortigate_ipsec_tunnel_receive_bytes_total counter
	fortigate_ipsec_tunnel_receive_bytes_total{name="CommonP2",p2serial="22",parent="My VPN",vdom="root"} 0
	fortigate_ipsec_tunnel_receive_bytes_total{name="CommonP2",p2serial="23",parent="My VPN",vdom="root"} 4.782292004e+09
	fortigate_ipsec_tunnel_receive_bytes_total{name="CommonP2",p2serial="24",parent="My VPN",vdom="root"} 3.82868846e+08
	fortigate_ipsec_tunnel_receive_bytes_total{name="CommonP2",p2serial="25",parent="My VPN",vdom="root"} 1.581264e+06
	fortigate_ipsec_tunnel_receive_bytes_total{name="mgmt",p2serial="1",parent="My VPN",vdom="root"} 0
	fortigate_ipsec_tunnel_receive_bytes_total{name="some-network",p2serial="14",parent="My VPN",vdom="root"} 274832
	# HELP fortigate_ipsec_tunnel_transmit_bytes_total Total number of bytes transmitted over the IPsec tunnel
	# TYPE fortigate_ipsec_tunnel_transmit_bytes_total counter
	fortigate_ipsec_tunnel_transmit_bytes_total{name="CommonP2",p2serial="22",parent="My VPN",vdom="root"} 0
	fortigate_ipsec_tunnel_transmit_bytes_total{name="CommonP2",p2serial="23",parent="My VPN",vdom="root"} 1.57533239e+09
	fortigate_ipsec_tunnel_transmit_bytes_total{name="CommonP2",p2serial="24",parent="My VPN",vdom="root"} 5.53928639e+08
	fortigate_ipsec_tunnel_transmit_bytes_total{name="CommonP2",p2serial="25",parent="My VPN",vdom="root"} 3.1269542e+07
	fortigate_ipsec_tunnel_transmit_bytes_total{name="mgmt",p2serial="1",parent="My VPN",vdom="root"} 0
	fortigate_ipsec_tunnel_transmit_bytes_total{name="some-network",p2serial="14",parent="My VPN",vdom="root"} 112307
	# HELP fortigate_ipsec_tunnel_up Status of IPsec tunnel (0 - Down, 1 - Up)
	# TYPE fortigate_ipsec_tunnel_up gauge
	fortigate_ipsec_tunnel_up{name="CommonP2",p2serial="22",parent="My VPN",vdom="root"} 0
	fortigate_ipsec_tunnel_up{name="CommonP2",p2serial="23",parent="My VPN",vdom="root"} 1
	fortigate_ipsec_tunnel_up{name="CommonP2",p2serial="24",parent="My VPN",vdom="root"} 1
	fortigate_ipsec_tunnel_up{name="CommonP2",p2serial="25",parent="My VPN",vdom="root"} 1
	fortigate_ipsec_tunnel_up{name="mgmt",p2serial="1",parent="My VPN",vdom="root"} 0
	fortigate_ipsec_tunnel_up{name="some-network",p2serial="14",parent="My VPN",vdom="root"} 1
	# HELP fortigate_ipsec_tunnels_total Total number of configured IPsec tunnels
	# TYPE fortigate_ipsec_tunnels_total gauge
	fortigate_ipsec_tunnels_total{vdom="root"} 1
	# HELP fortigate_ipsec_tunnels_up Number of active IPsec tunnels
	# TYPE fortigate_ipsec_tunnels_up gauge
	fortigate_ipsec_tunnels_up{vdom="root"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestVPNIPSecWithDialup(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/vpn/ipsec", "testdata/ipsec-with-dialup.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeVPNIPSec, c, r) {
		t.Errorf("probeVPNIPSec() returned non-success")
	}

	em := `
	# HELP fortigate_ipsec_connections_total Total number of configured IPsec connections (proxy IDs)
	# TYPE fortigate_ipsec_connections_total gauge
	fortigate_ipsec_connections_total{vdom="root"} 3
	# HELP fortigate_ipsec_connections_up Number of active IPsec connections (proxy IDs)
	# TYPE fortigate_ipsec_connections_up gauge
	fortigate_ipsec_connections_up{vdom="root"} 3
	# HELP fortigate_ipsec_tunnel_receive_bytes_total Total number of bytes received over the IPsec tunnel
	# TYPE fortigate_ipsec_tunnel_receive_bytes_total counter
	fortigate_ipsec_tunnel_receive_bytes_total{name="nixvpn-split-client1",p2serial="101",parent="nixvpn-split",vdom="root"} 1.024e+06
	fortigate_ipsec_tunnel_receive_bytes_total{name="nixvpn-split-client2",p2serial="102",parent="nixvpn-split",vdom="root"} 512000
	fortigate_ipsec_tunnel_receive_bytes_total{name="tunnel_1-sub",p2serial="1",parent="tunnel_1",vdom="root"} 1.429824e+07
	# HELP fortigate_ipsec_tunnel_transmit_bytes_total Total number of bytes transmitted over the IPsec tunnel
	# TYPE fortigate_ipsec_tunnel_transmit_bytes_total counter
	fortigate_ipsec_tunnel_transmit_bytes_total{name="nixvpn-split-client1",p2serial="101",parent="nixvpn-split",vdom="root"} 2.048e+06
	fortigate_ipsec_tunnel_transmit_bytes_total{name="nixvpn-split-client2",p2serial="102",parent="nixvpn-split",vdom="root"} 1.024e+06
	fortigate_ipsec_tunnel_transmit_bytes_total{name="tunnel_1-sub",p2serial="1",parent="tunnel_1",vdom="root"} 1.424856e+07
	# HELP fortigate_ipsec_tunnel_up Status of IPsec tunnel (0 - Down, 1 - Up)
	# TYPE fortigate_ipsec_tunnel_up gauge
	fortigate_ipsec_tunnel_up{name="nixvpn-split-client1",p2serial="101",parent="nixvpn-split",vdom="root"} 1
	fortigate_ipsec_tunnel_up{name="nixvpn-split-client2",p2serial="102",parent="nixvpn-split",vdom="root"} 1
	fortigate_ipsec_tunnel_up{name="tunnel_1-sub",p2serial="1",parent="tunnel_1",vdom="root"} 1
	# HELP fortigate_ipsec_tunnels_total Total number of configured IPsec tunnels
	# TYPE fortigate_ipsec_tunnels_total gauge
	fortigate_ipsec_tunnels_total{vdom="root"} 2
	# HELP fortigate_ipsec_tunnels_up Number of active IPsec tunnels
	# TYPE fortigate_ipsec_tunnels_up gauge
	fortigate_ipsec_tunnels_up{vdom="root"} 2
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
