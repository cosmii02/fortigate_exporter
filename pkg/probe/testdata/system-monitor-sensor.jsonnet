# /api/v2/monitor/system/monitor-sensor?vdom=root (FortiGate 90G, FortiOS 7.6.7)
{
  "http_method": "GET",
  "results": {
    "hw-monitor": [
      {
        "id": "power.psu_[1]",
        "name": "PSU [1]",
        "type": "power",
        "value": 1.0,
        "alarm": false,
        "thresholds": {}
      },
      {
        "id": "power.psu_[2]",
        "name": "PSU [2]",
        "type": "power",
        "value": 1.0,
        "alarm": false,
        "thresholds": {}
      },
      {
        "id": "fan.sys_fan_speed",
        "name": "SYS Fan Speed",
        "type": "fan",
        "value": 818.0,
        "alarm": false,
        "thresholds": {
          "lower_critical": 300.0,
          "upper_critical": 14600.0
        }
      },
      {
        "id": "temperature.tmp1_external_temperature",
        "name": "TMP1 External Temperature",
        "type": "temperature",
        "value": 39.875,
        "alarm": false,
        "thresholds": {}
      },
      {
        "id": "temperature.tmp4_temperature",
        "name": "TMP4 Temperature",
        "type": "temperature",
        "value": 31.0,
        "alarm": false,
        "thresholds": {
          "upper_non_critical": 98.875,
          "upper_critical": 113.0
        }
      },
      {
        "id": "temperature.cpu_on-die_temperature",
        "name": "CPU ON-DIE Temperature",
        "type": "temperature",
        "value": 44.214,
        "alarm": false,
        "thresholds": {
          "upper_non_critical": 103.25,
          "upper_critical": 118.0
        }
      },
      {
        "id": "temperature.b50182_temperature",
        "name": "B50182 Temperature",
        "type": "temperature",
        "value": 49.325,
        "alarm": false,
        "thresholds": {
          "upper_non_critical": 87.5,
          "upper_critical": 100.0
        }
      },
      {
        "id": "temperature.bcm84887_1_temperature",
        "name": "BCM84887_1 Temperature",
        "type": "temperature",
        "value": 40.0,
        "alarm": false,
        "thresholds": {
          "upper_non_critical": 87.5,
          "upper_critical": 100.0
        }
      },
      {
        "id": "temperature.bcm84887_2_temperature",
        "name": "BCM84887_2 Temperature",
        "type": "temperature",
        "value": 38.0,
        "alarm": false,
        "thresholds": {
          "upper_non_critical": 87.5,
          "upper_critical": 100.0
        }
      }
    ]
  },
  "vdom": "root",
  "path": "system",
  "name": "monitor-sensor",
  "status": "success",
  "serial": "FGT90GTK00000000",
  "version": "v7.6.7",
  "build": 3704
}
