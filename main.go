package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/omarabdelaz1z/go-monitor/util"
	"github.com/shirou/gopsutil/v3/net"
)

const (
	DELAY                = 1 * time.Second
	PREVIOUS_LINE string = "\033[1A\033[K" // hacky way to clear the previous line.
)

var (
	logger         *log.Logger = log.New(os.Stdout, fmt.Sprintf("[%s] ", time.Now().Format(time.RFC822)), 0)
	cumulativeStat *NetStat    = NewNetStat(0, 0)
)

func PrintStat(schan <-chan *NetStat, quit <-chan bool) {
	for {
		select {
		case <-quit:
			return
		case stat, ok := <-schan:
			if ok {
				cumulative := util.ByteCountSI(cumulativeStat.bytesTotal)

				logger.Printf(
					"%s %s\n",
					stat,
					util.CumulativeTextFunc("cumulative: %s", util.RateTextFunc(cumulative)),
				)
			}
		}
	}
}

func CaptureStat(buffer chan<- *NetStat, quit <-chan bool) {
	stats, err := net.IOCounters(false)

	if err != nil {
		logger.Panicf("failed to start capturing stats: %s", err)
	}

	lastStat := NewNetStat(stats[0].BytesSent, stats[0].BytesRecv)

	for {
		select {
		case <-quit:
			return
		default:
			stats, err := net.IOCounters(false)

			if err != nil {
				logger.Panicf("failed to get net stats: %s", err)
			}

			netstat := NewNetStat(stats[0].BytesSent, stats[0].BytesRecv)

			delta := netstat.Delta(lastStat)
			cumulativeStat.Incr(delta)

			buffer <- delta

			// replace previous stat with current stat.
			lastStat = netstat

			time.Sleep(DELAY)
			fmt.Print(PREVIOUS_LINE)
		}
	}
}

func CaptureInterrupt(sig <-chan os.Signal, quit chan<- bool) {
	<-sig
	quit <- true
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	quit := make(chan bool, 1)
	buffer := make(chan *NetStat)

	go CaptureInterrupt(sig, quit)
	go CaptureStat(buffer, quit)
	go PrintStat(buffer, quit)

	<-quit
	close(quit)
	close(buffer)
	close(sig)

	fmt.Print("Captured: \n")
	logger.Print(cumulativeStat)
}
