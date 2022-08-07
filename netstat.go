package main

import (
	"fmt"

	"github.com/omarabdelaz1z/go-monitor/util"
)

type NetStat struct {
	bytesSent  uint64
	bytesRecv  uint64
	bytesTotal uint64
}

func (netStat *NetStat) Incr(new *NetStat) {
	netStat.bytesRecv += new.bytesRecv
	netStat.bytesSent += new.bytesSent
	netStat.bytesTotal += new.bytesTotal
}

func (current *NetStat) Delta(previous *NetStat) *NetStat {
	return &NetStat{
		bytesSent:  current.bytesSent - previous.bytesSent,
		bytesRecv:  current.bytesRecv - previous.bytesRecv,
		bytesTotal: current.bytesTotal - previous.bytesTotal,
	}
}

func NewNetStat(sent uint64, recv uint64) *NetStat {
	return &NetStat{
		bytesSent:  sent,
		bytesRecv:  recv,
		bytesTotal: sent + recv,
	}
}

func (netStat *NetStat) String() string {
	sent := util.ByteCountSI(netStat.bytesSent)
	recv := util.ByteCountSI(netStat.bytesRecv)
	total := util.ByteCountSI(netStat.bytesTotal)

	return fmt.Sprintf("%s %s %s",
		util.UploadTextFunc("upload: %s", util.RateTextFunc(sent)),
		util.DownloadTextFunc("download: %s", util.RateTextFunc(recv)),
		util.TotalTextFunc("total: %s", util.RateTextFunc(total)),
	)
}
