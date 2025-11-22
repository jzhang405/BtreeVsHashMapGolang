module github.com/datastructures/bplus-vs-hash

go 1.21

// 这是一个演示项目，展示了文档《是否存在一种数据结构能同时拥有btree和hash表的所有优点？》中提到的各种数据结构
// 的生产级Go语言实现。
//
// 主要包含以下数据结构：
// - B+树：适用于磁盘存储和范围查询
// - 跳表：融合型结构，内存友好
// - 默克尔树：适用于区块链和分布式存储
// - 可扩展哈希：磁盘优化的哈希表
// - 布隆过滤器：概率性数据结构
// - 优化的B+树：结合B+树和哈希的优点
//
// 使用示例：
//
// 	package main
//
// 	import (
// 		"fmt"
// 		"github.com/datastructures/bplus-vs-hash/pkg/datastructures"
// 	)
//
// 	func main() {
// 		// 创建B+树
// 		tree := datastructures.NewBPlusTree(64, func(a, b interface{}) int {
// 			if a.(int) < b.(int) {
// 				return -1
// 			} else if a.(int) > b.(int) {
// 				return 1
// 			}
// 			return 0
// 		})
//
// 		// 插入数据
// 		tree.Insert(100, "value_100")
// 		tree.Insert(200, "value_200")
//
// 		// 查询数据
// 		if value, found := tree.Search(100); found {
// 			fmt.Println("Found:", value)
// 		}
//
// 		// 范围查询
// 		results, _ := tree.RangeQuery(50, 150)
// 		for _, kv := range results {
// 			fmt.Printf("Key: %v, Value: %v\n", kv.Key, kv.Value)
// 		}
// 	}
//
// 运行测试：
//
// 	go test -v -run TestCorrectness
// 	go test -bench=. -benchmem
//
// 参考文档：
// - README.md: 完整的使用指南和API参考
// - examples/example_usage.go: 详细的使用示例
// - benchmark_test.go: 性能测试和对比

require ()
