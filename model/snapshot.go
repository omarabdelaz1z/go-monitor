package model

type Stat struct {
	Sent     uint64
	Received uint64
	Total    uint64
}

type Snapshot struct {
	Timestamp int64
	Stat      Stat
}

type MonthStat struct {
	Month string
	Stat  Stat
}
