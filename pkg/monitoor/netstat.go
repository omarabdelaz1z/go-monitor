package monitoor

// A Netstat represents a network statistics snapshot.
type NetStat struct {
	BytesSent  uint64
	BytesRecv  uint64
	BytesTotal uint64
}
