package datastructures

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

// MerkleNode 默克尔树节点
type MerkleNode struct {
	hash     string      // 节点哈希值
	data     []byte      // 叶子节点的数据
	children []*MerkleNode // 子节点（最多2个）
	parent   *MerkleNode   // 父节点指针
	isLeaf   bool          // 是否为叶子节点
}

// NewMerkleNode 创建新的默克尔树节点
func NewMerkleNode(data []byte, left, right *MerkleNode) *MerkleNode {
	node := &MerkleNode{
		data:     data,
		children: make([]*MerkleNode, 0, 2),
		isLeaf:   left == nil && right == nil,
	}

	if left != nil {
		node.children = append(node.children, left)
		left.parent = node
	}

	if right != nil {
		node.children = append(node.children, right)
		right.parent = node
	}

	node.hash = node.computeHash()
	return node
}

// computeHash 计算节点哈希值
func (n *MerkleNode) computeHash() string {
	if n.isLeaf {
		// 叶子节点：直接对数据哈希
		hash := sha256.Sum256(n.data)
		return hex.EncodeToString(hash[:])
	}

	// 内部节点：对子节点哈希连接后哈希
	hashData := ""
	for _, child := range n.children {
		hashData += child.hash
	}

	hash := sha256.Sum256([]byte(hashData))
	return hex.EncodeToString(hash[:])
}

// MerkleTree 默克尔树
// 特点：
// - 叶子节点存储数据的哈希值
// - 非叶子节点存储子节点哈希的聚合值
// - 快速验证数据完整性和一致性
// - 支持范围查询（按叶子节点顺序）
// - 常用于区块链、分布式存储
type MerkleTree struct {
	root     *MerkleNode // 根节点
	leaves   []*MerkleNode // 所有叶子节点
	data     [][]byte     // 原始数据
	mu       sync.RWMutex // 读写锁
	count    int64        // 数据块数量
}

// NewMerkleTree 从数据块创建默克尔树
func NewMerkleTree(data [][]byte) *MerkleTree {
	if len(data) == 0 {
		return &MerkleTree{
			root:   nil,
			leaves: make([]*MerkleNode, 0),
			data:   make([][]byte, 0),
			count:  0,
		}
	}

	mt := &MerkleTree{
		data:   make([][]byte, len(data)),
		leaves: make([]*MerkleNode, 0, len(data)),
		count:  int64(len(data)),
	}

	// 复制数据
	copy(mt.data, data)

	// 构建叶子节点
	leaves := make([]*MerkleNode, len(data))
	for i, d := range data {
		leaves[i] = NewMerkleNode(d, nil, nil)
		mt.leaves = append(mt.leaves, leaves[i])
	}

	// 递归构建树
	mt.root = buildMerkleTree(leaves)

	return mt
}

// buildMerkleTree 递归构建默克尔树
func buildMerkleTree(nodes []*MerkleNode) *MerkleNode {
	if len(nodes) == 1 {
		return nodes[0]
	}

	nextLevel := make([]*MerkleNode, 0, (len(nodes)+1)/2)

	for i := 0; i < len(nodes); i += 2 {
		right := i + 1
		if right >= len(nodes) {
			// 奇数个节点，最后一个节点复制
			right = i
		}

		parent := NewMerkleNode(nil, nodes[i], nodes[right])
		nextLevel = append(nextLevel, parent)
	}

	return buildMerkleTree(nextLevel)
}

// VerifyData 验证单个数据块的完整性
func (mt *MerkleTree) VerifyData(index int, data []byte) bool {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	if index < 0 || index >= len(mt.leaves) {
		return false
	}

	expectedHash := mt.leaves[index].hash
	actualHash := sha256.Sum256(data)
	actualHashStr := hex.EncodeToString(actualHash[:])

	return expectedHash == actualHashStr
}

// VerifyRoot 验证根哈希
func (mt *MerkleTree) VerifyRoot(expectedRootHash string) bool {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	return mt.root != nil && mt.root.hash == expectedRootHash
}

// GetProof 获取数据块的完整性证明
// 返回从该叶子节点到根的所有兄弟节点的哈希值
func (mt *MerkleTree) GetProof(index int) ([][]byte, []string, error) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	if index < 0 || index >= len(mt.leaves) {
		return nil, nil, fmt.Errorf("index out of range")
	}

	var proof []string
	var hashes [][]byte

	node := mt.leaves[index]

	// 从叶子节点向上遍历到根节点
	for node.parent != nil {
		parent := node.parent
		var sibling *MerkleNode

		// 找到兄弟节点
		if len(parent.children) == 2 {
			if parent.children[0] == node {
				sibling = parent.children[1]
			} else {
				sibling = parent.children[0]
			}

			proof = append(proof, sibling.hash)
			hashes = append(hashes, []byte(sibling.hash))
		}

		node = parent
	}

	// 反转数组（从叶子到根改为从根到叶子）
	for i, j := 0, len(proof)-1; i < j; i, j = i+1, j-1 {
		proof[i], proof[j] = proof[j], proof[i]
		hashes[i], hashes[j] = hashes[j], hashes[i]
	}

	return hashes, proof, nil
}

