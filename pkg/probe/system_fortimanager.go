package probe

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type fortimanagerStatusResponse struct {
	Results map[string]interface{} `json:"results"`
	VDOM    string                 `json:"vdom"`
}

type fortimanagerStatus struct {
	VDOM           string
	Mode           string
	StatusID       *int
	RegistrationID *int
}

var (
	fortimanagerStatusText = map[string]int{
		"down":         0,
		"disconnected": 0,
		"disabled":     0,
		"handshake":    1,
		"connecting":   1,
		"inprogress":   1,
		"up":           2,
		"connected":    2,
		"established":  2,
	}
	fortimanagerRegistrationStatusText = map[string]int{
		"unknown":      0,
		"inprogress":   1,
		"in-progress":  1,
		"registering":  1,
		"registered":   2,
		"authorized":   2,
		"unregistered": 3,
		"unauthorized": 3,
	}
)

func probeSystemFortimanagerStatus(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		fortimanagerStatusDesc = prometheus.NewDesc(
			"fortigate_fortimanager_connection_status",
			"Fortimanager status ID",
			[]string{"vdom", "mode", "status"}, nil,
		)
		fortimanagerRegistrationDesc = prometheus.NewDesc(
			"fortigate_fortimanager_registration_status",
			"Fortimanager registration status ID",
			[]string{"vdom", "mode", "status"}, nil,
		)
	)

	statusResponses, skip, err := fetchFortimanagerStatus(c, meta)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	if skip {
		log.Printf("Fortimanager status endpoint not available (HTTP 404), skipping probe")
		return nil, true
	}

	statuses := normalizeFortimanagerStatuses(statusResponses)
	if len(statuses) == 0 {
		log.Printf("Fortimanager status endpoint returned no recognizable fields, skipping probe")
		return nil, true
	}

	m := []prometheus.Metric{}
	for _, r := range statuses {
		mode := r.Mode
		if mode == "" {
			mode = "unknown"
		}

		if r.StatusID != nil {
			statusDown, statusHandshake, statusUp := 0.0, 0.0, 0.0
			switch *r.StatusID {
			case 1:
				statusHandshake = 1.0
			case 2:
				statusUp = 1.0
			default:
				statusDown = 1.0
			}

			m = append(m, prometheus.MustNewConstMetric(fortimanagerStatusDesc, prometheus.GaugeValue, statusDown, r.VDOM, mode, "down"))
			m = append(m, prometheus.MustNewConstMetric(fortimanagerStatusDesc, prometheus.GaugeValue, statusHandshake, r.VDOM, mode, "handshake"))
			m = append(m, prometheus.MustNewConstMetric(fortimanagerStatusDesc, prometheus.GaugeValue, statusUp, r.VDOM, mode, "up"))
		}

		if r.RegistrationID != nil {
			registrationUnknown, registrationInProgress, registrationRegistered, registrationUnregistered := 0.0, 0.0, 0.0, 0.0
			switch *r.RegistrationID {
			case 1:
				registrationInProgress = 1.0
			case 2:
				registrationRegistered = 1.0
			case 3:
				registrationUnregistered = 1.0
			default:
				registrationUnknown = 1.0
			}

			m = append(m, prometheus.MustNewConstMetric(fortimanagerRegistrationDesc, prometheus.GaugeValue, registrationUnknown, r.VDOM, mode, "unknown"))
			m = append(m, prometheus.MustNewConstMetric(fortimanagerRegistrationDesc, prometheus.GaugeValue, registrationInProgress, r.VDOM, mode, "inprogress"))
			m = append(m, prometheus.MustNewConstMetric(fortimanagerRegistrationDesc, prometheus.GaugeValue, registrationRegistered, r.VDOM, mode, "registered"))
			m = append(m, prometheus.MustNewConstMetric(fortimanagerRegistrationDesc, prometheus.GaugeValue, registrationUnregistered, r.VDOM, mode, "unregistered"))
		}
	}

	return m, true
}

func fetchFortimanagerStatus(c http.FortiHTTP, meta *TargetMetadata) ([]fortimanagerStatusResponse, bool, error) {
	preferred := []string{
		"api/v2/monitor/system/fortimanager/status",
		"api/v2/monitor/system/central-management/status",
	}
	if versionAtLeast(meta, 7, 6) {
		preferred[0], preferred[1] = preferred[1], preferred[0]
	}

	had404 := false
	for _, path := range preferred {
		var res []fortimanagerStatusResponse
		if err := c.Get(path, "vdom=*", &res); err != nil {
			var apiErr http.APIError
			if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
				had404 = true
				continue
			}
			return nil, false, err
		}
		return res, false, nil
	}

	if had404 {
		return nil, true, nil
	}
	return nil, false, errors.New("fortimanager status endpoint unavailable")
}

func normalizeFortimanagerStatuses(rs []fortimanagerStatusResponse) []fortimanagerStatus {
	statuses := make([]fortimanagerStatus, 0, len(rs))
	for _, r := range rs {
		statusID, hasStatus := normalizeFortimanagerValue(r.Results, []string{"fortimanager_status_id", "status", "status_id", "conn_status", "conn-status", "connection_status", "connection-status"}, fortimanagerStatusText)
		registrationID, hasRegistration := normalizeFortimanagerValue(r.Results, []string{"fortimanager_registration_status_id", "registration_status", "registration-status", "registration_status_id"}, fortimanagerRegistrationStatusText)

		if !hasStatus && !hasRegistration {
			log.Printf("Fortimanager status response for VDOM %q missing recognizable status fields", r.VDOM)
			continue
		}

		mode := ""
		if v, ok := r.Results["mode"].(string); ok {
			mode = v
		}

		statuses = append(statuses, fortimanagerStatus{
			VDOM:           r.VDOM,
			Mode:           mode,
			StatusID:       statusID,
			RegistrationID: registrationID,
		})
	}
	return statuses
}

func normalizeFortimanagerValue(fields map[string]interface{}, keys []string, textMap map[string]int) (*int, bool) {
	for _, k := range keys {
		if v, ok := fields[k]; ok {
			if normalized, ok := normalizeStatusValue(v, textMap); ok {
				return &normalized, true
			}
		}
	}
	return nil, false
}

func normalizeStatusValue(value interface{}, textMap map[string]int) (int, bool) {
	switch v := value.(type) {
	case float64:
		return int(v), true
	case int:
		return v, true
	case string:
		normalized := strings.ToLower(v)
		normalized = strings.ReplaceAll(normalized, " ", "")
		normalized = strings.ReplaceAll(normalized, "-", "")
		if mapped, ok := textMap[normalized]; ok {
			return mapped, true
		}
		if intValue, err := strconv.Atoi(normalized); err == nil {
			return intValue, true
		}
	}
	return 0, false
}
