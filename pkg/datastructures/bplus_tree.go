package datastructures

import (
	"fmt"
	"sync"
)

// BPlusTree 生产级别的B+树实现
// 特点：
// - 所有数据存储在叶子节点
// - 叶子节点通过指针连接，支持高效范围查询
// - 磁盘友好，适合大数据场景
// - 支持泛型键值对

type (
	// Comparator 比较函数，返回值：-1(小于)、0(等于)、1(大于)
	Comparator func(a, b any) int

	// KeyValue 键值对
	KeyValue struct {
		Key   any
		Value any
	}

	// TreeNode B+树节点
	TreeNode struct {
		keys     []any       // 排序后的键列表
		values   []KeyValue  // 叶子节点的值列表
		children []*TreeNode // 子节点列表（内部节点）
		isLeaf   bool        // 是否为叶子节点
		next     *TreeNode   // 叶子节点链表指针（仅叶子节点使用）
		parent   *TreeNode   // 父节点指针
	}
)

// BPlusTree B+树主体结构
type BPlusTree struct {
	root       *TreeNode // 根节点
	order      int       // 阶数（每个节点最大子节点数）
	minKeys    int       // 最小键数
	minChildren int      // 最小子节点数
	comparator Comparator // 比较函数
	mu         sync.RWMutex // 读写锁，支持并发访问
	count      int64      // 总键数
}

// NewBPlusTree 创建新的B+树
// order: 树的阶数，建议值：64-256（根据磁盘块大小调整）
// comparator: 比较函数，用于键的排序
func NewBPlusTree(order int, comparator Comparator) *BPlusTree {
	if order < 3 {
		panic("order must be >= 3")
	}
	if comparator == nil {
		panic("comparator is required")
	}

	return &BPlusTree{
		root:       &TreeNode{isLeaf: true},
		order:      order,
		minKeys:    order/2 - 1,
		minChildren: order / 2,
		comparator: comparator,
		count:      0,
	}
}

// Insert 插入键值对
func (t *BPlusTree) Insert(key any, value any) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if key == nil {
		return fmt.Errorf("key cannot be nil")
	}

	// 查找叶子节点
	leaf := t.findLeafNode(key)

	// 检查是否已存在该键
	for i, k := range leaf.keys {
		cmp := t.comparator(k, key)
		if cmp == 0 {
			// 更新已存在的键
			leaf.values[i].Value = value
			return nil
		}
	}

	// 插入新键值对
	t.insertIntoLeaf(leaf, key, value)
	t.count++

	return nil
}

// Search 查找值
func (t *BPlusTree) Search(key any) (any, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if key == nil {
		return nil, false
	}

	leaf := t.findLeafNode(key)

	// 在叶子节点中查找键
	for i, k := range leaf.keys {
		cmp := t.comparator(k, key)
		if cmp == 0 {
			return leaf.values[i].Value, true
		}
	}

	return nil, false
}

// Delete 删除键值对
func (t *BPlusTree) Delete(key any) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if key == nil {
		return false
	}

	leaf := t.findLeafNode(key)

	// 查找键的位置
	idx := -1
	for i, k := range leaf.keys {
		if t.comparator(k, key) == 0 {
			idx = i
			break
		}
	}

	if idx == -1 {
		return false // 键不存在
	}

	// 从叶子节点中删除
	t.deleteFromLeaf(leaf, idx)
	t.count--

	return true
}

// RangeQuery 范围查询 [start, end)
func (t *BPlusTree) RangeQuery(start, end any) ([]KeyValue, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if start == nil || end == nil {
		return nil, fmt.Errorf("start and end cannot be nil")
	}

	if t.comparator(start, end) >= 0 {
		return nil, fmt.Errorf("start must be less than end")
	}

	var result []KeyValue
	leaf := t.findLeafNode(start)

	// 遍历叶子节点链表
	for leaf != nil {
		for i, key := range leaf.keys {
			if t.comparator(key, start) >= 0 && t.comparator(key, end) < 0 {
				result = append(result, leaf.values[i])
			}
		}
		leaf = leaf.next
	}

	return result, nil
}

// ScanAll 顺序遍历所有键值对
func (t *BPlusTree) ScanAll() []KeyValue {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []KeyValue

	// 找到最左叶子节点
	leaf := t.root
	for !leaf.isLeaf {
		leaf = leaf.children[0]
	}

	// 遍历所有叶子节点
	for leaf != nil {
		result = append(result, leaf.values...)
		leaf = leaf.next
	}

	return result
}

