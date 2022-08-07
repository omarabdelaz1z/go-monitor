package main

import "testing"

func TestIncr(t *testing.T) {
	A := NewNetStat(20, 5)
	B := NewNetStat(10, 3)
	A.Incr(B)

	if A.bytesSent != 30 || A.bytesRecv != 8 || A.bytesTotal != 38 {
		t.Errorf("Incr failed")
	}
}

func TestDelta(t *testing.T) {
	A := NewNetStat(20, 5)
	B := NewNetStat(10, 3)

	C := A.Delta(B)

	if C.bytesSent != 10 || C.bytesRecv != 2 || C.bytesTotal != 12 {
		t.Errorf("got: %d %d %d. %s", C.bytesSent, C.bytesRecv, C.bytesTotal, "expected: 10 3 12")
	}
}
