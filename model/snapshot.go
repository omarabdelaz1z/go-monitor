package model

type Snapshot struct {
	Timestamp int64  `json:"timestamp"`
	Sent      uint64 `json:"sent"`
	Received  uint64 `json:"received"`
	Total     uint64 `json:"total"`
}