// Size 返回树中键值对数量
func (t *BPlusTree) Size() int64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.count
}

// Height 返回树的高度
func (t *BPlusTree) Height() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	height := 1
	node := t.root
	for !node.isLeaf {
		height++
		node = node.children[0]
	}
	return height
}

// 内部方法：查找包含指定键的叶子节点
func (t *BPlusTree) findLeafNode(key any) *TreeNode {
	node := t.root

	// 从根节点向下查找，直到叶子节点
	for !node.isLeaf {
		// 找到第一个大于key的键的位置
		idx := 0
		for idx < len(node.keys) && t.comparator(key, node.keys[idx]) >= 0 {
			idx++
		}
		// 确保索引不越界
		if idx >= len(node.children) {
			idx = len(node.children) - 1
		}
		node = node.children[idx]
	}

	return node
}

// 内部方法：将键值对插入叶子节点
func (t *BPlusTree) insertIntoLeaf(leaf *TreeNode, key any, value any) {
	// 找到插入位置
	insertPos := 0
	for insertPos < len(leaf.keys) && t.comparator(leaf.keys[insertPos], key) < 0 {
		insertPos++
	}

	// 插入键
	leaf.keys = append(leaf.keys, nil)
	copy(leaf.keys[insertPos+1:], leaf.keys[insertPos:])
	leaf.keys[insertPos] = key

	// 插入值
	leaf.values = append(leaf.values, KeyValue{})
	copy(leaf.values[insertPos+1:], leaf.values[insertPos:])
	leaf.values[insertPos] = KeyValue{Key: key, Value: value}

	// 检查是否需要分裂
	if len(leaf.keys) > t.order-1 {
		t.splitLeafNode(leaf)
	}
}

// 内部方法：分裂叶子节点
func (t *BPlusTree) splitLeafNode(leaf *TreeNode) {
	// 找到分裂点
	splitPos := len(leaf.keys) / 2

	// 创建新叶子节点
	newLeaf := &TreeNode{
		isLeaf: true,
		next:   leaf.next,
		parent: leaf.parent,
	}

	// 移动键值对到新节点
	newLeaf.keys = append(newLeaf.keys, leaf.keys[splitPos:]...)
	newLeaf.values = append(newLeaf.values, leaf.values[splitPos:]...)

	// 更新原节点
	leaf.keys = leaf.keys[:splitPos]
	leaf.values = leaf.values[:splitPos]
	leaf.next = newLeaf

	// 如果这是根节点，创建新根节点
	if leaf.parent == nil {
		newRoot := &TreeNode{
			keys:     []any{newLeaf.keys[0]},
			children: []*TreeNode{leaf, newLeaf},
			isLeaf:   false,
		}
		leaf.parent = newRoot
		newLeaf.parent = newRoot
		t.root = newRoot
		return
	}

	// 更新父节点 - 找到原叶节点在父节点中的位置
	parent := leaf.parent
	leafPos := -1
	for i := 0; i < len(parent.children); i++ {
		if parent.children[i] == leaf {
			leafPos = i
			break
		}
	}

	if leafPos == -1 {
		// 这不应该发生
		panic("leaf node not found in parent")
	}

	// 在 leafPos 位置插入新键，在 leafPos+1 位置插入新子节点
	parent.keys = append(parent.keys, nil)
	copy(parent.keys[leafPos+1:], parent.keys[leafPos:])
	parent.keys[leafPos] = newLeaf.keys[0]

	parent.children = append(parent.children, nil)
	copy(parent.children[leafPos+2:], parent.children[leafPos+1:])
	parent.children[leafPos+1] = newLeaf
	newLeaf.parent = parent

	// 检查父节点是否需要分裂
	if len(parent.keys) > t.order-1 {
		t.splitInternalNode(parent)
	}
}

