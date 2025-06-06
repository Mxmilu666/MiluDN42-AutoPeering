package models

type PeerType string

const (
	PeerTypeIPv4Only                                PeerType = "ipv4_only"
	PeerTypeIPv6Only                                PeerType = "ipv6_only"
	PeerTypeDualStack                               PeerType = "ipv4_ipv6"
	PeerTypeMultiProtocol                           PeerType = "ipv6_multi_protocol"
	PeerTypeIPv6OnlyLocalLinkv6                     PeerType = "ipv6_only_local_linkv6"
	PeerTypeDualStackLocalLinkv6                    PeerType = "ipv4_ipv6_local_linkv6"
	PeerTypeMultiProtocolLocalLinkv6                PeerType = "ipv6_multi_protocol_local_linkv6"
	PeerTypeMultiProtocolExtendedNextHop            PeerType = "ipv6_multi_protocol_extended_next_hop"
	PeerTypeMultiProtocolExtendedNextHopLocalLinkv6 PeerType = "ipv6_multi_protocol_extended_next_hop_local_linkv6"
)
