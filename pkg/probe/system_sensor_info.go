package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type SystemSensorInfoResults struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
	Alarm bool    `json:"alarm"`
}

type SystemSensorInfo struct {
	Results []SystemSensorInfoResults `json:"results"`
}

type SystemMonitorSensor struct {
	Results struct {
		HWMonitor []SystemSensorInfoResults `json:"hw-monitor"`
	} `json:"results"`
}

func probeSystemSensorInfo(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		sensorTemperature = prometheus.NewDesc(
			"fortigate_sensor_temperature_celsius",
			"Sensor temperature in degree celsius",
			[]string{"name"}, nil,
		)
		sensorFan = prometheus.NewDesc(
			"fortigate_sensor_fan_rpm",
			"Sensor fan rotation speed in RPM",
			[]string{"name"}, nil,
		)
		sensorVoltage = prometheus.NewDesc(
			"fortigate_sensor_voltage_volts",
			"Sensor voltage in volts",
			[]string{"name"}, nil,
		)
		sensorPower = prometheus.NewDesc(
			"fortigate_sensor_power_status",
			"Power supply status as reported by the sensor (1 = present and OK)",
			[]string{"name"}, nil,
		)
		sensorAlarm = prometheus.NewDesc(
			"fortigate_sensor_alarm",
			"Whether the sensor is outside its thresholds according to FortiOS (1 = alarm)",
			[]string{"name", "type"}, nil,
		)
	)

	var sensors []SystemSensorInfoResults
	var res SystemSensorInfo
	if err := c.Get("api/v2/monitor/system/sensor-info", "vdom=root", &res); err == nil {
		sensors = res.Results
	} else {
		// Desktop models (e.g. 40F, 60F, 90G) answer sensor-info with 404 but
		// expose the same IPMC sensors through monitor-sensor.
		var ms SystemMonitorSensor
		if err2 := c.Get("api/v2/monitor/system/monitor-sensor", "vdom=root", &ms); err2 != nil {
			log.Printf("Warning: %v; fallback: %v", err, err2)
			return nil, false
		}
		sensors = ms.Results.HWMonitor
	}

	m := []prometheus.Metric{}
	for _, r := range sensors {
		switch r.Type {
		case "temperature":
			m = append(m, prometheus.MustNewConstMetric(sensorTemperature, prometheus.GaugeValue, r.Value, r.Name))
		case "fan":
			m = append(m, prometheus.MustNewConstMetric(sensorFan, prometheus.GaugeValue, r.Value, r.Name))
		case "voltage":
			m = append(m, prometheus.MustNewConstMetric(sensorVoltage, prometheus.GaugeValue, r.Value, r.Name))
		case "power":
			m = append(m, prometheus.MustNewConstMetric(sensorPower, prometheus.GaugeValue, r.Value, r.Name))
		default:
			continue
		}
		alarm := 0.0
		if r.Alarm {
			alarm = 1.0
		}
		m = append(m, prometheus.MustNewConstMetric(sensorAlarm, prometheus.GaugeValue, alarm, r.Name, r.Type))
	}

	return m, true
}
