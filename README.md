# elasticp
elasticp is a lightweight, elastic goroutine pool for Go. It dynamically adjusts the number of worker goroutines based on workload, enabling efficient task processing under varying loads. Supports single-task and batch processing.


## Benchmark

Install graphviz
```shell
brew install graphviz
```

Run microbenchmark
```shell
go test -bench=. -benchtime=3s -count=1 -cpuprofile=cpu.prof -memprofile=mem.prof
```

```shell
go tool pprof -http=: mem.prof
```

```shell
go tool pprof -http=: cpu.prof
```
