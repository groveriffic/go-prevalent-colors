# Go Prevalent Colors

This is an exercise program to download and process image urls listed in a file.
The output is a CSV containing original URL and top three most prevalent colors.

Docker is used to simulate constrained resources (cpu and memory).  Scripts are
included to make running everything easy and repeatable.

## Tests

Run tests with

    scripts/test

## Build

    scripts/build

## Run

    scripts/run

The run script uses a docker container with limited cpu and memory

## Performance Tuning

Monitoring during run with

    docker stats

Using 1 cpu, 512mb ram and 5 workers.  I saw at or near 100% cpu and memory
utilization for the full run.  The number of workers should be tuned for your
system specs.  Too many workers and you may see swapping or the process OOM
killed.  Too few workers and it may become IO bound and not utilize the
available CPU resources.

## Profiling

Analyzing the cpu profiling output

    go tool pprof -http localhost:8080 cpu.prof

Most of the CPU time is spent in `color Counter Inc` method.  This would be the
first target for optimization.
