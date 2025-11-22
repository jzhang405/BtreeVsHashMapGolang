package datastructures

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// HashFunc 哈希函数类型
type HashFunc func(data []byte) uint32

// defaultHash 默认哈希函数
func defaultHash(data []byte) uint32 {
	h := fnv.New32a()
	h.Write(data)
	return h.Sum32()
}

// HashBucket 哈希桶
type HashBucket struct {
	keys    []interface{} // 桶中的键
	values  []interface{} // 桶中的值
	localDepth int       // 局部深度
}

// NewHashBucket 创建新的哈希桶
func NewHashBucket() *HashBucket {
	return &HashBucket{
		keys:      make([]interface{}, 0),
		values:    make([]interface{}, 0),
		localDepth: 0,
	}
}

// isFull 检查桶是否已满
func (b *HashBucket) isFull(capacity int) bool {
	return len(b.keys) >= capacity
}

// isEmpty 检查桶是否为空
func (b *HashBucket) isEmpty() bool {
	return len(b.keys) == 0
}

// ExtendibleHash 可扩展哈希表
// 特点：
// - 磁盘优化的哈希表实现
// - O(1)等值查询
// - 动态扩容，无需全量重新哈希
// - 减少随机IO，适合磁盘存储
// - 通过目录和桶的分离实现可扩展性
type ExtendibleHash struct {
	buckets   []*HashBucket // 桶数组（目录）
	directory []*HashBucket // 目录指针
	globalDepth int        // 全局深度
	bucketCapacity  int    // 桶容量
	hashFunc        HashFunc // 哈希函数
	mu              sync.RWMutex // 读写锁
	count           int64   // 总键数
}

// NewExtendibleHash 创建新的可扩展哈希表
// bucketCapacity: 桶容量，建议值：4-64（根据磁盘块大小调整）
// hashFunc: 哈希函数
func NewExtendibleHash(bucketCapacity int, hashFunc HashFunc) *ExtendibleHash {
	if bucketCapacity <= 0 {
		panic("bucketCapacity must be > 0")
	}

	if hashFunc == nil {
		hashFunc = defaultHash
	}

	// 初始目录大小为2^globalDepth
	initialSize := 1 // 初始globalDepth=0，大小为1
	buckets := make([]*HashBucket, initialSize)
	directory := make([]*HashBucket, initialSize)

	// 创建初始桶
	bucket := NewHashBucket()
	buckets[0] = bucket
	directory[0] = bucket

	return &ExtendibleHash{
		buckets:      buckets,
		directory:    directory,
		globalDepth:  0,
		bucketCapacity: bucketCapacity,
		hashFunc:      hashFunc,
		count:        0,
	}
}

// getBucketIndex 获取键对应的桶索引
func (eh *ExtendibleHash) getBucketIndex(key interface{}) (uint32, uint32) {
	if key == nil {
		return 0, 0
	}

	// 将键转换为字节数组
	keyBytes := []byte(fmt.Sprintf("%v", key))

	// 计算哈希值
	hashValue := eh.hashFunc(keyBytes)

	// 使用低globalDepth位作为索引
	index := hashValue & uint32((1<<eh.globalDepth)-1)
	fullHash := hashValue

	return index, fullHash
}

// Insert 插入键值对
func (eh *ExtendibleHash) Insert(key interface{}, value interface{}) error {
	eh.mu.Lock()
	defer eh.mu.Unlock()

	if key == nil {
		return fmt.Errorf("key cannot be nil")
	}

	index, _ := eh.getBucketIndex(key)
	bucket := eh.directory[index]

	// 检查桶中是否已存在该键
	for i, k := range bucket.keys {
		if fmt.Sprintf("%v", k) == fmt.Sprintf("%v", key) {
			bucket.values[i] = value
			return nil
		}
	}

	// 如果桶未满，直接插入
	if !bucket.isFull(eh.bucketCapacity) {
		bucket.keys = append(bucket.keys, key)
		bucket.values = append(bucket.values, value)
		eh.count++
		return nil
	}

	// 桶已满，需要分裂
	eh.splitBucket(bucket, index)
	eh.count++

	return nil
}

