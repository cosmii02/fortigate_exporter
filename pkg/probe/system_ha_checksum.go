package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type HAChecksumResults struct {
	IsManageMaster int    `json:"is_manage_master"`
	IsRootMaster   int    `json:"is_root_master"`
	SerialNo       string `json:"serial_no"`
}

type HAChecksum struct {
	Results []HAChecksumResults `json:"results"`
	VDOM    string              `json:"vdom"`
}

func probeSystemHAChecksum(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		isMaster = prometheus.NewDesc(
			"fortigate_ha_member_has_role",
			"Master/Slave information",
			[]string{"role", "serial", "hostname", "vdom"}, nil,
		)
	)

	var res HAChecksum
	if err := c.Get("api/v2/monitor/system/ha-checksums", "scope=global", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	type haStatisticsResult struct {
		Hostname string `json:"hostname"`
		SerialNo string `json:"serial_no"`
	}
	type haStatisticsResponse struct {
		Results []haStatisticsResult `json:"results"`
		VDOM    string               `json:"vdom"`
	}

	hostnames := map[string]string{}
	vdom := res.VDOM

	var stats haStatisticsResponse
	if err := c.Get("api/v2/monitor/system/ha-statistics", "", &stats); err != nil {
		log.Printf("Warning: failed to map HA serials to hostnames: %v", err)
	} else {
		if vdom == "" {
			vdom = stats.VDOM
		}
		for _, member := range stats.Results {
			hostnames[member.SerialNo] = member.Hostname
		}
	}

	m := []prometheus.Metric{}
	for _, response := range res.Results {
		hostname := hostnames[response.SerialNo]
		if hostname == "" {
			hostname = response.SerialNo
		}
		m = append(m, prometheus.MustNewConstMetric(isMaster, prometheus.GaugeValue, float64(response.IsManageMaster), "manage_master", response.SerialNo, hostname, vdom))
		m = append(m, prometheus.MustNewConstMetric(isMaster, prometheus.GaugeValue, float64(response.IsRootMaster), "root_master", response.SerialNo, hostname, vdom))
	}

	return m, true
}
