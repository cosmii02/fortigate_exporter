package probe

import (
	"log"
	"strconv"
	"strings"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type UserFirewall struct {
	Results []map[string]interface{} `json:"results"`
	VDOM    string                   `json:"vdom"`
}

func probeUserFirewall(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		firewallUsers = prometheus.NewDesc(
			"fortigate_user_firewall_info",
			"Info on authenticated firewall users",
			[]string{"vdom", "user", "src_ip", "group", "auth_server"}, nil,
		)
	)

	var res []UserFirewall
	if err := c.Get("api/v2/monitor/user/firewall", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	seen := map[string]struct{}{}
	for _, r := range res {
		for _, user := range r.Results {
			userName := pickFieldString(user, "user", "user_name", "username", "name")
			srcIP := pickFieldString(user, "src_ip", "srcip", "src-ip", "ip", "ipaddr")
			group := pickFieldString(user, "user_group", "group", "usergroup", "grp")
			authServer := pickFieldString(user, "auth_server", "auth-server", "authserver", "server")

			if userName == "" && srcIP == "" && group == "" && authServer == "" {
				continue
			}

			key := strings.Join([]string{r.VDOM, userName, srcIP, group, authServer}, "\x1f")
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}

			m = append(m, prometheus.MustNewConstMetric(
				firewallUsers,
				prometheus.GaugeValue,
				1,
				r.VDOM,
				userName,
				srcIP,
				group,
				authServer,
			))
		}
	}

	return m, true
}

func pickFieldString(fields map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if value, ok := fields[key]; ok {
			if normalized := normalizeFieldString(value); normalized != "" {
				return normalized
			}
		}
	}
	return ""
}

func normalizeFieldString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case []interface{}:
		parts := make([]string, 0, len(v))
		for _, item := range v {
			if normalized := normalizeFieldString(item); normalized != "" {
				parts = append(parts, normalized)
			}
		}
		if len(parts) > 0 {
			return strings.Join(parts, ",")
		}
	}
	return ""
}
