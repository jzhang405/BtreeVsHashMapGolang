package main

import (
	"fmt"
	"log"

	"github.com/datastructures/bplus-vs-hash/pkg/datastructures"
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
	fmt.Println("=== 数据结构演示 ===\n")

	// 演示B+树
	demoBPlusTree()

	// 演示跳表
	demoSkipList()

	// 演示可扩展哈希
	demoExtendibleHash()

	// 演示布隆过滤器
	demoBloomFilter()

	// 演示默克尔树
	demoMerkleTree()
}

func demoBPlusTree() {
	fmt.Println("--- B+树演示 ---")

	tree := datastructures.NewBPlusTree(64, intComparator)

	// 插入数据
	keys := []int{100, 200, 50, 150, 300, 250}
	for _, key := range keys {
		err := tree.Insert(key, fmt.Sprintf("value_%d", key))
		if err != nil {
			log.Printf("插入键 %d 失败: %v", key, err)
		}
	}

	fmt.Printf("✓ 插入 %d 个键值对\n", len(keys))

	// 等值查询
	value, found := tree.Search(100)
	if found {
		fmt.Printf("✓ 等值查询: Search(100) = %s\n", value)
	}

	// 范围查询
	results, err := tree.RangeQuery(50, 200)
	if err == nil {
		fmt.Printf("✓ 范围查询: RangeQuery(50, 200) 返回 %d 条记录\n", len(results))
	}

	fmt.Printf("✓ 树高度: %d\n", tree.Height())
	fmt.Printf("✓ 数据量: %d\n\n", tree.Size())
}

func demoSkipList() {
	fmt.Println("--- 跳表演示 ---")

	skipList := datastructures.NewDefaultSkipList(intComparator)

	// 插入数据
	keys := []int{100, 200, 50, 150, 300}
	for _, key := range keys {
		err := skipList.Insert(key, fmt.Sprintf("value_%d", key))
		if err != nil {
			log.Printf("插入键 %d 失败: %v", key, err)
		}
	}

	fmt.Printf("✓ 插入 %d 个键值对\n", len(keys))

	// 等值查询
	value, found := skipList.Search(100)
	if found {
		fmt.Printf("✓ 等值查询: Search(100) = %s\n", value)
	}

	// 范围查询
	results, err := skipList.RangeQuery(50, 200)
	if err == nil {
		fmt.Printf("✓ 范围查询: RangeQuery(50, 200) 返回 %d 条记录\n", len(results))
	}

	fmt.Printf("✓ 跳表层数: %d\n", skipList.Level())
	fmt.Printf("✓ 数据量: %d\n\n", skipList.Size())
}

func demoExtendibleHash() {
	fmt.Println("--- 可扩展哈希演示 ---")

	hashTable := datastructures.NewExtendibleHashWithDefault()

	// 插入数据
	keys := []int{100, 200, 50, 150, 300, 400}
	for _, key := range keys {
		err := hashTable.Insert(key, fmt.Sprintf("value_%d", key))
		if err != nil {
			log.Printf("插入键 %d 失败: %v", key, err)
		}
	}

	fmt.Printf("✓ 插入 %d 个键值对\n", len(keys))

	// 等值查询
	value, found := hashTable.Search(100)
	if found {
		fmt.Printf("✓ 等值查询: Search(100) = %s\n", value)
	}

	// 统计信息
	avg, max, min, fullCount := hashTable.GetBucketUsage()
	fmt.Printf("✓ 桶使用统计: 平均=%.2f, 最大=%d, 最小=%d, 满桶数=%d\n", avg, max, min, fullCount)
	fmt.Printf("✓ 全局深度: %d\n", hashTable.GlobalDepth())
	fmt.Printf("✓ 数据量: %d\n\n", hashTable.Size())
}

func demoBloomFilter() {
	fmt.Println("--- 布隆过滤器演示 ---")

	bloomFilter := datastructures.NewDefaultBloomFilter()

	// 添加元素
	keys := []int{100, 200, 50, 150, 300}
	fmt.Println("添加元素:")
	for _, key := range keys {
		bloomFilter.AddInt(key)
		fmt.Printf("  - 添加键 %d\n", key)
	}

	// 检查存在性
	fmt.Println("\n检查元素存在性:")
	for _, key := range []int{100, 150, 999} {
		exists := bloomFilter.ContainsInt(key)
		fmt.Printf("  - 键 %d: %v\n", key, map[bool]string{true: "可能存在", false: "一定不存在"}[exists])
	}

	fpr := bloomFilter.GetFalsePositiveRate()
	fmt.Printf("✓ 当前假阳性率: %.4f%%\n", fpr*100)
	fmt.Printf("✓ 元素数量: %d\n\n", bloomFilter.Size())
}

func demoMerkleTree() {
	fmt.Println("--- 默克尔树演示 ---")

	// 创建键值对数据
	kvs := []datastructures.KeyValue{
		{Key: "file1", Value: "content1"},
		{Key: "file2", Value: "content2"},
		{Key: "file3", Value: "content3"},
		{Key: "file4", Value: "content4"},
	}

	// 创建默克尔树
	merkleTree := datastructures.NewMerkleTreeFromKV(kvs)

	// 获取根哈希
	rootHash := merkleTree.GetRootHash()
	fmt.Printf("✓ 根哈希: %s\n", rootHash)

	// 验证数据完整性
	fmt.Println("\n验证数据完整性:")
	for i, kv := range kvs {
		data := []byte(fmt.Sprintf("%v:%v", kv.Key, kv.Value))
		isValid := merkleTree.VerifyData(i, data)
		fmt.Printf("  - 文件 %s: %s\n", kv.Key, map[bool]string{true: "✓ 有效", false: "✗ 无效"}[isValid])
	}

	// 获取完整性证明
	_, proof, err := merkleTree.GetProof(1)
	if err == nil {
		fmt.Printf("✓ 完整性证明长度: %d\n", len(proof))
	}

	fmt.Printf("✓ 数据块数量: %d\n", merkleTree.Size())
	fmt.Printf("✓ 树高度: %d\n\n", merkleTree.Height())
}
