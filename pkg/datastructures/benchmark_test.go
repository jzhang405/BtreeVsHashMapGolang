package datastructures

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// 基准测试通用参数
const (
	benchmarkSize = 100000
	smallSize     = 10000
)

// 通用比较函数
var intComparator Comparator = func(a, b interface{}) int {
	switch a.(type) {
	case int:
		if a.(int) < b.(int) {
			return -1
		} else if a.(int) > b.(int) {
			return 1
		}
	case int64:
		if a.(int64) < b.(int64) {
			return -1
		} else if a.(int64) > b.(int64) {
			return 1
		}
	}
	return 0
}

// 生成测试数据
func generateTestData(size int) []int {
	rand.Seed(time.Now().UnixNano())
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Int()
	}
	return data
}

// =============== B+树基准测试 ===============

func BenchmarkBPlusTreeInsert(b *testing.B) {
	data := generateTestData(benchmarkSize)
	tree := NewBPlusTree(64, intComparator)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree = NewBPlusTree(64, intComparator)
		b.StopTimer()
		keys := data[:b.N%benchmarkSize]
		b.StartTimer()

		for j, key := range keys {
			tree.Insert(key, fmt.Sprintf("value_%d", j))
		}
	}
}

func BenchmarkBPlusTreeSearch(b *testing.B) {
	tree := NewBPlusTree(64, intComparator)
	data := generateTestData(benchmarkSize)

	// 预填充数据
	for i, key := range data {
		tree.Insert(key, fmt.Sprintf("value_%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := data[i%benchmarkSize]
		tree.Search(key)
	}
}

func BenchmarkBPlusTreeRangeQuery(b *testing.B) {
	tree := NewBPlusTree(64, intComparator)
	data := generateTestData(benchmarkSize)

	// 预填充数据并排序
	for i, key := range data {
		tree.Insert(key, fmt.Sprintf("value_%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := data[i%benchmarkSize]
		end := start + 1000
		tree.RangeQuery(start, end)
	}
}

func BenchmarkBPlusTreeDelete(b *testing.B) {
	tree := NewBPlusTree(64, intComparator)
	data := generateTestData(benchmarkSize)

	// 预填充数据
	for i, key := range data {
		tree.Insert(key, fmt.Sprintf("value_%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := data[i%benchmarkSize]
		tree.Delete(key)
	}
}

// =============== 跳表基准测试 ===============

func BenchmarkSkipListInsert(b *testing.B) {
	data := generateTestData(benchmarkSize)
	skipList := NewDefaultSkipList(intComparator)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		skipList = NewDefaultSkipList(intComparator)
		b.StopTimer()
		keys := data[:b.N%benchmarkSize]
		b.StartTimer()

		for j, key := range keys {
			skipList.Insert(key, fmt.Sprintf("value_%d", j))
		}
	}
}

func BenchmarkSkipListSearch(b *testing.B) {
	skipList := NewDefaultSkipList(intComparator)
	data := generateTestData(benchmarkSize)

	// 预填充数据
	for i, key := range data {
		skipList.Insert(key, fmt.Sprintf("value_%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := data[i%benchmarkSize]
		skipList.Search(key)
	}
}

func BenchmarkSkipListRangeQuery(b *testing.B) {
	skipList := NewDefaultSkipList(intComparator)
	data := generateTestData(benchmarkSize)

	// 预填充数据
	for i, key := range data {
		skipList.Insert(key, fmt.Sprintf("value_%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := data[i%benchmarkSize]
		end := start + 1000
		skipList.RangeQuery(start, end)
	}
}

func BenchmarkSkipListDelete(b *testing.B) {
	skipList := NewDefaultSkipList(intComparator)
	data := generateTestData(benchmarkSize)

	// 预填充数据
	for i, key := range data {
		skipList.Insert(key, fmt.Sprintf("value_%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := data[i%benchmarkSize]
		skipList.Delete(key)
	}
}

// =============== 可扩展哈希基准测试 ===============

func BenchmarkExtendibleHashInsert(b *testing.B) {
	data := generateTestData(benchmarkSize)
	hashTable := NewExtendibleHashWithDefault()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hashTable = NewExtendibleHashWithDefault()
		b.StopTimer()
		keys := data[:b.N%benchmarkSize]
		b.StartTimer()

		for j, key := range keys {
			hashTable.Insert(key, fmt.Sprintf("value_%d", j))
		}
	}
}

func BenchmarkExtendibleHashSearch(b *testing.B) {
	hashTable := NewExtendibleHashWithDefault()
	data := generateTestData(benchmarkSize)

	// 预填充数据
	for i, key := range data {
		hashTable.Insert(key, fmt.Sprintf("value_%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := data[i%benchmarkSize]
		hashTable.Search(key)
	}
}

func BenchmarkExtendibleHashDelete(b *testing.B) {
	hashTable := NewExtendibleHashWithDefault()
	data := generateTestData(benchmarkSize)

	// 预填充数据
	for i, key := range data {
		hashTable.Insert(key, fmt.Sprintf("value_%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := data[i%benchmarkSize]
		hashTable.Delete(key)
	}
}

// =============== 布隆过滤器基准测试 ===============

func BenchmarkBloomFilterInsert(b *testing.B) {
	data := generateTestData(benchmarkSize)
	bloomFilter := NewDefaultBloomFilter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bloomFilter = NewDefaultBloomFilter()
		b.StopTimer()
		keys := data[:b.N%benchmarkSize]
		b.StartTimer()

		for _, key := range keys {
			bloomFilter.AddInt(key)
		}
	}
}

func BenchmarkBloomFilterSearch(b *testing.B) {
	bloomFilter := NewDefaultBloomFilter()
	data := generateTestData(benchmarkSize)

	// 预填充数据
	for _, key := range data {
		bloomFilter.AddInt(key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := data[i%benchmarkSize]
		bloomFilter.ContainsInt(key)
	}
}

// =============== 性能对比测试 ===============

// TestPerformanceComparison 性能对比测试
func TestPerformanceComparison(t *testing.T) {
	benchmarkSizes := []int{1000, 10000, 100000}

	for _, size := range benchmarkSizes {
		t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
			data := generateTestData(size)

			// B+树测试
			tree := NewBPlusTree(64, intComparator)
			start := time.Now()
			for i, key := range data {
				tree.Insert(key, fmt.Sprintf("value_%d", i))
			}
			btreeInsertTime := time.Since(start)

			start = time.Now()
			for _, key := range data {
				tree.Search(key)
			}
			btreeSearchTime := time.Since(start)

			start = time.Now()
			for i := 0; i < 100 && i < size; i++ {
				tree.RangeQuery(data[i], data[i]+1000)
			}
			btreeRangeTime := time.Since(start)

			// 跳表测试
			skipList := NewDefaultSkipList(intComparator)
			start = time.Now()
			for i, key := range data {
				skipList.Insert(key, fmt.Sprintf("value_%d", i))
			}
			skipListInsertTime := time.Since(start)

			start = time.Now()
			for _, key := range data {
				skipList.Search(key)
			}
			skipListSearchTime := time.Since(start)

			start = time.Now()
			for i := 0; i < 100 && i < size; i++ {
				skipList.RangeQuery(data[i], data[i]+1000)
			}
			skipListRangeTime := time.Since(start)

			// 可扩展哈希测试
			hashTable := NewExtendibleHashWithDefault()
			start = time.Now()
			for i, key := range data {
				hashTable.Insert(key, fmt.Sprintf("value_%d", i))
			}
			hashInsertTime := time.Since(start)

			start = time.Now()
			for _, key := range data {
				hashTable.Search(key)
			}
			hashSearchTime := time.Since(start)

			t.Logf("\n========== Size: %d ==========", size)
			t.Logf("B+Tree       - Insert: %v, Search: %v, Range: %v",
				btreeInsertTime, btreeSearchTime, btreeRangeTime)
			t.Logf("SkipList     - Insert: %v, Search: %v, Range: %v",
				skipListInsertTime, skipListSearchTime, skipListRangeTime)
			t.Logf("ExtHash      - Insert: %v, Search: %v",
				hashInsertTime, hashSearchTime)
		})
	}
}

// TestMemoryUsage 内存使用测试
func TestMemoryUsage(t *testing.T) {
	size := 100000
	data := generateTestData(size)

	// 测试B+树内存
	tree := NewBPlusTree(64, intComparator)
	for i, key := range data {
		tree.Insert(key, fmt.Sprintf("value_%d", i))
	}

	// 测试跳表内存
	skipList := NewDefaultSkipList(intComparator)
	for i, key := range data {
		skipList.Insert(key, fmt.Sprintf("value_%d", i))
	}

	// 测试可扩展哈希内存
	hashTable := NewExtendibleHashWithDefault()
	for i, key := range data {
		hashTable.Insert(key, fmt.Sprintf("value_%d", i))
	}

	t.Logf("B+Tree Height: %d", tree.Height())
	t.Logf("SkipList Level: %d, Height: %d", skipList.Level(), skipList.Height())
	t.Logf("ExtHash GlobalDepth: %d, BucketCount: %d", hashTable.GlobalDepth(), hashTable.BucketCount())
}

// TestConcurrency 并发性能测试
func TestConcurrency(t *testing.T) {
	size := 10000
	data := generateTestData(size)

	// 测试B+树并发
	tree := NewBPlusTree(64, intComparator)
	testConcurrentInsert(tree, data, t)

	// 测试跳表并发
	skipList := NewDefaultSkipList(intComparator)
	testConcurrentInsert(skipList, data, t)

	// 测试可扩展哈希并发
	hashTable := NewExtendibleHashWithDefault()
	testConcurrentInsert(hashTable, data, t)
}

func testConcurrentInsert(ds interface{}, data []int, t *testing.T) {
	// 并发插入测试
	done := make(chan bool)
	threads := 4
	chunkSize := len(data) / threads

	for i := 0; i < threads; i++ {
		go func(start, end int) {
			for j := start; j < end; j++ {
				key := data[j]
				switch tree := ds.(type) {
				case *BPlusTree:
					tree.Insert(key, fmt.Sprintf("value_%d", j))
				case *SkipList:
					tree.Insert(key, fmt.Sprintf("value_%d", j))
				case *ExtendibleHash:
					tree.Insert(key, fmt.Sprintf("value_%d", j))
				}
			}
			done <- true
		}(i*chunkSize, (i+1)*chunkSize)
	}

	for i := 0; i < threads; i++ {
		<-done
	}

	t.Logf("%T completed concurrent insert test", ds)
}

// TestDataStructureFeatures 功能特性测试
func TestDataStructureFeatures(t *testing.T) {
	data := generateTestData(100)

	// 测试B+树范围查询
	tree := NewBPlusTree(64, intComparator)
	for i, key := range data {
		tree.Insert(key, fmt.Sprintf("value_%d", i))
	}

	rangeData, err := tree.RangeQuery(data[10], data[20])
	if err != nil {
		t.Errorf("B+Tree range query failed: %v", err)
	}
	t.Logf("B+Tree range query returned %d items", len(rangeData))

	// 测试跳表范围查询
	skipList := NewDefaultSkipList(intComparator)
	for i, key := range data {
		skipList.Insert(key, fmt.Sprintf("value_%d", i))
	}

	rangeData, err = skipList.RangeQuery(data[10], data[20])
	if err != nil {
		t.Errorf("SkipList range query failed: %v", err)
	}
	t.Logf("SkipList range query returned %d items", len(rangeData))

	// 测试布隆过滤器
	bloomFilter := NewDefaultBloomFilter()
	for _, key := range data {
		bloomFilter.AddInt(key)
	}

	// 验证已存在的元素
	existCount := 0
	for _, key := range data {
		if bloomFilter.ContainsInt(key) {
			existCount++
		}
	}
	t.Logf("BloomFilter found %d/%d existing elements", existCount, len(data))

	// 测试默克尔树
	kvData := make([]KeyValue, len(data))
	for i, key := range data {
		kvData[i] = KeyValue{Key: key, Value: fmt.Sprintf("value_%d", i)}
	}

	merkleTree := NewMerkleTreeFromKV(kvData)
	rootHash := merkleTree.GetRootHash()

	// 验证数据完整性
	for i, key := range data {
		if !merkleTree.VerifyData(i, []byte(fmt.Sprintf("%v:value_%d", key, i))) {
			t.Errorf("MerkleTree verification failed for key %d", key)
		}
	}
	t.Logf("MerkleTree root hash: %s", rootHash)
}

// TestCorrectness 正确性测试
func TestCorrectness(t *testing.T) {
	data := []int{5, 3, 7, 1, 9, 4, 6, 8, 2}

	// 测试B+树
	tree := NewBPlusTree(4, intComparator)
	for _, key := range data {
		tree.Insert(key, fmt.Sprintf("value_%d", key))
	}

	// 验证所有键
	for _, key := range data {
		value, found := tree.Search(key)
		if !found {
			t.Errorf("B+Tree: key %d not found", key)
		} else {
			t.Logf("B+Tree: found key %d with value %v", key, value)
		}
	}

	// 测试删除
	if !tree.Delete(5) {
		t.Errorf("B+Tree: failed to delete key 5")
	}
	if value, _ := tree.Search(5); value != nil {
		t.Errorf("B+Tree: key 5 still exists after deletion")
	}

	// 测试跳表
	skipList := NewDefaultSkipList(intComparator)
	for _, key := range data {
		skipList.Insert(key, fmt.Sprintf("value_%d", key))
	}

	// 验证所有键
	for _, key := range data {
		value, found := skipList.Search(key)
		if !found {
			t.Errorf("SkipList: key %d not found", key)
		} else {
			t.Logf("SkipList: found key %d with value %v", key, value)
		}
	}

	// 测试可扩展哈希
	hashTable := NewExtendibleHashWithDefault()
	for _, key := range data {
		hashTable.Insert(key, fmt.Sprintf("value_%d", key))
	}

	// 验证所有键
	for _, key := range data {
		value, found := hashTable.Search(key)
		if !found {
			t.Errorf("ExtendibleHash: key %d not found", key)
		} else {
			t.Logf("ExtendibleHash: found key %d with value %v", key, value)
		}
	}
}

// BenchmarkComprehensive 全面性能基准测试
func BenchmarkComprehensive(b *testing.B) {
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			data := generateTestData(size)

			// B+树基准测试
			tree := NewBPlusTree(64, intComparator)
			b.Run("BPlusTree_Insert", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					tree = NewBPlusTree(64, intComparator)
					for j := 0; j < size && j < b.N; j++ {
						tree.Insert(data[j%size], fmt.Sprintf("value_%d", j))
					}
				}
			})

			skipList := NewDefaultSkipList(intComparator)
			b.Run("SkipList_Insert", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					skipList = NewDefaultSkipList(intComparator)
					for j := 0; j < size && j < b.N; j++ {
						skipList.Insert(data[j%size], fmt.Sprintf("value_%d", j))
					}
				}
			})

			hashTable := NewExtendibleHashWithDefault()
			b.Run("ExtendibleHash_Insert", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					hashTable = NewExtendibleHashWithDefault()
					for j := 0; j < size && j < b.N; j++ {
						hashTable.Insert(data[j%size], fmt.Sprintf("value_%d", j))
					}
				}
			})
		})
	}
}
