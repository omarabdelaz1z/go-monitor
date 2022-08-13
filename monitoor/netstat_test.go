package monitoor

import "testing"

func TestIncr(t *testing.T) {
	A := &NetStat{
		BytesSent:  20,
		BytesRecv:  5,
		BytesTotal: 25,
	}

	B := &NetStat{
		BytesSent:  10,
		BytesRecv:  3,
		BytesTotal: 13,
	}

	A.Incr(B)

	if A.BytesSent != 30 || A.BytesRecv != 8 || A.BytesTotal != 38 {
		t.Errorf("Incr failed")
	}
}

func TestDelta(t *testing.T) {
	A := &NetStat{
		BytesSent:  20,
		BytesRecv:  5,
		BytesTotal: 25,
	}

	B := &NetStat{
		BytesSent:  10,
		BytesRecv:  3,
		BytesTotal: 13,
	}

	C := A.Delta(B)

	if C.BytesSent != 10 || C.BytesRecv != 2 || C.BytesTotal != 12 {
		t.Errorf("got: %d %d %d. %s", C.BytesSent, C.BytesRecv, C.BytesTotal, "expected: 10 3 12")
	}
}
