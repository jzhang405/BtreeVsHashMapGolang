package datastructures

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// SkipNode 跳表节点
type SkipNode struct {
	key     any    // 键
	value   any    // 值
	forward []*SkipNode    // 前向指针数组，每一层的下一个节点
	height  int            // 节点高度（层数）
}

// NewSkipNode 创建新的跳表节点
func NewSkipNode(key, value any, height int) *SkipNode {
	return &SkipNode{
		key:     key,
		value:   value,
		forward: make([]*SkipNode, height),
		height:  height,
	}
}

// SkipList 跳表结构
// 特点：
// - 融合了数组和链表的优点
// - 有序结构，支持范围查询
// - O(log n)时间复杂度的操作
// - 实现简单，适合内存场景
type SkipList struct {
	head     *SkipNode    // 头节点
	comparator Comparator // 比较函数
	mu       sync.RWMutex // 读写锁
	maxLevel int          // 最大层数
	level    int          // 当前最大层数
	prob     float64      // 随机层数的概率因子 (0 < prob < 1)
	count    int64        // 元素总数
}

// NewSkipList 创建新的跳表
// maxLevel: 最大层数，建议16-32
// prob: 升层概率，建议0.5
// comparator: 比较函数
func NewSkipList(maxLevel int, prob float64, comparator Comparator) *SkipList {
	if maxLevel < 1 {
		panic("maxLevel must be >= 1")
	}
	if prob <= 0 || prob >= 1 {
		panic("prob must be in (0, 1)")
	}
	if comparator == nil {
		panic("comparator is required")
	}

	return &SkipList{
		head:       NewSkipNode(nil, nil, maxLevel),
		comparator: comparator,
		maxLevel:   maxLevel,
		level:      1,
		prob:       prob,
		count:      0,
	}
}

// randomLevel 生成随机层数
// 使用几何分布，概率为p的节点有第k层
func (s *SkipList) randomLevel() int {
	level := 1
	for level < s.maxLevel && rand.Float64() < s.prob {
		level++
	}
	return level
}

// Insert 插入键值对
func (s *SkipList) Insert(key any, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if key == nil {
		return fmt.Errorf("key cannot be nil")
	}

	// 查找插入位置和更新指针
	update := make([]*SkipNode, s.maxLevel)
	x := s.head

	// 从最高层开始查找，找到每层的插入位置
	for i := s.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && s.comparator(x.forward[i].key, key) < 0 {
			x = x.forward[i]
		}
		update[i] = x
	}

	// 如果键已存在，更新值
	if x.forward[0] != nil && s.comparator(x.forward[0].key, key) == 0 {
		x.forward[0].value = value
		return nil
	}

	// 生成随机层数
	newLevel := s.randomLevel()
	if newLevel > s.level {
		// 如果新层数超过当前最大层数，补充update数组
		for i := s.level; i < newLevel; i++ {
			update[i] = s.head
		}
		s.level = newLevel
	}

	// 创建新节点
	newNode := NewSkipNode(key, value, newLevel)

	// 更新指针
	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	s.count++
	return nil
}

// Search 查找值
func (s *SkipList) Search(key any) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if key == nil {
		return nil, false
	}

	x := s.head

	// 从最高层开始查找
	for i := s.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && s.comparator(x.forward[i].key, key) < 0 {
			x = x.forward[i]
		}
	}

	// 前进到第一层的下一个节点
	x = x.forward[0]

	// 检查是否找到
	if x != nil && s.comparator(x.key, key) == 0 {
		return x.value, true
	}

	return nil, false
}

// Delete 删除键值对
func (s *SkipList) Delete(key any) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if key == nil {
		return false
	}

	update := make([]*SkipNode, s.maxLevel)
	x := s.head

	// 查找要删除的节点和更新指针
	for i := s.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && s.comparator(x.forward[i].key, key) < 0 {
			x = x.forward[i]
		}
		update[i] = x
	}

	x = x.forward[0]

	// 如果找到要删除的节点
	if x != nil && s.comparator(x.key, key) == 0 {
		// 更新指针
		for i := 0; i < s.level; i++ {
			if update[i].forward[i] != x {
				break
			}
			update[i].forward[i] = x.forward[i]
		}

		// 移除最高层为空的头指针
		for s.level > 1 && s.head.forward[s.level-1] == nil {
			s.level--
		}

		s.count--
		return true
	}

	return false
}

// RangeQuery 范围查询 [start, end)
func (s *SkipList) RangeQuery(start, end any) ([]KeyValue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if start == nil || end == nil {
		return nil, fmt.Errorf("start and end cannot be nil")
	}

	if s.comparator(start, end) >= 0 {
		return nil, fmt.Errorf("start must be less than end")
	}

	var result []KeyValue
	x := s.head

	// 找到起始节点
	for i := s.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && s.comparator(x.forward[i].key, start) < 0 {
			x = x.forward[i]
		}
	}

	x = x.forward[0]

	// 遍历直到达到结束条件
	for x != nil && s.comparator(x.key, end) < 0 {
		result = append(result, KeyValue{Key: x.key, Value: x.value})
		x = x.forward[0]
	}

	return result, nil
}

// ScanAll 顺序遍历所有键值对
func (s *SkipList) ScanAll() []KeyValue {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []KeyValue
	x := s.head.forward[0]

	for x != nil {
		result = append(result, KeyValue{Key: x.key, Value: x.value})
		x = x.forward[0]
	}

	return result
}

// Size 返回元素数量
func (s *SkipList) Size() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.count
}

// Level 返回当前最大层数
func (s *SkipList) Level() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.level
}

// MaxLevel 返回最大层数
func (s *SkipList) MaxLevel() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.maxLevel
}

// Height 返回跳表的高度（统计值）
func (s *SkipList) Height() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	height := 0
	x := s.head
	for x.forward[0] != nil {
		height++
		x = x.forward[0]
	}
	return height
}

// String 返回跳表的字符串表示（用于调试）
func (s *SkipList) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result string
	result += fmt.Sprintf("SkipList(level=%d, count=%d):\n", s.level, s.count)

	// 显示每一层的节点
	for i := s.level - 1; i >= 0; i-- {
		result += fmt.Sprintf("Level %d: ", i)
		x := s.head
		for x.forward[i] != nil {
			result += fmt.Sprintf("%v -> ", x.forward[i].key)
			x = x.forward[i]
		}
		result += "nil\n"
	}

	return result
}

// NewDefaultSkipList 创建默认配置的跳表
// 使用合理的默认参数：
// - maxLevel: 16
// - prob: 0.5
// 适用于大多数场景
func NewDefaultSkipList(comparator Comparator) *SkipList {
	// 使用随机种子
	rand.Seed(time.Now().UnixNano())
	return NewSkipList(16, 0.5, comparator)
}