// splitBucket 分裂桶
func (eh *ExtendibleHash) splitBucket(bucket *HashBucket, index uint32) {
	// 创建两个新桶
	newBucket1 := NewHashBucket()
	newBucket2 := NewHashBucket()

	// 设置新桶的局部深度
	bucket.localDepth++
	newBucket1.localDepth = bucket.localDepth
	newBucket2.localDepth = bucket.localDepth

	// 重新分配键值对
	for i, key := range bucket.keys {
		keyBytes := []byte(fmt.Sprintf("%v", key))
		hashValue := eh.hashFunc(keyBytes)

		// 使用新的局部深度确定桶位置
		bit := (hashValue >> (bucket.localDepth - 1)) & 1
		if bit == 0 {
			newBucket1.keys = append(newBucket1.keys, key)
			newBucket1.values = append(newBucket1.values, bucket.values[i])
		} else {
			newBucket2.keys = append(newBucket2.keys, key)
			newBucket2.values = append(newBucket2.values, bucket.values[i])
		}
	}

	// 如果局部深度超过全局深度，需要扩展目录
	if bucket.localDepth > eh.globalDepth {
		eh.expandDirectory(bucket.localDepth)
	}

	// 更新目录指针
	bit := (index >> (bucket.localDepth - 1)) & 1
	if bit == 0 {
		eh.updateDirectoryPointers(index, newBucket1, newBucket2)
	} else {
		eh.updateDirectoryPointers(index, newBucket2, newBucket1)
	}
}

// expandDirectory 扩展目录
func (eh *ExtendibleHash) expandDirectory(newDepth int) {
	oldSize := len(eh.directory)
	newSize := oldSize * 2

	// 创建新的目录
	newDirectory := make([]*HashBucket, newSize)
	copy(newDirectory, eh.directory)

	// 添加新桶到buckets数组
	for i := 0; i < oldSize; i++ {
		bucketsUsed := make(map[*HashBucket]bool)
		for _, bucket := range newDirectory[i:] {
			if bucket != nil {
				bucketsUsed[bucket] = true
			}
		}
	}

	eh.directory = newDirectory
	eh.globalDepth = newDepth

	// 重新分配目录指针
	for i := 0; i < oldSize; i++ {
		if eh.directory[i] != nil {
			eh.directory[i+oldSize] = eh.directory[i]
		}
	}
}

// updateDirectoryPointers 更新目录指针
func (eh *ExtendibleHash) updateDirectoryPointers(index uint32, bucket1, bucket2 *HashBucket) {
	mask := uint32((1 << bucket1.localDepth) - 1)
	prefix := index & ^mask

	// 更新所有匹配前缀的目录项
	for i := uint32(0); i < uint32(len(eh.directory)); i++ {
		if (i & mask) == prefix {
			bit := (i >> (bucket1.localDepth - 1)) & 1
			if bit == 0 {
				eh.directory[i] = bucket1
			} else {
				eh.directory[i] = bucket2
			}
		}
	}
}

// Search 查找值
func (eh *ExtendibleHash) Search(key interface{}) (interface{}, bool) {
	eh.mu.RLock()
	defer eh.mu.RUnlock()

	if key == nil {
		return nil, false
	}

	index, _ := eh.getBucketIndex(key)
	bucket := eh.directory[index]

	// 在桶中查找键
	for i, k := range bucket.keys {
		if fmt.Sprintf("%v", k) == fmt.Sprintf("%v", key) {
			return bucket.values[i], true
		}
	}

	return nil, false
}

