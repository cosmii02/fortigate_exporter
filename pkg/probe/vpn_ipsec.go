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
		phase1Status = prometheus.NewDesc(
			"fortigate_ipsec_phase1_up",
			"Status of IPsec phase 1 tunnel (0 - Down, 1 - Up)",
			[]string{"vdom", "name", "type"}, nil,
		)
		phase2Status = prometheus.NewDesc(
			"fortigate_ipsec_phase2_up",
			"Status of IPsec phase 2 selector aggregated by parent/name (0 - Down, 1 - Up)",
			[]string{"vdom", "parent", "name"}, nil,
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
		tunnelActiveConnections = prometheus.NewDesc(
			"fortigate_ipsec_tunnel_connections_up",
			"Number of active client connections for this IPsec tunnel",
			[]string{"vdom", "tunnel", "type"}, nil,
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

	// Maps to track per-tunnel connection counts
	tunnelActiveConns := make(map[string]map[string]int) // vdom -> tunnel -> count
	tunnelTypes := make(map[string]map[string]string)    // vdom -> tunnel -> type

	for _, v := range res {
		// Initialize per-vdom tunnel connection map
		if tunnelActiveConns[v.VDOM] == nil {
			tunnelActiveConns[v.VDOM] = make(map[string]int)
		}
		if tunnelTypes[v.VDOM] == nil {
			tunnelTypes[v.VDOM] = make(map[string]string)
		}

		for _, i := range v.Results {
			// Count tunnels (now including dialup/client VPNs)
			vdomTotalTunnels[v.VDOM]++

			// Initialize connection count for this tunnel
			tunnelActiveConns[v.VDOM][i.Name] = 0
			tunnelTypes[v.VDOM][i.Name] = i.Type

			// Check if tunnel has any active proxy IDs
			hasActiveProxy := false
			p2StatusByName := make(map[string]float64)
			for _, t := range i.ProxyID {
				// Count connections
				vdomTotalConnections[v.VDOM]++

				if t.Status == "up" {
					vdomActiveConnections[v.VDOM]++
					tunnelActiveConns[v.VDOM][i.Name]++
					hasActiveProxy = true
				}

				s := 0.0
				if t.Status == "up" {
					s = 1.0
				}
				if p2StatusByName[t.Name] < s {
					p2StatusByName[t.Name] = s
				}
				m = append(m, prometheus.MustNewConstMetric(status, prometheus.GaugeValue, s, v.VDOM, t.Name, strconv.Itoa(t.P2serial), i.Name))
				m = append(m, prometheus.MustNewConstMetric(transmitted, prometheus.CounterValue, t.Outgoing, v.VDOM, t.Name, strconv.Itoa(t.P2serial), i.Name))
				m = append(m, prometheus.MustNewConstMetric(received, prometheus.CounterValue, t.Incoming, v.VDOM, t.Name, strconv.Itoa(t.P2serial), i.Name))
			}

			phase1Up := 0.0
			// If tunnel has at least one active proxy ID, count it as active
			if hasActiveProxy {
				phase1Up = 1.0
				vdomActiveTunnels[v.VDOM]++
			}
			m = append(m, prometheus.MustNewConstMetric(phase1Status, prometheus.GaugeValue, phase1Up, v.VDOM, i.Name, i.Type))
			for p2Name, p2Up := range p2StatusByName {
				m = append(m, prometheus.MustNewConstMetric(phase2Status, prometheus.GaugeValue, p2Up, v.VDOM, i.Name, p2Name))
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

	// Add per-tunnel connection counts
	for vdom, tunnels := range tunnelActiveConns {
		for tunnelName, activeCount := range tunnels {
			tunnelType := tunnelTypes[vdom][tunnelName]
			if tunnelType == "" {
				tunnelType = "automatic"
			}
			m = append(m, prometheus.MustNewConstMetric(tunnelActiveConnections, prometheus.GaugeValue, float64(activeCount), vdom, tunnelName, tunnelType))
		}
	}
	return m, true
}
