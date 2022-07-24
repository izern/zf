package test

func init() {

}

var TestYaml = `
allow-lan: true
external-controller: 0.0.0.0:9090
log-level: info
mode: Rule
port: 7890
proxies:
  - cipher: aes-256-gcm
    name: '[07-23]|oslook|斯洛伐克(SK)Slovakia/Bratislava_1'
    password: TEzjfAYq2IjtuoS
    port: 6697
    server: 141.164.39.146
    type: ss
  - name: '[07-23]|oslook|中国香港特别行政区(HK)Hongkong+SAR+China/Hong+Kong_42'
    password: 7b4066ae-accc-11eb-a8bf-f23c91cfbbc9
    port: 443
    server: ssl.tcpbbr.net
    skip-cert-verify: true
    type: trojan
    udp: true
proxy-groups:
  - name: 🚀 节点选择
    proxies:
      - ♻️ 自动选择
      - DIRECT
      - '[07-23]|oslook|美国(US)USA/LosAngeles_15'
      - '[07-23]|oslook|日本(JP)Japan/Tokyo_16'

    type: select
  - interval: 300
    name: ♻️ 自动选择
    proxies:
      - '[07-23]|oslook|日本(JP)Japan/Osaka_40'
      - '[07-23]|oslook|日本(JP)Japan/Osaka_41'
      - '[07-23]|oslook|中国香港特别行政区(HK)Hongkong+SAR+China/Hong+Kong_42'
    tolerance: 50
    type: url-test
    url: http://www.gstatic.com/generate_204
rules:
  - DOMAIN-SUFFIX,acl4.ssr,🎯 全球直连
  - DOMAIN-SUFFIX,ip6-localhost,🎯 全球直连
  - DOMAIN-SUFFIX,ip6-loopback,🎯 全球直连
  - DOMAIN-SUFFIX,local,🎯 全球直连
  - DOMAIN-SUFFIX,localhost,🎯 全球直连
  - IP-CIDR,10.0.0.0/8,🎯 全球直连,no-resolve
  - IP-CIDR,100.64.0.0/10,🎯 全球直连,no-resolve
  - IP-CIDR,127.0.0.0/8,🎯 全球直连,no-resolve
  - IP-CIDR,172.16.0.0/12,🎯 全球直连,no-resolve
  - IP-CIDR,192.168.0.0/16,🎯 全球直连,no-resolve
  - IP-CIDR,198.18.0.0/16,🎯 全球直连,no-resolve
  - IP-CIDR6,::1/128,🎯 全球直连,no-resolve
  - IP-CIDR6,fc00::/7,🎯 全球直连,no-resolve
  - IP-CIDR6,fe80::/10,🎯 全球直连,no-resolve
  - IP-CIDR6,fd00::/8,🎯 全球直连,no-resolve
  - GEOIP,CN,🎯 全球直连
  - MATCH,🐟 漏网之鱼
  - DOMAIN-SUFFIX,dyndns.com,🎯 全球直连
  - DOMAIN-SUFFIX,dyndns.org,🎯 全球直连
socks-port: 7891
`