// Delete 删除键值对
func (eh *ExtendibleHash) Delete(key interface{}) bool {
	eh.mu.Lock()
	defer eh.mu.Unlock()

	if key == nil {
		return false
	}

	index, _ := eh.getBucketIndex(key)
	bucket := eh.directory[index]

	// 查找并删除键
	for i, k := range bucket.keys {
		if fmt.Sprintf("%v", k) == fmt.Sprintf("%v", key) {
			bucket.keys = append(bucket.keys[:i], bucket.keys[i+1:]...)
			bucket.values = append(bucket.values[:i], bucket.values[i+1:]...)
			eh.count--
			return true
		}
	}

	return false
}

// GetBucketInfo 获取桶信息（用于调试和监控）
func (eh *ExtendibleHash) GetBucketInfo() map[int]int {
	eh.mu.RLock()
	defer eh.mu.RUnlock()

	bucketCounts := make(map[int]int)
	for _, bucket := range eh.directory {
		if bucket != nil {
			bucketCounts[bucket.localDepth]++
		}
	}

	return bucketCounts
}

// GetBucketUsage 获取桶使用率统计
func (eh *ExtendibleHash) GetBucketUsage() (avg float64, max int, min int, fullCount int) {
	eh.mu.RLock()
	defer eh.mu.RUnlock()

	if len(eh.directory) == 0 {
		return 0, 0, 0, 0
	}

	total := 0
	max = 0
	min = eh.bucketCapacity
	fullCount = 0

	for _, bucket := range eh.directory {
		if bucket != nil {
			size := len(bucket.keys)
			total += size
			if size > max {
				max = size
			}
			if size < min {
				min = size
			}
			if size >= eh.bucketCapacity {
				fullCount++
			}
		}
	}

	avg = float64(total) / float64(len(eh.directory))
	return
}

// Size 返回键值对数量
func (eh *ExtendibleHash) Size() int64 {
	eh.mu.RLock()
	defer eh.mu.RUnlock()
	return eh.count
}

// GlobalDepth 返回全局深度
func (eh *ExtendibleHash) GlobalDepth() int {
	eh.mu.RLock()
	defer eh.mu.RUnlock()
	return eh.globalDepth
}

// BucketCount 返回桶数量
func (eh *ExtendibleHash) BucketCount() int {
	eh.mu.RLock()
	defer eh.mu.RUnlock()
	return len(eh.directory)
}

// String 返回哈希表的字符串表示（用于调试）
func (eh *ExtendibleHash) String() string {
	eh.mu.RLock()
	defer eh.mu.RUnlock()

	result := fmt.Sprintf("ExtendibleHash(globalDepth=%d, bucketCount=%d, count=%d):\n",
		eh.globalDepth, len(eh.directory), eh.count)

	bucketInfo := make(map[*HashBucket]int)
	for _, bucket := range eh.directory {
		if bucket != nil {
			if _, ok := bucketInfo[bucket]; !ok {
				bucketInfo[bucket] = len(bucket.keys)
			}
		}
	}

	for bucket, size := range bucketInfo {
		result += fmt.Sprintf("  Bucket(localDepth=%d, size=%d)\n", bucket.localDepth, size)
	}

	return result
}

// NewExtendibleHashWithDefault 创建默认配置的可扩展哈希表
func NewExtendibleHashWithDefault() *ExtendibleHash {
	return NewExtendibleHash(4, defaultHash)
}

// Compact 压缩目录（可选操作，减少不使用的目录项）
func (eh *ExtendibleHash) Compact() {
	eh.mu.Lock()
	defer eh.mu.Unlock()

	// 统计目录使用情况
	bucketsUsed := make(map[*HashBucket]bool)
	for _, bucket := range eh.directory {
		if bucket != nil {
			bucketsUsed[bucket] = true
		}
	}

	// 如果所有桶的局部深度都小于全局深度，可以减少全局深度
	if eh.globalDepth > 0 {
		canReduce := true
		for bucket := range bucketsUsed {
			if bucket.localDepth == eh.globalDepth {
				canReduce = false
				break
			}
		}

		if canReduce {
			eh.globalDepth--
			newSize := 1 << eh.globalDepth
			eh.directory = eh.directory[:newSize]
		}
	}
}
