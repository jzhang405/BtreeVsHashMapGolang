package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yourusername/datastructures/pkg/datastructures"
)

// 整数比较函数
func intComparator(a, b any) int {
	if a.(int) < b.(int) {
		return -1
	} else if a.(int) > b.(int) {
		return 1
	}
	return 0
}

func main() {
	fmt.Println("=== 数据结构使用示例 ===\n")

	// 示例1: B+树基本使用
	exampleBPlusTree()

	// 示例2: 跳表基本使用
	exampleSkipList()

	// 示例3: 可扩展哈希基本使用
	exampleExtendibleHash()

	// 示例4: 布隆过滤器基本使用
	exampleBloomFilter()

	// 示例5: 默克尔树基本使用
	exampleMerkleTree()

	// 示例6: 性能对比
	examplePerformanceComparison()

	// 示例7: 实际应用场景
	exampleRealWorldScenarios()
}

func exampleBPlusTree() {
	fmt.Println("--- 示例1: B+树基本使用 ---")

	// 创建B+树
	tree := datastructures.NewBPlusTree(64, intComparator)

	// 插入数据
	keys := []int{100, 200, 50, 150, 300, 250}
	for _, key := range keys {
		err := tree.Insert(key, fmt.Sprintf("value_%d", key))
		if err != nil {
			log.Printf("插入键 %d 失败: %v", key, err)
		}
	}

	fmt.Printf("插入 %d 个键值对\n", len(keys))

	// 等值查询
	fmt.Println("\n等值查询:")
	for _, key := range []int{100, 150, 999} {
		if value, found := tree.Search(key); found {
			fmt.Printf("  找到键 %d: %s\n", key, value)
		} else {
			fmt.Printf("  键 %d 不存在\n", key)
		}
	}

	// 范围查询
	fmt.Println("\n范围查询 [50, 200):")
	results, err := tree.RangeQuery(50, 200)
	if err != nil {
		log.Printf("范围查询失败: %v", err)
	} else {
		for _, kv := range results {
			fmt.Printf("  键: %v, 值: %s\n", kv.Key, kv.Value)
		}
	}

	// 顺序遍历
	fmt.Println("\n顺序遍历所有数据:")
	allData := tree.ScanAll()
	for _, kv := range allData {
		fmt.Printf("  键: %v, 值: %s\n", kv.Key, kv.Value)
	}

	// 删除数据
	fmt.Println("\n删除键 150:")
	if tree.Delete(150) {
		fmt.Println("  删除成功")
		if _, found := tree.Search(150); !found {
			fmt.Println("  验证: 键已不存在")
		}
	} else {
		fmt.Println("  删除失败")
	}

	fmt.Printf("树的高度: %d\n", tree.Height())
	fmt.Printf("数据量: %d\n\n", tree.Size())
}

func exampleSkipList() {
	fmt.Println("--- 示例2: 跳表基本使用 ---")

	// 创建跳表
	skipList := datastructures.NewDefaultSkipList(intComparator)

	// 插入数据
	keys := []int{100, 200, 50, 150, 300, 250}
	for _, key := range keys {
		err := skipList.Insert(key, fmt.Sprintf("value_%d", key))
		if err != nil {
			log.Printf("插入键 %d 失败: %v", key, err)
		}
	}

	fmt.Printf("插入 %d 个键值对\n", len(keys))

	// 等值查询
	fmt.Println("\n等值查询:")
	for _, key := range []int{100, 150, 999} {
		if value, found := skipList.Search(key); found {
			fmt.Printf("  找到键 %d: %s\n", key, value)
		} else {
			fmt.Printf("  键 %d 不存在\n", key)
		}
	}

	// 范围查询
	fmt.Println("\n范围查询 [50, 200):")
	results, err := skipList.RangeQuery(50, 200)
	if err != nil {
		log.Printf("范围查询失败: %v", err)
	} else {
		for _, kv := range results {
			fmt.Printf("  键: %v, 值: %s\n", kv.Key, kv.Value)
		}
	}

	// 删除数据
	fmt.Println("\n删除键 150:")
	if skipList.Delete(150) {
		fmt.Println("  删除成功")
	}

	fmt.Printf("跳表层数: %d\n", skipList.Level())
	fmt.Printf("数据量: %d\n\n", skipList.Size())
}

