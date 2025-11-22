package datastructures

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash"
	"hash/fnv"
	"math"
	"sync"
)

// BloomFilter 布隆过滤器
// 特点：
// - 空间效率高的概率性数据结构
// - 有假阳性（false positive），无假阴性（false negative）
// - 用于快速判断元素是否在集合中
// - 常用于数据库查询优化、缓存穿透防护
type BloomFilter struct {
	bitArray  []byte     // 位数组
	m         uint       // 位数组大小（位数）
	k         uint       // 哈希函数数量
	count     uint64     // 已插入元素数量
	hashFuncs []hash.Hash32 // 哈希函数列表
	mu        sync.RWMutex // 读写锁
}

// NewBloomFilter 创建新的布隆过滤器
// expectedElements: 期望插入的元素数量
// falsePositiveRate: 期望的假阳性率 (0 < fpr < 1)
func NewBloomFilter(expectedElements uint, falsePositiveRate float64) *BloomFilter {
	if expectedElements == 0 {
		panic("expectedElements must be > 0")
	}
	if falsePositiveRate <= 0 || falsePositiveRate >= 1 {
		panic("falsePositiveRate must be in (0, 1)")
	}

	// 计算最优的位数组大小
	m := uint(-float64(expectedElements) * math.Log(falsePositiveRate) / (math.Log(2) * math.Log(2)))

	// 计算最优的哈希函数数量
	k := uint(float64(m) / float64(expectedElements) * math.Log(2))

	// 确保k至少为1
	if k == 0 {
		k = 1
	}

	// 初始化哈希函数
	hashFuncs := make([]hash.Hash32, k)
	for i := uint(0); i < k; i++ {
		// 为每个哈希函数使用不同的盐值
		h := fnv.New32a()
		salt := make([]byte, 4)
		binary.BigEndian.PutUint32(salt, uint32(i))
		h.Write(salt)
		hashFuncs[i] = h
	}

	return &BloomFilter{
		bitArray:  make([]byte, (m+7)/8),
		m:         m,
		k:         k,
		hashFuncs: hashFuncs,
		count:     0,
	}
}

// getHashPositions 获取元素对应的哈希位置
func (bf *BloomFilter) getHashPositions(data []byte) []uint {
	positions := make([]uint, bf.k)

	for i, hashFunc := range bf.hashFuncs {
		hashFunc.Reset()
		hashFunc.Write(data)
		hashValue := hashFunc.Sum32()

		// 使用哈希值计算位置
		position := hashValue % uint32(bf.m)
		positions[i] = uint(position)
	}

	return positions
}

// Add 添加元素
func (bf *BloomFilter) Add(data []byte) {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	positions := bf.getHashPositions(data)

	for _, pos := range positions {
		// 计算字节和位的位置
	_byte := pos / 8
	_bit := pos % 8

		// 设置位
		bf.bitArray[_byte] |= 1 << _bit
	}

	bf.count++
}

// AddString 添加字符串元素
func (bf *BloomFilter) AddString(s string) {
	bf.Add([]byte(s))
}

// AddInt 添加整数元素
func (bf *BloomFilter) AddInt(n int) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(n))
	bf.Add(b)
}

// AddFloat64 添加浮点数元素
func (bf *BloomFilter) AddFloat64(f float64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(f))
	bf.Add(b)
}

// AddGeneric 添加泛型元素
func (bf *BloomFilter) AddGeneric(item interface{}) {
	data := []byte(fmt.Sprintf("%v", item))
	bf.Add(data)
}

// Contains 检查元素是否存在
// 返回true表示可能存在，返回false表示一定不存在
func (bf *BloomFilter) Contains(data []byte) bool {
	bf.mu.RLock()
	defer bf.mu.RUnlock()

	positions := bf.getHashPositions(data)

	for _, pos := range positions {
		_byte := pos / 8
		_bit := pos % 8

		// 检查位是否被设置
		if (bf.bitArray[_byte] & (1 << _bit)) == 0 {
			return false
		}
	}

	return true
}

// ContainsString 检查字符串元素是否存在
func (bf *BloomFilter) ContainsString(s string) bool {
	return bf.Contains([]byte(s))
}

// ContainsInt 检查整数元素是否存在
func (bf *BloomFilter) ContainsInt(n int) bool {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(n))
	return bf.Contains(b)
}

// ContainsFloat64 检查浮点数元素是否存在
func (bf *BloomFilter) ContainsFloat64(f float64) bool {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(f))
	return bf.Contains(b)
}

// ContainsGeneric 检查泛型元素是否存在
func (bf *BloomFilter) ContainsGeneric(item interface{}) bool {
	data := []byte(fmt.Sprintf("%v", item))
	return bf.Contains(data)
}

