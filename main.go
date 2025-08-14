package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/logananthony/go-baseball/pkg/api"
	"github.com/logananthony/go-baseball/pkg/config"
)

func init() {
	// Loads .env when present (dev). In App Platform, set env vars in the UI/spec.
	_ = godotenv.Load()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db := config.ConnectDB()
	// Reasonable pool defaults; tweak to taste.
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := api.NewAPIServer(":"+port, db)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// package main

// import (
// 	"context"
// 	"crypto/sha256"
// 	"database/sql"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"math"
// 	"math/rand"
// 	"os"
// 	"runtime"
// 	"sync"
// 	"sync/atomic"
// 	"time"

// 	"github.com/joho/godotenv"
// 	_ "github.com/lib/pq"
// 	"github.com/logananthony/go-baseball/pkg/api"
// 	"github.com/logananthony/go-baseball/pkg/config"
// )

// func init() {
// 	_ = godotenv.Load()
// 	rand.Seed(time.Now().UnixNano())
// }

// // ------------------------ Flags ------------------------
// var (
// 	benchMode  = flag.Bool("bench", false, "Run benchmark instead of starting the API server")
// 	kernel     = flag.String("bench.kernel", "fma", "Benchmark kernel: fma | sha256 | matmul | simlike")
// 	duration   = flag.Duration("bench.duration", 30*time.Second, "How long to run the benchmark")
// 	workers    = flag.Int("bench.workers", 0, "Number of worker goroutines (default: GOMAXPROCS)")
// 	procs      = flag.Int("bench.procs", 0, "GOMAXPROCS (default: runtime.NumCPU())")
// 	inner      = flag.Int("bench.inner", 1_000_000, "Inner loop work (fma/sha256/simlike)")
// 	nSize      = flag.Int("bench.n", 160, "Problem size (matmul: NxN)")
// 	dbEvery    = flag.Int("bench.dbEvery", 10_000, "simlike: do a DB ping every N inner iterations (0 = never)")
// 	dbSQL      = flag.String("bench.sql", "SELECT 1", "simlike: SQL to run when DB step triggers")
// 	logPerSec  = flag.Bool("bench.logPerSec", true, "Print per-second throughput lines")
// 	httpListen = flag.String("listen", "", "Optional additional HTTP listener addr in bench mode (e.g. :8080)")
// 	flushEvery = flag.Int("bench.flush", 1024, "Flush local ops to global counter every N ops (lower = more accurate per-second, higher = less overhead)")
// )

// func main() {
// 	flag.Parse()

// 	if *procs <= 0 {
// 		*procs = runtime.NumCPU()
// 	}
// 	runtime.GOMAXPROCS(*procs)

// 	// DB used by server and simlike kernel
// 	db := config.ConnectDB()
// 	db.SetMaxOpenConns(15)
// 	db.SetMaxIdleConns(10)
// 	db.SetConnMaxLifetime(5 * time.Minute)
// 	defer db.Close()

// 	if *benchMode {
// 		if *httpListen != "" {
// 			go func() {
// 				s := api.NewAPIServer(*httpListen, db)
// 				if err := s.Run(); err != nil {
// 					log.Printf("bench-mode HTTP listener error: %v", err)
// 				}
// 			}()
// 		}
// 		runBenchmark(db)
// 		return
// 	}

// 	// Normal server mode
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}
// 	server := api.NewAPIServer(":"+port, db)
// 	if err := server.Run(); err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}
// }

// // ------------------------ Benchmark ------------------------

// type benchFn func(workerID int) int

// func runBenchmark(db *sql.DB) {
// 	if *workers <= 0 {
// 		*workers = runtime.GOMAXPROCS(0)
// 	}

// 	fn := selectKernel(*kernel, *inner, *nSize, db, *dbEvery, *dbSQL)

// 	ctx, cancel := context.WithTimeout(context.Background(), *duration)
// 	defer cancel()

// 	var total atomic.Uint64
// 	var wg sync.WaitGroup
// 	startWall := time.Now()

// 	fmt.Println("=== CPU Benchmark ===")
// 	fmt.Printf("Go: %s | GOMAXPROCS: %d | HW CPUs: %d\n", runtime.Version(), runtime.GOMAXPROCS(0), runtime.NumCPU())
// 	fmt.Printf("Kernel=%s | duration=%s | workers=%d | inner=%d | n=%d | dbEvery=%d | flushEvery=%d\n",
// 		*kernel, duration.String(), *workers, *inner, *nSize, *dbEvery, *flushEvery)
// 	fmt.Println("---------------------")

// 	flushN := uint64(1)
// 	if *flushEvery > 0 {
// 		flushN = uint64(*flushEvery)
// 	}

// 	for i := 0; i < *workers; i++ {
// 		wg.Add(1)
// 		go func(id int) {
// 			defer wg.Done()
// 			local := uint64(0)
// 			for {
// 				select {
// 				case <-ctx.Done():
// 					if local > 0 {
// 						total.Add(local)
// 					}
// 					return
// 				default:
// 					ops := fn(id)
// 					local += uint64(ops)
// 					if local >= flushN {
// 						total.Add(local)
// 						local = 0
// 					}
// 				}
// 			}
// 		}(i)
// 	}

