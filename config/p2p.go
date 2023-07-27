package config

import "time"

// P2PConfig stores configuration related to peer-to-peer networking.
type P2PConfig struct {
	ListenAddress     string // Address to listen for incoming connections
	Seeds             string // Comma separated list of seed nodes to connect to
	BlockedPeers      string // Comma separated list of nodes to ignore
	AllowedPeers      string // Comma separated list of nodes to whitelist
	reAdvertisePeriod time.Duration
	peerLimit         int32
	txTopicSuffix     string
}