// GetFalsePositiveRate 获取当前假阳性率
func (bf *BloomFilter) GetFalsePositiveRate() float64 {
	bf.mu.RLock()
	defer bf.mu.RUnlock()

	// 使用标准公式计算假阳性率
	// FPR = (1 - e^(-kn/m))^k
	// 其中 k 是哈希函数数量，n 是元素数量，m 是位数组大小
	expValue := math.Exp(-float64(bf.k) * float64(bf.count) / float64(bf.m))
	fpr := math.Pow(1-expValue, float64(bf.k))

	return fpr
}

// Clear 清空布隆过滤器
func (bf *BloomFilter) Clear() {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	for i := range bf.bitArray {
		bf.bitArray[i] = 0
	}
	bf.count = 0
}

// Size 返回已插入元素数量
func (bf *BloomFilter) Size() uint64 {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	return bf.count
}

// BitArraySize 返回位数组大小（字节数）
func (bf *BloomFilter) BitArraySize() uint {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	return uint(len(bf.bitArray))
}

// BitSize 返回位数组大小（位数）
func (bf *BloomFilter) BitSize() uint {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	return bf.m
}

// HashFuncCount 返回哈希函数数量
func (bf *BloomFilter) HashFuncCount() uint {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	return bf.k
}

// Merge 合并另一个布隆过滤器
func (bf *BloomFilter) Merge(other *BloomFilter) error {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	if other.m != bf.m || other.k != bf.k {
		return fmt.Errorf("cannot merge bloom filters with different parameters")
	}

	for i, b := range other.bitArray {
		bf.bitArray[i] |= b
	}

	bf.count += other.count
	return nil
}

// Clone 克隆布隆过滤器
func (bf *BloomFilter) Clone() *BloomFilter {
	bf.mu.RLock()
	defer bf.mu.RUnlock()

	newBf := &BloomFilter{
		bitArray:  make([]byte, len(bf.bitArray)),
		m:         bf.m,
		k:         bf.k,
		count:     bf.count,
		hashFuncs: make([]hash.Hash32, len(bf.hashFuncs)),
	}

	copy(newBf.bitArray, bf.bitArray)
	copy(newBf.hashFuncs, bf.hashFuncs)

	return newBf
}

// Serialize 序列化布隆过滤器
func (bf *BloomFilter) Serialize() ([]byte, error) {
	bf.mu.RLock()
	defer bf.mu.RUnlock()

	data := struct {
		BitArray []byte
		M        uint
		K        uint
		Count    uint64
	}{
		BitArray: bf.bitArray,
		M:        bf.m,
		K:        bf.k,
		Count:    bf.count,
	}

	return json.Marshal(data)
}

// Deserialize 反序列化布隆过滤器
func Deserialize(data []byte) (*BloomFilter, error) {
	var bfData struct {
		BitArray []byte
		M        uint
		K        uint
		Count    uint64
	}

	if err := json.Unmarshal(data, &bfData); err != nil {
		return nil, err
	}

	// 重新初始化哈希函数
	hashFuncs := make([]hash.Hash32, bfData.K)
	for i := uint(0); i < bfData.K; i++ {
		h := fnv.New32a()
		salt := make([]byte, 4)
		binary.BigEndian.PutUint32(salt, uint32(i))
		h.Write(salt)
		hashFuncs[i] = h
	}

	return &BloomFilter{
		bitArray:  bfData.BitArray,
		m:         bfData.M,
		k:         bfData.K,
		count:     bfData.Count,
		hashFuncs: hashFuncs,
	}, nil
}

// String 返回布隆过滤器的字符串表示
func (bf *BloomFilter) String() string {
	bf.mu.RLock()
	defer bf.mu.RUnlock()

	return fmt.Sprintf("BloomFilter(m=%d, k=%d, count=%d, fpr=%.6f)",
		bf.m, bf.k, bf.count, bf.GetFalsePositiveRate())
}

// NewOptimalBloomFilter 根据元素数量和假阳性率创建最优布隆过滤器
func NewOptimalBloomFilter(expectedElements uint, falsePositiveRate float64) *BloomFilter {
	return NewBloomFilter(expectedElements, falsePositiveRate)
}

// NewDefaultBloomFilter 创建默认配置的布隆过滤器
// 默认：100000个元素，1%假阳性率
func NewDefaultBloomFilter() *BloomFilter {
	return NewBloomFilter(100000, 0.01)
}

// NewSmallBloomFilter 创建小型布隆过滤器
// 适用于内存受限的场景
func NewSmallBloomFilter(expectedElements uint) *BloomFilter {
	return NewBloomFilter(expectedElements, 0.05) // 5%假阳性率
}

// NewLargeBloomFilter 创建大型布隆过滤器
// 适用于对精度要求较高的场景
func NewLargeBloomFilter(expectedElements uint) *BloomFilter {
	return NewBloomFilter(expectedElements, 0.001) // 0.1%假阳性率
}
