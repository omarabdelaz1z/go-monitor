package main

import (
	"testing"

	m "github.com/omarabdelaz1z/go-monitor/pkg/monitoor"
)

func TestIncr(t *testing.T) {
	A := &m.NetStat{
		BytesSent:  20,
		BytesRecv:  5,
		BytesTotal: 25,
	}

	B := &m.NetStat{
		BytesSent:  10,
		BytesRecv:  3,
		BytesTotal: 13,
	}

	C := Incr(A, B)

	if C.BytesSent != 30 || C.BytesRecv != 8 || C.BytesTotal != 38 {
		t.Errorf("Incr failed")
	}
}

func TestDelta(t *testing.T) {
	A := &m.NetStat{
		BytesSent:  20,
		BytesRecv:  5,
		BytesTotal: 25,
	}

	B := &m.NetStat{
		BytesSent:  10,
		BytesRecv:  3,
		BytesTotal: 13,
	}

	C := Delta(A, B)

	if C.BytesSent != 10 || C.BytesRecv != 2 || C.BytesTotal != 12 {
		t.Errorf("got: %d %d %d. %s", C.BytesSent, C.BytesRecv, C.BytesTotal, "expected: 10 3 12")
	}
}