// 	var last uint64
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// loop:
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			break loop
// 		case <-ticker.C:
// 			if *logPerSec {
// 				cur := total.Load()
// 				delta := cur - last
// 				last = cur
// 				fmt.Printf("t+%2ds: %9d ops/sec\n", int(time.Since(startWall).Seconds()), delta)
// 			}
// 		}
// 	}
// 	wg.Wait()

// 	elapsed := time.Since(startWall)
// 	final := total.Load()
// 	opsPerSec := float64(final) / elapsed.Seconds()
// 	fmt.Println("---------------------")
// 	fmt.Printf("Elapsed: %v | Total ops: %d | Avg throughput: %.0f ops/sec\n", elapsed, final, opsPerSec)
// }

// func selectKernel(name string, inner, n int, db *sql.DB, dbEvery int, sqlText string) benchFn {
// 	switch name {
// 	case "fma":
// 		return fmaKernel(inner)
// 	case "sha256":
// 		return sha256Kernel(inner)
// 	case "matmul":
// 		return matmulKernel(n)
// 	case "simlike":
// 		return simlikeKernel(inner, dbEvery, db, sqlText)
// 	default:
// 		log.Fatalf("unknown -bench.kernel=%q (use fma|sha256|matmul|simlike)", name)
// 		return nil
// 	}
// }

// // ---------- Kernels ----------

// func fmaKernel(inner int) benchFn {
// 	return func(workerID int) int {
// 		a, b, c, d := 1.0000001, 1.0000003, 0.9999997, 0.5000001
// 		e, f := 1.41421356, 3.14159265
// 		for i := 0; i < inner; i++ {
// 			a = a*b + c
// 			d = d*e + f
// 			b = math.Sin(a) + 1.00000001
// 			c = math.Cos(d) + 0.99999999
// 		}
// 		if a+d == 0 {
// 			log.Print("")
// 		}
// 		return 1
// 	}
// }
// func sha256Kernel(inner int) benchFn {
// 	bufs := sync.Pool{
// 		New: func() any {
// 			b := make([]byte, 64<<10)
// 			rand.Read(b)
// 			return b
// 		},
// 	}
// 	return func(workerID int) int {
// 		data := bufs.Get().([]byte)
// 		defer bufs.Put(data)
// 		chunk := 1000 // report every 1,000 iterations
// 		totalOps := 0
// 		for i := 0; i < inner; i++ {
// 			h := sha256.Sum256(data)
// 			copy(data[:32], h[:])
// 			totalOps++
// 			if totalOps%chunk == 0 {
// 				return chunk // flush partial progress
// 			}
// 		}
// 		return totalOps % chunk
// 	}
// }

// func matmulKernel(n int) benchFn {
// 	type mats struct {
// 		A, B, C [][]float64
// 		N       int
// 	}
// 	newMats := func(n int) mats {
// 		A := make([][]float64, n)
// 		B := make([][]float64, n)
// 		C := make([][]float64, n)
// 		for i := 0; i < n; i++ {
// 			A[i] = make([]float64, n)
// 			B[i] = make([]float64, n)
// 			C[i] = make([]float64, n)
// 			for j := 0; j < n; j++ {
// 				A[i][j] = rand.Float64()
// 				B[i][j] = rand.Float64()
// 			}
// 		}
// 		return mats{A: A, B: B, C: C, N: n}
// 	}
// 	var pool sync.Pool
// 	pool.New = func() any { return newMats(n) }

// 	return func(workerID int) int {
// 		m := pool.Get().(mats)
// 		N := m.N
// 		for i := 0; i < N; i++ {
// 			rowA := m.A[i]
// 			rowC := m.C[i]
// 			for j := 0; j < N; j++ {
// 				sum := 0.0
// 				for k := 0; k < N; k++ {
// 					sum += rowA[k] * m.B[k][j]
// 				}
// 				rowC[j] = sum
// 			}
// 		}
// 		if m.C[0][0] == -1 {
// 			log.Print("")
// 		}
// 		pool.Put(m)
// 		return 1
// 	}
// }

// func simlikeKernel(inner, dbEvery int, db *sql.DB, sqlText string) benchFn {
// 	if dbEvery < 0 {
// 		dbEvery = 0
// 	}
// 	return func(workerID int) int {
// 		a, b, c, d := 1.0000001, 1.0000003, 0.9999997, 0.5000001
// 		e, f := 1.41421356, 3.14159265
// 		for i := 1; i <= inner; i++ {
// 			a = a*b + c
// 			d = d*e + f
// 			b = math.Sin(a) + 1.00000001
// 			c = math.Cos(d) + 0.99999999

// 			if dbEvery > 0 && i%dbEvery == 0 && db != nil {
// 				if err := pingSQL(db, sqlText); err != nil {
// 					log.Printf("[simlike] DB step error: %v", err)
// 				}
// 			}
// 		}
// 		if a+d == 0 {
// 			log.Print("")
// 		}
// 		return 1
// 	}
// }

// func pingSQL(db *sql.DB, sqlText string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()
// 	var one int
// 	return db.QueryRowContext(ctx, sqlText).Scan(&one)
// }
