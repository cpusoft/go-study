package main

import (
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/badger/v4/options"
	"github.com/panjf2000/ants/v2"
)

// 全局DB实例（高并发复用）
var db *badger.DB
var (
	writeCount uint64 // 原子计数：写入总数
	readCount  uint64 // 原子计数：读取总数
)

// 初始化高并发配置的Badger
func initDB() error {
	opts := badger.DefaultOptions("./badger_high_concurrency")
	// 高并发核心配置
	opts.MemTableSize = 256 * 1024 * 1024 // 内存表256MB
	opts.NumMemtables = runtime.NumCPU()  // 内存表缓冲写操作
	opts.ValueLogFileSize = 1 << 30       // Value日志1GB
	opts.NumCompactors = runtime.NumCPU() // 压缩器数=CPU核心数
	opts.SyncWrites = false               // 关闭同步写，提升写性能
	opts.Compression = options.None       // 关闭压缩，减少CPU开销
	opts.NumLevelZeroTables = 10          // 增大L0表阈值，减少压缩触发
	opts.NumLevelZeroTablesStall = 20     // 增大stall阈值，避免写阻塞

	var err error
	db, err = badger.Open(opts)
	return err
}

// 批量写入（高并发写核心）
func batchWrite(prefix string, start, end int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Batch批量写：减少事务提交开销
	err := db.Batch(func(b *badger.Batch) error {
		for i := start; i < end; i++ {
			key := []byte(prefix + ":" + string(rune(i)))
			value := []byte("data:" + string(rune(i)))
			if err := b.Set(key, value); err != nil {
				return err
			}
			atomic.AddUint64(&writeCount, 1) // 原子计数，避免并发冲突
		}
		return nil
	})

	if err != nil {
		log.Printf("批量写入失败: %v", err)
	}
}

// 并发读取
func concurrentRead(prefix string, start, end int, wg *sync.WaitGroup) {
	defer wg.Done()

	// View只读事务：轻量级，无锁
	err := db.View(func(txn *badger.Txn) error {
		for i := start; i < end; i++ {
			key := []byte(prefix + ":" + string(rune(i)))
			item, err := txn.Get(key)
			if err != nil {
				if err == badger.ErrKeyNotFound {
					return nil
				}
				return err
			}
			// 按需读取Value（不拷贝，直接使用）
			_ = item.Value(func(v []byte) error {
				// 业务逻辑处理...
				return nil
			})
			atomic.AddUint64(&readCount, 1)
		}
		return nil
	})

	if err != nil {
		log.Printf("并发读取失败: %v", err)
	}
}

func main() {
	// 1. 初始化DB
	if err := initDB(); err != nil {
		log.Fatalf("初始化DB失败: %v", err)
	}
	defer db.Close()

	// 2. 配置并发参数
	totalKeys := 100000 // 总写入Key数
	goroutineNum := 10  // 并发goroutine数
	keysPerGoroutine := totalKeys / goroutineNum

	// 3. 高并发写入
	var writeWg sync.WaitGroup
	writeWg.Add(goroutineNum)
	writeStart := time.Now()

	// goroutine池：限制并发数
	pool, _ := ants.NewPool(goroutineNum)
	for i := 0; i < goroutineNum; i++ {
		start := i * keysPerGoroutine
		end := start + keysPerGoroutine
		if i == goroutineNum-1 {
			end = totalKeys // 最后一个goroutine处理剩余Key
		}
		_ = pool.Submit(func() {
			batchWrite("user", start, end, &writeWg)
		})
	}
	writeWg.Wait()
	pool.Release()

	// 打印写入性能
	writeDur := time.Since(writeStart)
	log.Printf("写入完成：%d条 | 耗时%v | 吞吐量%.2f条/秒",
		writeCount, writeDur, float64(writeCount)/writeDur.Seconds())

	// 4. 高并发读取
	var readWg sync.WaitGroup
	readWg.Add(goroutineNum)
	readStart := time.Now()

	readPool, _ := ants.NewPool(goroutineNum)
	for i := 0; i < goroutineNum; i++ {
		start := i * keysPerGoroutine
		end := start + keysPerGoroutine
		if i == goroutineNum-1 {
			end = totalKeys
		}
		_ = readPool.Submit(func() {
			concurrentRead("user", start, end, &readWg)
		})
	}
	readWg.Wait()
	readPool.Release()

	// 打印读取性能
	readDur := time.Since(readStart)
	log.Printf("读取完成：%d条 | 耗时%v | 吞吐量%.2f条/秒",
		readCount, readDur, float64(readCount)/readDur.Seconds())
}
