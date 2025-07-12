# `elasticp`

`elasticp` is a lightweight **adaptive goroutine pool** written in Go.
It provides a simple mechanism for executing high volumes of small CPU-bound tasks using a worker pool with elastic capacity.

### 💡 What is it?

This is a **proof of concept** for benchmarking goroutine pools vs spawning native goroutines — especially under high concurrency.

It includes:

* A dynamic pool of workers (`Grow`, `Shrink`)
* A `Submit` method that dispatches tasks in a round-robin-like fashion with non-blocking fallback
* A benchmark suite comparing raw goroutines vs `elasticp` under load

---

### 🧪 Benchmarks

Install graphviz
```shell
brew install graphviz
```

Run microbenchmark
```shell
go test -bench=. -benchtime=3s -count=1 -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof
```

```shell
go tool pprof -http=: mem.prof
```

```shell
go tool pprof -http=: cpu.prof
```

Available benchmarks:

| Benchmark                          | Description                                              |
| ---------------------------------- | -------------------------------------------------------- |
| `BenchmarkSequential100kTasks`     | Runs all tasks sequentially, without goroutines          |
| `BenchmarkRawGoroutines_100kTasks` | Spawns 100,000 native goroutines for 100,000 small tasks |
| `BenchmarkElasticpPool100kTasks`   | Uses the pool to dispatch 100,000 parallel tasks         |

### Results
```shell
BenchmarkSequential100kTasks-8              3946            834117 ns/op        16007171 B/op          1 allocs/op
BenchmarkRawGoroutines_100kTasks-8           310          11590516 ns/op        24007883 B/op     100003 allocs/op
BenchmarkElasticpPool100kTasks-8             542           6619665 ns/op        16007207 B/op          2 allocs/op
```

This benchmark setup shows how `elasticp` handles **100,000 concurrent tasks** working over shared memory — and how it compares against both sequential and native goroutine execution models.

---

### 🧬 Task Model

A `Task` is a unit of work with:

* An `Input` slice of float64s
* An `Output` slice to write the result to
* A `WaitGroup` to signal completion

```go
type Task struct {
	Input  []float64
	Output []float64
	Wg     *sync.WaitGroup
}
```

The pool workers perform simple vector addition:

```go
for i, v := range task.Input {
	task.Output[i] += v
}
```

---

### ⚙️ API Overview

#### Create a new pool

```go
pool := elasticp.New()
```

#### Grow the pool

```go
pool.Grow(8) // spawn 8 workers
```

#### Submit work

```go
pool.Submit(elasticp.Task{
	Input:  input,
	Output: output,
	Wg:     &wg,
})
```

#### Shrink the pool

```go
pool.Shrink(4) // remove 4 workers
```

---

### 📈 Why?

This POC helps measure:

* The overhead of goroutine creation vs reuse
* The tradeoffs between pool saturation and raw concurrency
* Fair scheduling under CPU-bound workloads

---

### 🧠 Notes

* This project is designed for benchmarking, not production.
* Tasks must be **independent and stateless**.
* You must call `Wg.Add()` before `Submit()` and `Wg.Wait()` afterward to block until all tasks finish.