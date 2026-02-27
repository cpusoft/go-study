package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/badger/v4/options"
)

type BadgerManager struct {
	db *badger.DB
}

func NewBadgerManager(dbPath string) (*BadgerManager, error) {
	// 优化配置以适应高并发场景
	opts := badger.DefaultOptions(dbPath)

	opts = opts.WithMemTableSize(256 * 1024 * 1024) // 128MB内存表, <=小于系统内存/4
	opts = opts.WithNumMemtables(runtime.NumCPU())  // cpunum个内存表
	opts = opts.WithValueLogFileSize(1 << 30)       // 1G日志文件
	opts = opts.WithNumCompactors(runtime.NumCPU()) // 增加压缩器数量
	opts = opts.WithSyncWrites(false)               // 关闭同步写，提升性能
	opts = opts.WithCompression(options.None)       // 关闭压缩，减少CPU开销
	opts = opts.WithNumLevelZeroTables(10)          // 增大L0表阈值，减少压缩触发
	opts = opts.WithNumLevelZeroTablesStall(20)     // 增大stall阈值，避免写阻塞

	opts = opts.WithNumGoroutines(runtime.NumCPU()) // 增加并发goroutine数量
	opts = opts.WithBlockCacheSize(100 << 20)       // 100MB块缓存
	opts = opts.WithIndexCacheSize(50 << 20)        // 50MB索引缓存
	opts = opts.WithValueThreshold(1024)            // 1KB以下的值内联存储

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &BadgerManager{db: db}, nil
}

func (bm *BadgerManager) Close() {
	bm.db.Close()
}

// 批量写入 - 最高效的写入方式
func (bm *BadgerManager) BatchWriteConcurrent(data map[string]string) error {
	// 创建批量写入器
	batch := bm.db.NewWriteBatch()
	defer batch.Cancel() // 确保在出错时取消

	for key, value := range data {
		if err := batch.Set([]byte(key), []byte(value)); err != nil {
			return err
		}
	}

	// 刷新批量写入
	return batch.Flush()
}

// 并发批量写入示例
func (bm *BadgerManager) ConcurrentBatchWrites() {
	var wg sync.WaitGroup
	workerCount := 10

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			data := make(map[string]string)
			for j := 0; j < 1000; j++ {
				key := fmt.Sprintf("worker%d_key%d", workerID, j)
				value := fmt.Sprintf("value_%d_%d", workerID, j)
				data[key] = value
			}

			if err := bm.BatchWriteConcurrent(data); err != nil {
				log.Printf("Worker %d error: %v", workerID, err)
			} else {
				log.Printf("Worker %d completed successfully", workerID)
			}
		}(i)
	}

	wg.Wait()
}

// 带重试机制的事务处理
func (bm *BadgerManager) UpdateWithRetry(key string, updateFunc func(current []byte) []byte, maxRetries int) error {
	for retry := 0; retry < maxRetries; retry++ {
		err := bm.db.Update(func(txn *badger.Txn) error {
			// 读取当前值
			item, err := txn.Get([]byte(key))
			if err != nil && err != badger.ErrKeyNotFound {
				return err
			}

			var currentValue []byte
			if err == nil {
				currentValue, err = item.ValueCopy(nil)
				if err != nil {
					return err
				}
			}

			// 应用更新函数
			newValue := updateFunc(currentValue)

			// 写入新值
			return txn.Set([]byte(key), newValue)
		})

		if err == nil {
			return nil // 成功
		}

		if err == badger.ErrConflict {
			// 事务冲突，等待后重试
			time.Sleep(time.Duration(retry) * 10 * time.Millisecond)
			continue
		}

		return err // 其他错误
	}

	return fmt.Errorf("max retries exceeded")
}

// 并发事务处理示例
func (bm *BadgerManager) ConcurrentTransactionalUpdates() {
	const key = "concurrent_counter"
	var wg sync.WaitGroup
	workerCount := 20

	// 初始化计数器
	bm.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte("0"))
	})

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			updateFunc := func(current []byte) []byte {
				var count int
				if current != nil {
					fmt.Sscanf(string(current), "%d", &count)
				}
				count++
				return []byte(fmt.Sprintf("%d", count))
			}

			if err := bm.UpdateWithRetry(key, updateFunc, 5); err != nil {
				log.Printf("Worker %d failed: %v", workerID, err)
			}
		}(i)
	}

	wg.Wait()

	// 验证结果
	bm.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, _ := item.ValueCopy(nil)
		log.Printf("Final counter value: %s", val)
		return nil
	})
}

