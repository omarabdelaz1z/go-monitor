package monitoor

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
)

const (
	ALL_INTERFACE bool = false // "all"
)

// An overall network statistics at the current time.
func Brief() (*NetStat, error) {
	stats, err := net.IOCounters(ALL_INTERFACE)

	if err != nil {
		return nil, fmt.Errorf("failed to capture network stat: %v", err)
	}

	return &NetStat{
		BytesSent:  stats[0].BytesSent,
		BytesRecv:  stats[0].BytesRecv,
		BytesTotal: stats[0].BytesSent + stats[0].BytesRecv,
	}, nil
}