// VerifyProof 验证完整性证明
// proof: 兄弟节点哈希值数组（从根到叶子）
// targetHash: 目标数据的哈希值
// rootHash: 期望的根哈希值
func VerifyProof(data []byte, proof []string, rootHash string) bool {
	// 计算数据块的哈希值
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	// 从叶子开始向上验证
	currentHash := hashStr

	for _, siblingHash := range proof {
		// 连接当前哈希和兄弟节点哈希
		combined := currentHash + siblingHash
		combinedHash := sha256.Sum256([]byte(combined))
		currentHash = hex.EncodeToString(combinedHash[:])
	}

	return currentHash == rootHash
}

// GetRootHash 获取根哈希值
func (mt *MerkleTree) GetRootHash() string {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	if mt.root == nil {
		return ""
	}
	return mt.root.hash
}

// RangeQuery 范围查询，返回指定范围内的数据块
func (mt *MerkleTree) RangeQuery(start, end int) ([][]byte, error) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	if start < 0 || end > len(mt.data) || start >= end {
		return nil, fmt.Errorf("invalid range")
	}

	result := make([][]byte, 0, end-start)
	for i := start; i < end; i++ {
		result = append(result, mt.data[i])
	}

	return result, nil
}

// GetAllData 获取所有数据
func (mt *MerkleTree) GetAllData() [][]byte {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	result := make([][]byte, len(mt.data))
	copy(result, mt.data)
	return result
}

// UpdateData 更新指定索引的数据块
func (mt *MerkleTree) UpdateData(index int, newData []byte) error {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	if index < 0 || index >= len(mt.leaves) {
		return fmt.Errorf("index out of range")
	}

	// 更新数据
	mt.data[index] = newData

	// 重新计算从叶子节点到根节点的哈希值
	node := mt.leaves[index]
	node.data = newData
	node.hash = node.computeHash()

	// 向上更新父节点
	for node.parent != nil {
		node = node.parent
		node.hash = node.computeHash()
	}

	mt.root = node

	return nil
}

// Size 返回数据块数量
func (mt *MerkleTree) Size() int64 {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	return mt.count
}

// Height 返回树的高度
func (mt *MerkleTree) Height() int {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	if mt.root == nil {
		return 0
	}

	height := 0
	node := mt.root
	for node != nil {
		height++
		if len(node.children) > 0 {
			node = node.children[0]
		} else {
			break
		}
	}
	return height
}

// String 返回默克尔树的字符串表示（用于调试）
func (mt *MerkleTree) String() string {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	var result string
	result += fmt.Sprintf("MerkleTree(count=%d, root=%s):\n", mt.count, mt.root.hash)
	result += mt.root.String("")
	return result
}

// MerkleNode的String方法
func (n *MerkleNode) String(prefix string) string {
	result := prefix

	if n.isLeaf {
		result += fmt.Sprintf("Leaf(hash=%s, data=%x)\n", n.hash, n.data)
	} else {
		result += fmt.Sprintf("Node(hash=%s)\n", n.hash)
		for _, child := range n.children {
			result += child.String(prefix + "  ")
		}
	}

	return result
}

// NewMerkleTreeFromKV 从键值对创建默克尔树
func NewMerkleTreeFromKV(kvs []KeyValue) *MerkleTree {
	data := make([][]byte, len(kvs))
	for i, kv := range kvs {
		// 将键值对序列化为字节数组
		kvBytes := []byte(fmt.Sprintf("%v:%v", kv.Key, kv.Value))
		data[i] = kvBytes
	}
	return NewMerkleTree(data)
}

// BinaryMerkleTree 二进制默克尔树版本
// 优化的实现，要求数据块数量为2的幂次方
type BinaryMerkleTree struct {
	root     *MerkleNode
	leaves   []*MerkleNode
	data     [][]byte
	mu       sync.RWMutex
	count    int64
}

// NewBinaryMerkleTree 创建二进制默克尔树
func NewBinaryMerkleTree(data [][]byte) (*BinaryMerkleTree, error) {
	if len(data) == 0 {
		return &BinaryMerkleTree{}, nil
	}

	// 检查是否为2的幂次方
	if len(data)&(len(data)-1) != 0 {
		return nil, fmt.Errorf("data length must be power of 2")
	}

	mt := &BinaryMerkleTree{
		data:   make([][]byte, len(data)),
		leaves: make([]*MerkleNode, 0, len(data)),
		count:  int64(len(data)),
	}

	copy(mt.data, data)

	// 构建叶子节点
	leaves := make([]*MerkleNode, len(data))
	for i, d := range data {
		leaves[i] = NewMerkleNode(d, nil, nil)
		mt.leaves = append(mt.leaves, leaves[i])
	}

	// 使用二进制方式构建树
	mt.root = mt.buildBinaryTree(leaves, 0, len(leaves))

	return mt, nil
}

// buildBinaryTree 二进制方式构建树
func (mt *BinaryMerkleTree) buildBinaryTree(nodes []*MerkleNode, start, end int) *MerkleNode {
	if end-start == 1 {
		return nodes[start]
	}

	mid := (start + end) / 2
	left := mt.buildBinaryTree(nodes, start, mid)
	right := mt.buildBinaryTree(nodes, mid, end)

	return NewMerkleNode(nil, left, right)
}
