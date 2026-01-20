# api/v2/monitor/user/firewall?vdom=*
[
  {
    "http_method": "GET",
    "results": [
      {
        "user_name": "alice",
        "src_ip": "10.0.0.5",
        "user_group": "Employees",
        "auth_type": "ldap",
        "duration_secs": 120,
        "auth_server": "LDAP1"
      },
      {
        "user": "bob",
        "ip": "10.0.0.6",
        "group": "Contractors",
        "type": "radius",
        "duration_secs": 30
      },
      {
        "user": "dave",
        "src_ip": "10.0.0.7",
        "group": "Employees",
        "auth_type": 3,
        "duration_secs": 60
      }
    ],
    "vdom": "root",
    "path": "user",
    "name": "firewall",
    "status": "success",
    "serial": "FG10000000000000",
    "version": "v7.0.0",
    "build": 1
  },
  {
    "http_method": "GET",
    "results": [
      {
        "username": "carol",
        "srcip": "10.0.1.7",
        "group": "Employees",
        "auth_type": "local",
        "duration_secs": 90
      }
    ],
    "vdom": "vdom2",
    "path": "user",
    "name": "firewall",
    "status": "success",
    "serial": "FG10000000000000",
    "version": "v7.0.0",
    "build": 1
  }
]
