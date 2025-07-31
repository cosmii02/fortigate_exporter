package probe

import (
	"log"
	"strconv"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeVPNIPSec(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		status = prometheus.NewDesc(
			"fortigate_ipsec_tunnel_up",
			"Status of IPsec tunnel (0 - Down, 1 - Up)",
			[]string{"vdom", "name", "p2serial", "parent"}, nil,
		)
		transmitted = prometheus.NewDesc(
			"fortigate_ipsec_tunnel_transmit_bytes_total",
			"Total number of bytes transmitted over the IPsec tunnel",
			[]string{"vdom", "name", "p2serial", "parent"}, nil,
		)
		received = prometheus.NewDesc(
			"fortigate_ipsec_tunnel_receive_bytes_total",
			"Total number of bytes received over the IPsec tunnel",
			[]string{"vdom", "name", "p2serial", "parent"}, nil,
		)
		activeTunnels = prometheus.NewDesc(
			"fortigate_ipsec_tunnels_up",
			"Number of active IPsec tunnels",
			[]string{"vdom"}, nil,
		)
		totalTunnels = prometheus.NewDesc(
			"fortigate_ipsec_tunnels_total",
			"Total number of configured IPsec tunnels",
			[]string{"vdom"}, nil,
		)
		activeConnections = prometheus.NewDesc(
			"fortigate_ipsec_connections_up",
			"Number of active IPsec connections (proxy IDs)",
			[]string{"vdom"}, nil,
		)
		totalConnections = prometheus.NewDesc(
			"fortigate_ipsec_connections_total",
			"Total number of configured IPsec connections (proxy IDs)",
			[]string{"vdom"}, nil,
		)
	)

	type proxyid struct {
		Name     string  `json:"p2name"`
		P2serial int     `json:"p2serial"`
		Status   string  `json:"status"`
		Incoming float64 `json:"incoming_bytes"`
		Outgoing float64 `json:"outgoing_bytes"`
	}
	type tunnel struct {
		Name    string    `json:"name"`
		Type    string    `json:"type"`
		ProxyID []proxyid `json:"proxyid"`
	}
	type ipsecResult struct {
		Results []tunnel `json:"results"`
		VDOM    string
	}
	var res []ipsecResult
	if err := c.Get("api/v2/monitor/vpn/ipsec", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	// Maps to track statistics per VDOM
	vdomActiveTunnels := make(map[string]int)
	vdomTotalTunnels := make(map[string]int)
	vdomActiveConnections := make(map[string]int)
	vdomTotalConnections := make(map[string]int)

	for _, v := range res {
		for _, i := range v.Results {
			/*
			  type 'dialup' seems to be client vpn.
			  Not sure exactly what the difference is between probeVPNSsl
			*/
			if i.Type == "dialup" {
				continue
			}

			// Count tunnels
			vdomTotalTunnels[v.VDOM]++

			// Check if tunnel has any active proxy IDs
			hasActiveProxy := false
			for _, t := range i.ProxyID {
				// Count connections
				vdomTotalConnections[v.VDOM]++

				if t.Status == "up" {
					vdomActiveConnections[v.VDOM]++
					hasActiveProxy = true
				}

				s := 0.0
				if t.Status == "up" {
					s = 1.0
				}
				m = append(m, prometheus.MustNewConstMetric(status, prometheus.GaugeValue, s, v.VDOM, t.Name, strconv.Itoa(t.P2serial), i.Name))
				m = append(m, prometheus.MustNewConstMetric(transmitted, prometheus.CounterValue, t.Outgoing, v.VDOM, t.Name, strconv.Itoa(t.P2serial), i.Name))
				m = append(m, prometheus.MustNewConstMetric(received, prometheus.CounterValue, t.Incoming, v.VDOM, t.Name, strconv.Itoa(t.P2serial), i.Name))
			}

			// If tunnel has at least one active proxy ID, count it as active
			if hasActiveProxy {
				vdomActiveTunnels[v.VDOM]++
			}
		}
	}

	// Add VDOM-level statistics
	for vdom := range vdomTotalTunnels {
		m = append(m, prometheus.MustNewConstMetric(activeTunnels, prometheus.GaugeValue, float64(vdomActiveTunnels[vdom]), vdom))
		m = append(m, prometheus.MustNewConstMetric(totalTunnels, prometheus.GaugeValue, float64(vdomTotalTunnels[vdom]), vdom))
		m = append(m, prometheus.MustNewConstMetric(activeConnections, prometheus.GaugeValue, float64(vdomActiveConnections[vdom]), vdom))
		m = append(m, prometheus.MustNewConstMetric(totalConnections, prometheus.GaugeValue, float64(vdomTotalConnections[vdom]), vdom))
	}
	return m, true
}
