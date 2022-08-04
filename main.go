package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/shirou/gopsutil/v3/net"
)

var (
	delay time.Duration = 1 * time.Second
	total *NetStat      = &NetStat{
		bytesSent:  0,
		bytesRecv:  0,
		bytesTotal: 0,
	}
	logger *log.Logger = log.New(os.Stdout, fmt.Sprintf("[ %s ] ", time.Now().Format(time.RFC822)), 0)

	RateTextFunc     = color.New(color.FgBlack).Add(color.BgHiWhite).SprintFunc()
	DownloadTextFunc = color.New(color.FgGreen).SprintfFunc()
	UploadTextFunc   = color.New(color.FgCyan).SprintfFunc()
	TotalTextFunc    = color.New(color.FgYellow).SprintfFunc()
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

func PrettyStat(stat *NetStat) {
	prettySent := ByteRepr(stat.bytesSent)
	prettyRecv := ByteRepr(stat.bytesRecv)
	prettyTotal := ByteRepr(stat.bytesTotal)

	logger.Printf(
		"%s %s %s\n",
		DownloadTextFunc("download: %s", RateTextFunc(prettyRecv)),
		UploadTextFunc("upload: %s", RateTextFunc(prettySent)),
		TotalTextFunc("total: %s", RateTextFunc(prettyTotal)),
	)
}

// This function is written & modified from Go blueprints.
// A website that provide code for com­mon tasks is a collection of handy code examples.
// https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/

func ByteRepr(bytes uint64) string {
	const unit = 1000

	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0

	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "kM"[exp])
}

func Cls() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig

		Cls()

		fmt.Print("Capture:\n")
		PrettyStat(total)
		os.Exit(0)
	}()

	if stats, err := net.IOCounters(false); err != nil {
		panic(err)
	} else {
		lastNetStat := NewNetStat(stats[0].BytesSent, stats[0].BytesRecv)

		for {
			stats, err := net.IOCounters(false)

			if err != nil {
				panic(err)
			}

			currentNetStat := NewNetStat(stats[0].BytesSent, stats[0].BytesRecv)
			delta := currentNetStat.Delta(lastNetStat)

			Cls()
			PrettyStat(delta)

			total.Incr(delta)
			lastNetStat = currentNetStat

			time.Sleep(delay)
		}
	}
}