func exampleExtendibleHash() {
	fmt.Println("--- 示例3: 可扩展哈希基本使用 ---")

	// 创建可扩展哈希表
	hashTable := datastructures.NewExtendibleHashWithDefault()

	// 插入数据
	keys := []int{100, 200, 50, 150, 300, 250, 400, 450, 500}
	for _, key := range keys {
		err := hashTable.Insert(key, fmt.Sprintf("value_%d", key))
		if err != nil {
			log.Printf("插入键 %d 失败: %v", key, err)
		}
	}

	fmt.Printf("插入 %d 个键值对\n", len(keys))

	// 等值查询
	fmt.Println("\n等值查询:")
	for _, key := range []int{100, 150, 999} {
		if value, found := hashTable.Search(key); found {
			fmt.Printf("  找到键 %d: %s\n", key, value)
		} else {
			fmt.Printf("  键 %d 不存在\n", key)
		}
	}

	// 删除数据
	fmt.Println("\n删除键 150:")
	if hashTable.Delete(150) {
		fmt.Println("  删除成功")
	}

	// 获取统计信息
	avg, max, min, fullCount := hashTable.GetBucketUsage()
	fmt.Printf("\n桶使用统计:\n")
	fmt.Printf("  平均: %.2f, 最大: %d, 最小: %d, 满桶数: %d\n", avg, max, min, fullCount)
	fmt.Printf("  全局深度: %d\n", hashTable.GlobalDepth())
	fmt.Printf("  桶数量: %d\n", hashTable.BucketCount())
	fmt.Printf("  数据量: %d\n\n", hashTable.Size())
}

func exampleBloomFilter() {
	fmt.Println("--- 示例4: 布隆过滤器基本使用 ---")

	// 创建布隆过滤器（10000元素，1%假阳性率）
	bloomFilter := datastructures.NewBloomFilter(10000, 0.01)

	// 添加元素
	keys := []int{100, 200, 50, 150, 300, 250}
	fmt.Println("添加元素:")
	for _, key := range keys {
		bloomFilter.AddInt(key)
		fmt.Printf("  添加键 %d\n", key)
	}

	// 检查存在性
	fmt.Println("\n检查元素存在性:")
	for _, key := range append(keys, 999) {
		exists := bloomFilter.ContainsInt(key)
		fmt.Printf("  键 %d: 可能存在=%v\n", key, exists)
	}

	// 获取假阳性率
	fpr := bloomFilter.GetFalsePositiveRate()
	fmt.Printf("\n当前假阳性率: %.4f%%\n", fpr*100)
	fmt.Printf("元素数量: %d\n\n", bloomFilter.Size())
}

func exampleMerkleTree() {
	fmt.Println("--- 示例5: 默克尔树基本使用 ---")

	// 创建键值对数据
	kvs := []datastructures.KeyValue{
		{Key: "file1.txt", Value: "content1"},
		{Key: "file2.txt", Value: "content2"},
		{Key: "file3.txt", Value: "content3"},
		{Key: "file4.txt", Value: "content4"},
	}

	// 创建默克尔树
	merkleTree := datastructures.NewMerkleTreeFromKV(kvs)

	// 获取根哈希
	rootHash := merkleTree.GetRootHash()
	fmt.Printf("根哈希: %s\n", rootHash)

	// 验证数据完整性
	fmt.Println("\n验证数据完整性:")
	for i, kv := range kvs {
		data := []byte(fmt.Sprintf("%v:%v", kv.Key, kv.Value))
		isValid := merkleTree.VerifyData(i, data)
		fmt.Printf("  文件 %s: %v\n", kv.Key, map[bool]string{true: "有效", false: "无效"}[isValid])
	}

	// 获取完整性证明
	fmt.Println("\n获取完整性证明 (file2.txt):")
	_, proof, err := merkleTree.GetProof(1)
	if err != nil {
		log.Printf("获取证明失败: %v", err)
	} else {
		fmt.Printf("  证明长度: %d\n", len(proof))
		for i, p := range proof {
			fmt.Printf("  步骤 %d: %s\n", i+1, p)
		}
	}

	// 验证证明
	fmt.Println("\n验证证明:")
	data := []byte("file2.txt:content2")
	isValid = datastructures.VerifyProof(data, proof, rootHash)
	fmt.Printf("  证明有效性: %v\n\n", isValid)
}

func examplePerformanceComparison() {
	fmt.Println("--- 示例6: 性能对比 ---")

	// 测试数据
	sizes := []int{1000, 10000, 10000}
	data := generateTestData(sizes[0])

	// 测试B+树
	fmt.Println("\nB+树性能测试:")
	tree := datastructures.NewBPlusTree(64, intComparator)
	start := time.Now()
	for i, key := range data {
		tree.Insert(key, fmt.Sprintf("value_%d", i))
	}
	insertTime := time.Since(start)
	fmt.Printf("  插入时间: %v\n", insertTime)

	start = time.Now()
	for _, key := range data[:1000] {
		tree.Search(key)
	}
	searchTime := time.Since(start)
	fmt.Printf("  1000次查询时间: %v\n", searchTime)

	// 测试跳表
	fmt.Println("\n跳表性能测试:")
	skipList := datastructures.NewDefaultSkipList(intComparator)
	start = time.Now()
	for i, key := range data {
		skipList.Insert(key, fmt.Sprintf("value_%d", i))
	}
	insertTime = time.Since(start)
	fmt.Printf("  插入时间: %v\n", insertTime)

	start = time.Now()
	for _, key := range data[:1000] {
		skipList.Search(key)
	}
	searchTime = time.Since(start)
	fmt.Printf("  1000次查询时间: %v\n", searchTime)

	// 测试可扩展哈希
	fmt.Println("\n可扩展哈希性能测试:")
	hashTable := datastructures.NewExtendibleHashWithDefault()
	start = time.Now()
	for i, key := range data {
		hashTable.Insert(key, fmt.Sprintf("value_%d", i))
	}
	insertTime = time.Since(start)
	fmt.Printf("  插入时间: %v\n", insertTime)

	start = time.Now()
	for _, key := range data[:1000] {
		hashTable.Search(key)
	}
	searchTime = time.Since(start)
	fmt.Printf("  1000次查询时间: %v\n\n", searchTime)
}

