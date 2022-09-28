package util

import "fmt"

// A website that provide code for com­mon tasks is a collection of handy code examples.
// https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func ByteCountSI(bytes uint64) string {
	const unit = 1000

	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0

	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "kMGTPE"[exp])
}
