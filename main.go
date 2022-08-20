package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/omarabdelaz1z/go-monitor/model"
	m "github.com/omarabdelaz1z/go-monitor/monitoor"
	"github.com/omarabdelaz1z/go-monitor/util"
)

var (
	cumulativeStat *m.NetStat = &m.NetStat{
		BytesSent:  0,
		BytesRecv:  0,
		BytesTotal: 0,
	}
)

const (
	DELAY                = 1 * time.Second
	PREVIOUS_LINE string = "\033[1A\033[K" // hacky way to clear the previous line.

	PERIOD int64 = 1 * 60 * 60 // 1 hour

	DB_DATASOURCE string = "monitor.db"
	DB_DRIVER     string = "sqlite3"
)

func DisplayStat(schan <-chan *m.NetStat, quit <-chan bool) {
	for {
		select {
		case <-quit:
			return
		case stat, ok := <-schan:
			if ok {
				cumulative := util.ByteCountSI(cumulativeStat.BytesTotal)

				log.Printf(
					"%s %s\n",
					stat,
					util.CumulativeTextFunc("cumulative: %s", util.RateTextFunc(cumulative)),
				)
			}
		}
	}
}

func CaptureStat(buffer chan<- *m.NetStat, quit chan bool) {
	lastStat, err := m.Brief()

	var (
		t0       = time.Now().Unix()
		periodic = &m.NetStat{
			BytesSent:  0,
			BytesRecv:  0,
			BytesTotal: 0,
		}
	)

	if err != nil {
		quit <- true
	}

	for {
		select {
		case <-quit:
			return
		default:
			netstat, err := m.Brief()

			if err != nil {
				quit <- true
			}

			delta := netstat.Delta(lastStat)
			buffer <- delta // send the delta to the display goroutine.

			t1 := time.Now().Unix()

			if t1-t0 >= PERIOD {
				t0 = t1

				model.Insert(&model.Snapshot{
					Timestamp: t1,
					Sent:      periodic.BytesSent,
					Received:  periodic.BytesRecv,
					Total:     periodic.BytesTotal,
				})

				periodic.BytesSent = 0
				periodic.BytesRecv = 0
				periodic.BytesTotal = 0
			}

			periodic.Incr(delta)
			cumulativeStat.Incr(delta)

			// replace previous stat with current stat
			// to reflect the next measurement.
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
	buffer := make(chan *m.NetStat)

	err := model.InitDb(DB_DRIVER, DB_DATASOURCE)

	if err != nil {
		quit <- true
	}

	go CaptureInterrupt(sig, quit)
	go CaptureStat(buffer, quit)
	go DisplayStat(buffer, quit)

	<-quit
	close(quit)
	close(buffer)
	close(sig)

	fmt.Println("Captured: ")
	log.Print(cumulativeStat)
}