func exampleRealWorldScenarios() {
	fmt.Println("--- 示例7: 实际应用场景 ---")

	// 场景1: 用户ID查找（等值查询为主）
	fmt.Println("\n场景1: 用户ID查找系统")
	fmt.Println("推荐: 使用可扩展哈希")
	userCache := datastructures.NewExtendibleHashWithDefault()
	userIDs := []int{1001, 1002, 1003, 1004, 1005}
	for _, id := range userIDs {
		userCache.Insert(id, fmt.Sprintf("User_%d", id))
	}

	if user, found := userCache.Search(1002); found {
		fmt.Printf("  找到用户: %s\n", user)
	}

	// 场景2: 股票价格范围查询（范围查询为主）
	fmt.Println("\n场景2: 股票价格监控系统")
	fmt.Println("推荐: 使用B+树或跳表")
	priceTree := datastructures.NewBPlusTree(64, intComparator)
	prices := []int{100, 105, 110, 95, 120, 115}
	for i, price := range prices {
		priceTree.Insert(price, fmt.Sprintf("Stock_%d", i))
	}

	if results, err := priceTree.RangeQuery(100, 115); err == nil {
		fmt.Printf("  价格范围 [100, 115) 的股票数: %d\n", len(results))
	}

	// 场景3: 数据库索引优化（等值+范围查询）
	fmt.Println("\n场景3: 数据库索引系统")
	fmt.Println("推荐: 使用优化的B+树")
	optTree := datastructures.NewBPlusTreeOptimized(64, 4, intComparator)
	for i := 0; i < 10; i++ {
		optTree.Insert(i*10, fmt.Sprintf("Record_%d", i))
	}

	// 等值查询
	if value, _ := optTree.SearchFast(50); value != nil {
		fmt.Printf("  快速查找键 50: %s\n", value)
	}

	// 范围查询
	if results, err := optTree.RangeQuery(20, 70); err == nil {
		fmt.Printf("  范围查询 [20, 70): %d 条记录\n", len(results))
	}

	// 场景4: 缓存穿透防护
	fmt.Println("\n场景4: 缓存系统")
	fmt.Println("推荐: 使用布隆过滤器 + B+树")
	bf := datastructures.NewBloomFilter(1000, 0.01)
	tree := datastructures.NewBPlusTree(64, intComparator)

	// 预填充缓存
	cacheKeys := []int{1, 2, 3, 4, 5}
	for _, key := range cacheKeys {
		bf.AddInt(key)
		tree.Insert(key, fmt.Sprintf("Cached_Data_%d", key))
	}

	// 查询缓存
	queryKeys := []int{1, 2, 6, 7}
	for _, key := range queryKeys {
		if bf.ContainsInt(key) {
			if value, _ := tree.Search(key); value != nil {
				fmt.Printf("  键 %d: 缓存命中 - %s\n", key, value)
			} else {
				fmt.Printf("  键 %d: 缓存穿透\n", key)
			}
		} else {
			fmt.Printf("  键 %d: 布隆过滤器判断不存在\n", key)
		}
	}

	// 场景5: 分布式存储验证
	fmt.Println("\n场景5: 分布式存储系统")
	fmt.Println("推荐: 使用默克尔树")
	merkleTree := datastructures.NewMerkleTreeFromKV([]datastructures.KeyValue{
		{Key: "chunk1", Value: "data_chunk_1"},
		{Key: "chunk2", Value: "data_chunk_2"},
		{Key: "chunk3", Value: "data_chunk_3"},
	})

	rootHash := merkleTree.GetRootHash()
	fmt.Printf("  存储系统根哈希: %s\n", rootHash)

	// 验证数据完整性
	data := []byte("chunk2:data_chunk_2")
	if hashes, proof, _ := merkleTree.GetProof(1); len(proof) > 0 {
		isValid := datastructures.VerifyProof(data, proof, rootHash)
		fmt.Printf("  数据块验证: %v\n", map[bool]string{true: "通过", false: "失败"}[isValid])
	}
}

func generateTestData(size int) []int {
	rand.Seed(time.Now().UnixNano())
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Int()
	}
	return data
}
