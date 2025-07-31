# api/v2/monitor/vpn/ipsec?vdom=*
[
  {
    "http_method":"GET",
    "results":[
      {
        "proxyid":[
          {
            "proxy_src":[
              {
                "subnet":"0.0.0.0\/0.0.0.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "proxy_dst":[
              {
                "subnet":"0.0.0.0\/0.0.0.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "status":"up",
            "p2name":"tunnel_1-sub",
            "p2serial":1,
            "expire":11279,
            "incoming_bytes": 14298240,
            "outgoing_bytes": 14248560
          }
        ],
        "name":"tunnel_1",
        "comments":"",
        "wizard-type":"custom",
        "creation_time":270801,
        "type":"automatic",
        "incoming_bytes": 14298240,
        "outgoing_bytes": 14248560,
        "rgwy":"1.2.3.4"
      },
      {
        "proxyid":[
          {
            "proxy_src":[
              {
                "subnet":"192.168.1.0\/255.255.255.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "proxy_dst":[
              {
                "subnet":"10.0.0.0\/255.0.0.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "status":"up",
            "p2name":"nixvpn-split-client1",
            "p2serial":101,
            "expire":3600,
            "incoming_bytes": 1024000,
            "outgoing_bytes": 2048000
          },
          {
            "proxy_src":[
              {
                "subnet":"192.168.1.0\/255.255.255.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "proxy_dst":[
              {
                "subnet":"10.0.0.0\/255.0.0.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "status":"up",
            "p2name":"nixvpn-split-client2",
            "p2serial":102,
            "expire":3600,
            "incoming_bytes": 512000,
            "outgoing_bytes": 1024000
          }
        ],
        "name":"nixvpn-split",
        "comments":"Client VPN tunnel",
        "wizard-type":"dialup",
        "creation_time":270801,
        "type":"dialup",
        "incoming_bytes": 1536000,
        "outgoing_bytes": 3072000,
        "rgwy":"0.0.0.0"
      }
    ],
    "vdom":"root",
    "path":"vpn",
    "name":"ipsec",
    "status":"success",
    "serial":"FGT61FT000000000",
    "version":"v6.0.10",
    "build":365
  }
]
