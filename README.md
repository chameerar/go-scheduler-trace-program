# Go Scheduler Trace Program

A multi-threaded Go program designed to observe and analyze Go's scheduler behavior using the `GODEBUG=schedtrace` output.

## Overview

This program creates multiple goroutines performing different types of concurrent work:

- **Worker Pool** (6 goroutines): Process CPU-intensive tasks sequentially
- **Producer-Consumer** (1 producer + 4 consumers): Continuous data generation and processing
- **Monitor** (1 goroutine): Periodic status checks
- **Compute Workers** (8 goroutines): Sustained CPU load with heavy computation

Total: ~20 active goroutines producing realistic scheduling patterns.

## Building & Running

### Standard Run
```bash
cd /Users/chameerar/DEV/go/schedtrace
go build -o schedtrace main.go
./schedtrace
```

### With Scheduler Trace
```bash
GODEBUG=schedtrace=10 GOMAXPROCS=2 ./schedtrace
```

The program runs for 60 seconds, creating rich scheduler activity for analysis.

## Understanding Scheduler Trace Output

The `GODEBUG=schedtrace` flag outputs periodic scheduler snapshots like:
```
SCHED 6233ms: gomaxprocs=2 idleprocs=2 threads=4 spinningthreads=0 needspinning=0 idlethreads=2 runqueue=0 [ 0 0 ] schedticks=[ 16 16 ]
```

### Field Reference

| Field | Meaning |
|-------|---------|
| **SCHED 6233ms** | Timestamp in milliseconds since start |
| **gomaxprocs=2** | Max OS threads allowed (set via GOMAXPROCS) |
| **idleprocs=2** | Number of idle processors with no work |
| **threads=4** | Total OS threads created by runtime |
| **spinningthreads=0** | Threads actively waiting for work |
| **needspinning=0** | Whether new spinning thread is needed |
| **idlethreads=2** | Number of idle OS threads |
| **runqueue=0** | Goroutines in global run queue |
| **[ 0 0 ]** | Per-processor run queue lengths |
| **schedticks=[ 16 16 ]** | Scheduling decisions per processor |

### What to Look For

- **High idleprocs**: Program isn't keeping all processors busy
- **Growing runqueue**: More work arriving than processors can handle
- **Balanced schedticks**: Fair scheduling across processors
- **Spinningthreads > 0**: Runtime is spinning looking for work

## Example Analysis

With `GOMAXPROCS=2`, you should see:
- Initially low idleprocs as CPU-intensive work runs
- Eventually higher idleprocs when waiting on channels/timers
- Regular spikes in runqueue as new tasks arrive

## Customization

Edit [main.go](main.go) to adjust:
- Number of workers: Line ~92 (`wg.Add(6)`)
- Task count: Line ~98 (loop to `20`)
- Number of consumers: Line ~108 (`wg.Add(4)`)
- Number of compute workers: Line ~119 (`wg.Add(8)`)
- Program duration: Line ~88 (`60*time.Second`)
- GOMAXPROCS value: Modify the environment variable when running
