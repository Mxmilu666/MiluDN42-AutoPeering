package database

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Peer 代表一个 peering 连接
type Peer struct {
	IPv4      string    `json:"ipv4"`
	IPv6      string    `json:"ipv6"`
	Endpoint  string    `json:"endpoint"`
	PublicKey string    `json:"public_key"`
	PeerTime  time.Time `json:"peer_time"`
	NodeName  string    `json:"node_name"`
}

// PeerNode 代表一个 ASN 下的所有 peering 节点
type PeerNode struct {
	ASN   string `json:"asn"`
	Peers []Peer `json:"peers"`
}

var (
	peersFile = "data/peers.json"
	peersLock sync.Mutex
)

// SavePeerNode 写入或追加一个 ASN 的 peer 信息
func SavePeerNode(asn string, peer Peer) error {
	peersLock.Lock()
	defer peersLock.Unlock()

	nodes, _ := LoadPeerNodes() // 忽略读取错误，文件不存在时视为空
	found := false
	for i, n := range nodes {
		if n.ASN == asn {
			nodes[i].Peers = append(nodes[i].Peers, peer)
			found = true
			break
		}
	}
	if !found {
		nodes = append(nodes, PeerNode{
			ASN:   asn,
			Peers: []Peer{peer},
		})
	}

	f, err := os.OpenFile(peersFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(nodes)
}

// LoadPeerNodes 读取所有 ASN 的 peering 节点
func LoadPeerNodes() ([]PeerNode, error) {
	peersLock.Lock()
	defer peersLock.Unlock()

	f, err := os.Open(peersFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []PeerNode{}, nil
		}
		return nil, err
	}
	defer f.Close()
	var nodes []PeerNode
	err = json.NewDecoder(f).Decode(&nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// DeletePeerNode 删除指定 ASN 下的 peer（按 NodeName 匹配）
func DeletePeerNode(asn, nodeName string) error {
	peersLock.Lock()
	defer peersLock.Unlock()

	nodes, err := LoadPeerNodes()
	if err != nil {
		return err
	}
	changed := false
	for i, n := range nodes {
		if n.ASN == asn {
			newPeers := make([]Peer, 0, len(n.Peers))
			for _, p := range n.Peers {
				if p.NodeName != nodeName {
					newPeers = append(newPeers, p)
				} else {
					changed = true
				}
			}
			nodes[i].Peers = newPeers
		}
	}
	if changed {
		f, err := os.OpenFile(peersFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		return enc.Encode(nodes)
	}
	return nil // 没有变化也视为成功
}