// 高并发只读查询
func (bm *BadgerManager) ConcurrentReads(keys []string, resultChan chan<- map[string]string) {
	var wg sync.WaitGroup
	batchSize := len(keys) / 10 // 分批处理

	for i := 0; i < len(keys); i += batchSize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()

			end := start + batchSize
			if end > len(keys) {
				end = len(keys)
			}

			batchKeys := keys[start:end]
			results := make(map[string]string)

			// 每个goroutine使用独立的事务
			err := bm.db.View(func(txn *badger.Txn) error {
				for _, key := range batchKeys {
					item, err := txn.Get([]byte(key))
					if err != nil {
						if err == badger.ErrKeyNotFound {
							results[key] = ""
							continue
						}
						return err
					}

					err = item.Value(func(val []byte) error {
						results[key] = string(val)
						return nil
					})
					if err != nil {
						return err
					}
				}
				return nil
			})

			if err == nil {
				// 将结果发送到通道
				for k, v := range results {
					resultChan <- map[string]string{k: v}
				}
			}
		}(i)
	}

	wg.Wait()
	close(resultChan)
}

// 并行迭代处理
func (bm *BadgerManager) ParallelIteration(prefix string, processor func(key, value string) error) error {
	const numWorkers = 8
	var wg sync.WaitGroup
	workChan := make(chan []string, numWorkers)
	errChan := make(chan error, 1)

	// 启动worker
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range workChan {
				for _, key := range batch {
					err := bm.db.View(func(txn *badger.Txn) error {
						item, err := txn.Get([]byte(key))
						if err != nil {
							return err
						}
						return item.Value(func(val []byte) error {
							return processor(key, string(val))
						})
					})
					if err != nil {
						select {
						case errChan <- err:
						default:
						}
						return
					}
				}
			}
		}()
	}

	// 收集所有键
	var allKeys []string
	err := bm.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // 只获取键
		it := txn.NewIterator(opts)
		defer it.Close()

		prefixBytes := []byte(prefix)
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			item := it.Item()
			key := string(item.KeyCopy(nil))
			allKeys = append(allKeys, key)
		}
		return nil
	})

	if err != nil {
		close(workChan)
		return err
	}

	// 分批发送到worker
	batchSize := (len(allKeys) + numWorkers - 1) / numWorkers
	for i := 0; i < len(allKeys); i += batchSize {
		end := i + batchSize
		if end > len(allKeys) {
			end = len(allKeys)
		}
		workChan <- allKeys[i:end]
	}

	close(workChan)
	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}
func main() {
	// 初始化Badger管理器
	manager, err := NewBadgerManager("./badgerdb")
	if err != nil {
		log.Fatal(err)
	}
	defer manager.Close()

	// 示例1: 并发批量写入
	fmt.Println("Starting concurrent batch writes...")
	manager.ConcurrentBatchWrites()

	// 示例2: 并发事务更新
	fmt.Println("Starting concurrent transactional updates...")
	manager.ConcurrentTransactionalUpdates()

	// 示例3: 准备测试数据
	testKeys := make([]string, 100)
	for i := 0; i < 100; i++ {
		testKeys[i] = fmt.Sprintf("test_key_%d", i)
	}

	// 示例4: 并发读取
	resultChan := make(chan map[string]string, 100)
	go manager.ConcurrentReads(testKeys, resultChan)

	// 处理结果
	for result := range resultChan {
		for k, v := range result {
			if v != "" {
				fmt.Printf("Key: %s, Value: %s\n", k, v)
			}
		}
	}

	// 示例5: 并行迭代
	fmt.Println("Starting parallel iteration...")
	manager.ParallelIteration("test_", func(key, value string) error {
		fmt.Printf("Processed: %s -> %s\n", key, value)
		return nil
	})
}