// 内部方法：分裂内部节点
func (t *BPlusTree) splitInternalNode(node *TreeNode) {
	// 找到分裂点（注意：内部节点的键不会移动到新节点）
	splitPos := len(node.keys) / 2

	// 创建新内部节点
	newNode := &TreeNode{
		isLeaf:   false,
		parent:   node.parent,
	}

	// 移动键到新节点（不包括中间的键）
	newNode.keys = append(newNode.keys, node.keys[splitPos+1:]...)

	// 移动子节点到新节点
	midKey := node.keys[splitPos]
	newNode.children = append(newNode.children, node.children[splitPos+1:]...)

	// 更新子节点的父指针
	for _, child := range newNode.children {
		child.parent = newNode
	}

	// 更新原节点
	node.keys = node.keys[:splitPos]
	node.children = node.children[:splitPos+1]

	// 如果这是根节点，创建新根节点
	if node.parent == nil {
		newRoot := &TreeNode{
			keys:     []any{midKey},
			children: []*TreeNode{node, newNode},
			isLeaf:   false,
		}
		node.parent = newRoot
		newNode.parent = newRoot
		t.root = newRoot
		return
	}

	// 更新父节点 - 找到原节点在父节点中的位置
	parent := node.parent
	nodePos := -1
	for i := 0; i < len(parent.children); i++ {
		if parent.children[i] == node {
			nodePos = i
			break
		}
	}

	if nodePos == -1 {
		// 这不应该发生
		panic("node not found in parent")
	}

	// 在 nodePos 位置插入新键，在 nodePos+1 位置插入新子节点
	parent.keys = append(parent.keys, nil)
	copy(parent.keys[nodePos+1:], parent.keys[nodePos:])
	parent.keys[nodePos] = midKey

	parent.children = append(parent.children, nil)
	copy(parent.children[nodePos+2:], parent.children[nodePos+1:])
	parent.children[nodePos+1] = newNode
	newNode.parent = parent

	// 检查祖父节点是否需要分裂
	if len(parent.keys) > t.order-1 {
		t.splitInternalNode(parent)
	}
}

// 内部方法：从叶子节点删除
func (t *BPlusTree) deleteFromLeaf(leaf *TreeNode, idx int) {
	// 从叶子节点中删除键值对
	leaf.keys = append(leaf.keys[:idx], leaf.keys[idx+1:]...)
	leaf.values = append(leaf.values[:idx], leaf.values[idx+1:]...)

	// 检查是否需要合并或借键
	if len(leaf.keys) < t.minKeys && leaf.parent != nil {
		t.rebalanceLeafNode(leaf)
	}
}

// 内部方法：重新平衡叶子节点
func (t *BPlusTree) rebalanceLeafNode(leaf *TreeNode) {
	parent := leaf.parent

	// 找到在父节点中的位置
	pos := 0
	for pos < len(parent.children) && parent.children[pos] != leaf {
		pos++
	}

	// 尝试从左兄弟节点借键
	if pos > 0 && len(parent.children[pos-1].keys) > t.minKeys {
		leftSibling := parent.children[pos-1]

		// 从左兄弟借最后一个键
		borrowedKey := leftSibling.keys[len(leftSibling.keys)-1]
		borrowedValue := leftSibling.values[len(leftSibling.values)-1]

		// 在当前节点前插入借来的键
		leaf.keys = append([]any{borrowedKey}, leaf.keys...)
		leaf.values = append([]KeyValue{borrowedValue}, leaf.values...)

		// 从左兄弟删除借出的键
		leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
		leftSibling.values = leftSibling.values[:len(leftSibling.values)-1]

		// 更新父节点中的键
		parent.keys[pos-1] = leaf.keys[0]
		return
	}

	// 尝试从右兄弟节点借键
	if pos < len(parent.children)-1 && len(parent.children[pos+1].keys) > t.minKeys {
		rightSibling := parent.children[pos+1]

		// 从右兄弟借第一个键
		borrowedKey := rightSibling.keys[0]
		borrowedValue := rightSibling.values[0]

		// 在当前节点后插入借来的键
		leaf.keys = append(leaf.keys, borrowedKey)
		leaf.values = append(leaf.values, borrowedValue)

		// 从右兄弟删除借出的键
		rightSibling.keys = rightSibling.keys[1:]
		rightSibling.values = rightSibling.values[1:]

		// 更新父节点中的键
		parent.keys[pos] = rightSibling.keys[0]
		return
	}

	// 合并节点
	if pos > 0 {
		// 与左兄弟合并
		leftSibling := parent.children[pos-1]
		leftSibling.keys = append(leftSibling.keys, leaf.keys...)
		leftSibling.values = append(leftSibling.values, leaf.values...)
		leftSibling.next = leaf.next

		// 从父节点删除键和子节点
		t.deleteFromInternalNode(parent, pos-1, leaf)
	} else {
		// 与右兄弟合并
		rightSibling := parent.children[pos+1]
		leaf.keys = append(leaf.keys, rightSibling.keys...)
		leaf.values = append(leaf.values, rightSibling.values...)
		leaf.next = rightSibling.next

		// 从父节点删除键和子节点
		t.deleteFromInternalNode(parent, pos, rightSibling)
	}
}

