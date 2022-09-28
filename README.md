## Bandwidth Monitoring

A monitoring bandwidth tool implemented in Go

- Currently monitors 'all' the network interfaces on an operating system.
- Soon to support per-process monitoring.

Monitoor Core

- Monitoring with [gopsutil](https://github.com/shirou/gopsutil).
- Data Persistance with `sqlite3`
- Goroutines with: channels, errgroup (a better waitgroup)
- Graceful Shutdown with `os/signal`
- Native logging with `log`

<p align="center">
  <img src="./doc/demo-run.png" style="zoom:50%; alt='Logging Preview'" />
</p>
