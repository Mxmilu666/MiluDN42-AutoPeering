# Peer API 所有类型 JSON 示例

## 1. 仅 IPv4（ipv4_only）
```json
{
  "asn": "4242420000",
  "ipv4": "10.0.0.1",
  "ipv6": "",
  "public_ip": "1.2.3.4",
  "extended_next_hop": false,
  "routes": "ipv4"
}
```

## 2. 仅 IPv6（ipv6_only）
```json
{
  "asn": "4242420000",
  "ipv4": "",
  "ipv6": "fd00:1::1",
  "public_ip": "1.2.3.4",
  "extended_next_hop": false,
  "routes": "ipv6"
}
```

## 3. 仅 local-link（ipv6_only_local_linkv6）
```json
{
  "asn": "4242420000",
  "ipv4": "",
  "ipv6": "fe80::1",
  "public_ip": "1.2.3.4",
  "extended_next_hop": false,
  "routes": "ipv6"
}
```

## 4. IPv4 + IPv6（ipv4_ipv6）
```json
{
  "asn": "4242420000",
  "ipv4": "10.0.0.1",
  "ipv6": "fd00:1::1",
  "public_ip": "1.2.3.4",
  "extended_next_hop": false,
  "routes": "both"
}
```

## 5. IPv4 + local-link（ipv4_ipv6_local_linkv6）
```json
{
  "asn": "4242420000",
  "ipv4": "10.0.0.1",
  "ipv6": "fe80::1",
  "public_ip": "1.2.3.4",
  "extended_next_hop": false,
  "routes": "ipv4"
}
```

## 6. IPv6 Multi-Protocol（ipv6_multi_protocol）
```json
{
  "asn": "4242420000",
  "ipv4": "",
  "ipv6": "fd00:1::1",
  "public_ip": "1.2.3.4",
  "extended_next_hop": false,
  "routes": "both"
}
```

## 7. IPv6 Multi-Protocol + local-link（ipv6_multi_protocol_local_linkv6）
```json
{
  "asn": "4242420000",
  "ipv4": "",
  "ipv6": "fe80::1",
  "public_ip": "1.2.3.4",
  "extended_next_hop": false,
  "routes": "both"
}
```

## 8. IPv6 Multi-Protocol Extended Next Hop（ipv6_multi_protocol_extended_next_hop）
```json
{
  "asn": "4242420000",
  "ipv4": "10.0.0.1",
  "ipv6": "fd00:1::1",
  "public_ip": "1.2.3.4",
  "extended_next_hop": true,
  "routes": "both"
}
```

## 9. IPv6 Multi-Protocol Extended Next Hop + local-link（ipv6_multi_protocol_extended_next_hop_local_linkv6）
```json
{
  "asn": "4242420000",
  "ipv4": "10.0.0.1",
  "ipv6": "fe80::1",
  "public_ip": "1.2.3.4",
  "extended_next_hop": true,
  "routes": "both"
}
```