// 内部方法：从内部节点删除
func (t *BPlusTree) deleteFromInternalNode(parent *TreeNode, pos int, child *TreeNode) {
	// 删除键和子节点
	parent.keys = append(parent.keys[:pos], parent.keys[pos+1:]...)
	parent.children = append(parent.children[:pos+1], parent.children[pos+2:]...)

	// 检查是否需要重新平衡
	if len(parent.keys) < t.minKeys && parent.parent != nil {
		t.rebalanceInternalNode(parent)
	}

	// 如果根节点只有一个子节点，简化树
	if len(parent.keys) == 0 && parent.parent == nil {
		t.root = parent.children[0]
		parent.children[0].parent = nil
	}
}

// 内部方法：重新平衡内部节点
func (t *BPlusTree) rebalanceInternalNode(node *TreeNode) {
	parent := node.parent

	// 找到在父节点中的位置
	pos := 0
	for pos < len(parent.children) && parent.children[pos] != node {
		pos++
	}

	// 尝试从左兄弟节点借键
	if pos > 0 && len(parent.children[pos-1].keys) > t.minKeys {
		leftSibling := parent.children[pos-1]

		// 从父节点借最后一个键到当前节点
		borrowedKey := parent.keys[pos-1]
		node.keys = append([]any{borrowedKey}, node.keys...)

		// 从左兄弟借最后一个子节点
		lastChild := leftSibling.children[len(leftSibling.children)-1]
		lastChild.parent = node
		node.children = append([]*TreeNode{lastChild}, node.children...)

		// 从左兄弟删除借出的键和子节点
		leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
		leftSibling.children = leftSibling.children[:len(leftSibling.children)-1]

		// 更新父节点中的键
		parent.keys[pos-1] = leftSibling.keys[len(leftSibling.keys)-1]
		return
	}

	// 尝试从右兄弟节点借键
	if pos < len(parent.children)-1 && len(parent.children[pos+1].keys) > t.minKeys {
		rightSibling := parent.children[pos+1]

		// 从父节点借第一个键到当前节点
		borrowedKey := parent.keys[pos]
		node.keys = append(node.keys, borrowedKey)

		// 从右兄弟借第一个子节点
		firstChild := rightSibling.children[0]
		firstChild.parent = node
		node.children = append(node.children, firstChild)

		// 从右兄弟删除借出的键和子节点
		rightSibling.keys = rightSibling.keys[1:]
		rightSibling.children = rightSibling.children[1:]

		// 更新父节点中的键
		parent.keys[pos] = rightSibling.keys[0]
		return
	}

	// 合并节点
	if pos > 0 {
		// 与左兄弟合并
		leftSibling := parent.children[pos-1]
		parentKey := parent.keys[pos-1]

		leftSibling.keys = append(leftSibling.keys, parentKey)
		leftSibling.keys = append(leftSibling.keys, node.keys...)
		leftSibling.children = append(leftSibling.children, node.children...)

		// 更新子节点的父指针
		for _, child := range node.children {
			child.parent = leftSibling
		}

		// 从父节点删除键和子节点
		t.deleteFromInternalNode(parent, pos-1, node)
	} else {
		// 与右兄弟合并
		rightSibling := parent.children[pos+1]
		parentKey := parent.keys[pos]

		node.keys = append(node.keys, parentKey)
		node.keys = append(node.keys, rightSibling.keys...)
		node.children = append(node.children, rightSibling.children...)

		// 更新子节点的父指针
		for _, child := range rightSibling.children {
			child.parent = node
		}

		// 从父节点删除键和子节点
		t.deleteFromInternalNode(parent, pos, rightSibling)
	}
}

// String 返回树的字符串表示（用于调试）
func (t *BPlusTree) String() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.root.String("")
}

// TreeNode的String方法
func (n *TreeNode) String(prefix string) string {
	result := prefix

	if n.isLeaf {
		result += "Leaf("
		for i, key := range n.keys {
			if i > 0 {
				result += ", "
			}
			result += fmt.Sprintf("%v", key)
		}
		result += ")\n"
	} else {
		result += "Internal("
		for i, key := range n.keys {
			if i > 0 {
				result += ", "
			}
			result += fmt.Sprintf("%v", key)
		}
		result += ")\n"

		for _, child := range n.children {
			result += child.String(prefix + "  ")
		}
	}

	return result
}
