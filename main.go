package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/omarabdelaz1z/go-monitor/model"
	m "github.com/omarabdelaz1z/go-monitor/monitoor"
	"github.com/omarabdelaz1z/go-monitor/util"
)

var (
	mu sync.RWMutex

	cumulativeStat *m.NetStat = &m.NetStat{
		BytesSent:  0,
		BytesRecv:  0,
		BytesTotal: 0,
	}
	periodicStat *m.NetStat = &m.NetStat{
		BytesSent:  0,
		BytesRecv:  0,
		BytesTotal: 0,
	}
)

const (
	DELAY  time.Duration = 1 * time.Second
	PERIOD time.Duration = 1 * time.Hour

	DSN    string = "monitor.db"
	DRIVER string = "sqlite3"

	PREVIOUS_LINE string = "\033[1A\033[K" // hacky way to clear the previous line.
)

func Display(schan <-chan *m.NetStat, quit <-chan bool) {
	for {
		select {
		case <-quit:
			return
		case stat, ok := <-schan:
			if ok {
				mu.RLock()
				currentTotal := cumulativeStat.BytesTotal
				mu.RUnlock()

				cumulative := util.ByteCountSI(currentTotal)

				log.Printf(
					"%s %s\n",
					Pretty(stat),
					util.CumulativeTextFunc("cumulative: %s", util.RateTextFunc(cumulative)),
				)
			}
		}
	}
}

// TODO: a side effect resulted from fmt.Print.
func Monitor(buffer chan<- *m.NetStat, quit chan bool) {
	lastStat, err := m.Brief()

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

			delta := Delta(netstat, lastStat)
			buffer <- delta // send the delta to the display goroutine.

			mu.Lock()

			periodicStat = Incr(periodicStat, delta)
			cumulativeStat = Incr(cumulativeStat, delta)

			mu.Unlock()

			lastStat = netstat // record the next stat.

			time.Sleep(DELAY)
			fmt.Print(PREVIOUS_LINE)
		}
	}
}

// TODO: the function does two things.
func Persist(ticker *time.Ticker, quit <-chan bool) {
	for {
		select {
		case <-ticker.C:
			mu.RLock()

			model.Insert(&model.Snapshot{
				Timestamp: time.Now().Unix(),
				Sent:      periodicStat.BytesSent,
				Received:  periodicStat.BytesRecv,
				Total:     periodicStat.BytesTotal,
			})

			mu.RUnlock()

			mu.Lock()

			periodicStat.BytesSent = 0
			periodicStat.BytesRecv = 0
			periodicStat.BytesTotal = 0

			mu.Unlock()

		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func Shutdown(signals <-chan os.Signal, quit chan<- bool) {
	<-signals
	quit <- true
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	buffer := make(chan *m.NetStat)

	quit := make(chan bool)

	err := model.InitDb(DRIVER, DSN)

	if err != nil {
		quit <- true
	}

	ticker := time.NewTicker(PERIOD)

	go Monitor(buffer, quit)
	go Display(buffer, quit)
	go Persist(ticker, quit)
	go Shutdown(signals, quit)

	<-quit

	defer func() {
		close(quit)
		close(buffer)
		close(signals)
	}()

	fmt.Println("\ncaptured: ")
	log.Print(Pretty(cumulativeStat))
}
