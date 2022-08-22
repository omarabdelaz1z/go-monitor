package main

import (
	"fmt"

	m "github.com/omarabdelaz1z/go-monitor/monitoor"
	"github.com/omarabdelaz1z/go-monitor/util"
)

// Increment the current netstat by the other netstat.
func Incr(current, new *m.NetStat) *m.NetStat {
	return &m.NetStat{
		BytesSent:  current.BytesSent + new.BytesSent,
		BytesRecv:  current.BytesRecv + new.BytesRecv,
		BytesTotal: current.BytesTotal + new.BytesTotal,
	}
}

// A netstat of the delta between the current netstat and the other netstat.
func Delta(current, previous *m.NetStat) *m.NetStat {
	return &m.NetStat{
		BytesSent:  current.BytesSent - previous.BytesSent,
		BytesRecv:  current.BytesRecv - previous.BytesRecv,
		BytesTotal: current.BytesTotal - previous.BytesTotal,
	}
}

// Formatted string representation of the netstat.
func Pretty(stat *m.NetStat) string {
	sent := util.ByteCountSI(stat.BytesSent)
	recv := util.ByteCountSI(stat.BytesRecv)
	total := util.ByteCountSI(stat.BytesTotal)

	return fmt.Sprintf("%s %s %s",
		util.UploadTextFunc("upload: %s", util.RateTextFunc(sent)),
		util.DownloadTextFunc("download: %s", util.RateTextFunc(recv)),
		util.TotalTextFunc("total: %s", util.RateTextFunc(total)),
	)
}
