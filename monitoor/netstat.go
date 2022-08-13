package monitoor

import (
	"fmt"

	"github.com/omarabdelaz1z/go-monitor/util"
)

// A Netstat represents a network statistics snapshot.
type NetStat struct {
	BytesSent  uint64
	BytesRecv  uint64
	BytesTotal uint64
}

// Increment the current netstat by the other netstat.
func (netStat *NetStat) Incr(new *NetStat) {
	netStat.BytesRecv += new.BytesRecv
	netStat.BytesSent += new.BytesSent
	netStat.BytesTotal += new.BytesTotal
}

// A netstat of the delta between the current netstat and the other netstat.
func (current *NetStat) Delta(previous *NetStat) *NetStat {
	return &NetStat{
		BytesSent:  current.BytesSent - previous.BytesSent,
		BytesRecv:  current.BytesRecv - previous.BytesRecv,
		BytesTotal: current.BytesTotal - previous.BytesTotal,
	}
}

// Formatted string representation of the netstat.
func (netStat *NetStat) String() string {
	sent := util.ByteCountSI(netStat.BytesSent)
	recv := util.ByteCountSI(netStat.BytesRecv)
	total := util.ByteCountSI(netStat.BytesTotal)

	return fmt.Sprintf("%s %s %s",
		util.UploadTextFunc("upload: %s", util.RateTextFunc(sent)),
		util.DownloadTextFunc("download: %s", util.RateTextFunc(recv)),
		util.TotalTextFunc("total: %s", util.RateTextFunc(total)),
	)
}
