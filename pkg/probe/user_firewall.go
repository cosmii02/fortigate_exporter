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
		firewallUserDuration = prometheus.NewDesc(
			"fortigate_user_firewall_session_duration_seconds",
			"Duration of authenticated firewall sessions in seconds",
			[]string{"vdom", "user", "src_ip", "group", "auth_server"}, nil,
		)
	)

	var res []UserFirewall
	if err := c.Get("api/v2/monitor/user/firewall", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	type firewallUserLabels struct {
		VDOM       string
		User       string
		SrcIP      string
		Group      string
		AuthServer string
	}
	type firewallUserMetric struct {
		Labels      firewallUserLabels
		Duration    float64
		HasDuration bool
	}

	userMetrics := map[string]*firewallUserMetric{}
	for _, r := range res {
		for _, user := range r.Results {
			userName := pickFieldString(user, "user", "user_name", "username", "name")
			srcIP := pickFieldString(user, "src_ip", "srcip", "src-ip", "ip", "ipaddr")
			group := pickFieldString(user, "user_group", "group", "usergroup", "grp")
			authServer := pickFieldString(user, "auth_server", "auth-server", "authserver", "server")

			if userName == "" && srcIP == "" && group == "" && authServer == "" {
				continue
			}

			labels := firewallUserLabels{
				VDOM:       r.VDOM,
				User:       userName,
				SrcIP:      srcIP,
				Group:      group,
				AuthServer: authServer,
			}
			key := strings.Join([]string{labels.VDOM, labels.User, labels.SrcIP, labels.Group, labels.AuthServer}, "\x1f")
			entry, ok := userMetrics[key]
			if !ok {
				entry = &firewallUserMetric{Labels: labels}
				userMetrics[key] = entry
			}

			if duration, ok := pickFieldFloat(user, "duration_secs", "duration_sec", "duration_seconds", "duration"); ok {
				if !entry.HasDuration || duration > entry.Duration {
					entry.Duration = duration
					entry.HasDuration = true
				}
			}
		}
	}

	m := []prometheus.Metric{}
	for _, entry := range userMetrics {
		m = append(m, prometheus.MustNewConstMetric(
			firewallUsers,
			prometheus.GaugeValue,
			1,
			entry.Labels.VDOM,
			entry.Labels.User,
			entry.Labels.SrcIP,
			entry.Labels.Group,
			entry.Labels.AuthServer,
		))
		if entry.HasDuration {
			m = append(m, prometheus.MustNewConstMetric(
				firewallUserDuration,
				prometheus.GaugeValue,
				entry.Duration,
				entry.Labels.VDOM,
				entry.Labels.User,
				entry.Labels.SrcIP,
				entry.Labels.Group,
				entry.Labels.AuthServer,
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

func pickFieldFloat(fields map[string]interface{}, keys ...string) (float64, bool) {
	for _, key := range keys {
		if value, ok := fields[key]; ok {
			if normalized, ok := normalizeFieldFloat(value); ok {
				return normalized, true
			}
		}
	}
	return 0, false
}

func normalizeFieldFloat(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return 0, false
		}
		if parsed, err := strconv.ParseFloat(trimmed, 64); err == nil {
			return parsed, true
		}
	}
	return 0, false
}
