package util

import "github.com/fatih/color"

var (
	// temporary, it will be changed in the future.
	RateTextFunc       = color.New(color.FgBlack).Add(color.BgHiWhite).SprintFunc()
	DownloadTextFunc   = color.New(color.FgGreen).SprintfFunc()
	UploadTextFunc     = color.New(color.FgCyan).SprintfFunc()
	TotalTextFunc      = color.New(color.FgYellow).SprintfFunc()
	CumulativeTextFunc = color.New(color.FgRed).SprintfFunc()
)
