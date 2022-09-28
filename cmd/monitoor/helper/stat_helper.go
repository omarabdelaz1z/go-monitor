package helper

import m "github.com/omarabdelaz1z/go-monitor/pkg/monitoor"

// Increment the current netstat by the other netstat.
func Incr(current, new *m.NetStat) m.NetStat {
	return m.NetStat{
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

func UpdateWith(old *m.NetStat, new m.NetStat) {
	old.BytesSent = new.BytesSent
	old.BytesRecv = new.BytesRecv
	old.BytesTotal = new.BytesTotal
}
