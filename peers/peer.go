package peers

import (
	"encoding/binary"
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func Unmarshal(peersBin []byte) ([]Peer, error) {
	const peerSize = 6 // 4 bytes for IP, 2 for port

	peersCount := len(peersBin) / 6

	peers := make([]Peer, peersCount)
	for i := 0; i < peersCount; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peersBin[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset+4 : offset+6])
	}
	return peers, nil
}

func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}
