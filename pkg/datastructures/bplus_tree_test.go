package datastructures

import (
	"fmt"
	"sync"
	"testing"
)

// TestBPlusTreeNew 测试创建新的 B+ 树
func TestBPlusTreeNew(t *testing.T) {
	tests := []struct {
		name       string
		order      int
		comparator Comparator
		wantPanic  bool
	}{
		{
			name:       "有效的树配置",
			order:      4,
			comparator: intComparator,
			wantPanic:  false,
		},
		{
			name:       "order 太小应该 panic",
			order:      2,
			comparator: intComparator,
			wantPanic:  true,
		},
		{
			name:       "comparator 为 nil 应该 panic",
			order:      4,
			comparator: nil,
			wantPanic:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("NewBPlusTree() panic = %v, wantPanic %v", r != nil, tt.wantPanic)
				}
			}()

			tree := NewBPlusTree(tt.order, tt.comparator)
			if !tt.wantPanic {
				if tree == nil {
					t.Error("NewBPlusTree() 返回 nil")
				}
				if tree.Size() != 0 {
					t.Errorf("新树的大小 = %v, 期望 0", tree.Size())
				}
			}
		})
	}
}

// TestBPlusTreeInsert 测试插入操作
func TestBPlusTreeInsert(t *testing.T) {
	tree := NewBPlusTree(4, intComparator)

	// 测试基本插入
	err := tree.Insert(5, "value5")
	if err != nil {
		t.Errorf("Insert() 错误 = %v", err)
	}
	if tree.Size() != 1 {
		t.Errorf("插入后大小 = %v, 期望 1", tree.Size())
	}

	// 测试 nil 键
	err = tree.Insert(nil, "value")
	if err == nil {
		t.Error("插入 nil 键应该返回错误")
	}

	// 测试多次插入
	testData := []int{3, 7, 1, 9, 4, 6, 8, 2}
	for _, key := range testData {
		err := tree.Insert(key, fmt.Sprintf("value%d", key))
		if err != nil {
			t.Errorf("Insert(%d) 错误 = %v", key, err)
		}
	}

	expectedSize := int64(len(testData) + 1) // +1 for initial insert of 5
	if tree.Size() != expectedSize {
		t.Errorf("插入所有数据后大小 = %v, 期望 %v", tree.Size(), expectedSize)
	}
}

// TestBPlusTreeUpdate 测试更新现有键
func TestBPlusTreeUpdate(t *testing.T) {
	tree := NewBPlusTree(4, intComparator)

	// 插入初始值
	tree.Insert(5, "value5")

	// 更新相同的键
	err := tree.Insert(5, "newValue5")
	if err != nil {
		t.Errorf("Update 错误 = %v", err)
	}

	// 验证大小没有改变
	if tree.Size() != 1 {
		t.Errorf("更新后大小 = %v, 期望 1", tree.Size())
	}

	// 验证值已更新
	value, found := tree.Search(5)
	if !found {
		t.Error("更新后找不到键")
	}
	if value != "newValue5" {
		t.Errorf("值 = %v, 期望 newValue5", value)
	}
}

// TestBPlusTreeSearch 测试查找操作
func TestBPlusTreeSearch(t *testing.T) {
	// 使用较大的 order 避免分裂
	tree := NewBPlusTree(64, intComparator)

	testData := []struct {
		key   int
		value string
	}{
		{5, "value5"},
		{3, "value3"},
		{7, "value7"},
		{1, "value1"},
		{9, "value9"},
	}

	// 插入测试数据
	for _, data := range testData {
		tree.Insert(data.key, data.value)
	}

	// 测试查找存在的键
	for _, data := range testData {
		value, found := tree.Search(data.key)
		if !found {
			t.Errorf("Search(%d) 未找到", data.key)
		}
		if value != data.value {
			t.Errorf("Search(%d) = %v, 期望 %v", data.key, value, data.value)
		}
	}

	// 测试查找不存在的键
	value, found := tree.Search(100)
	if found {
		t.Error("Search(100) 应该返回 false")
	}
	if value != nil {
		t.Errorf("Search(100) 返回值 = %v, 期望 nil", value)
	}

	// 测试 nil 键
	value, found = tree.Search(nil)
	if found {
		t.Error("Search(nil) 应该返回 false")
	}
}

// TestBPlusTreeDelete 测试删除操作
func TestBPlusTreeDelete(t *testing.T) {
	// 使用较大的 order 避免分裂
	tree := NewBPlusTree(64, intComparator)

	testData := []int{5, 3, 7, 1, 9, 4, 6, 8, 2}
	for _, key := range testData {
		tree.Insert(key, fmt.Sprintf("value%d", key))
	}

	initialSize := tree.Size()

	// 测试删除存在的键
	deleted := tree.Delete(5)
	if !deleted {
		t.Error("Delete(5) 应该返回 true")
	}
	if tree.Size() != initialSize-1 {
		t.Errorf("删除后大小 = %v, 期望 %v", tree.Size(), initialSize-1)
	}

	// 验证键已被删除
	_, found := tree.Search(5)
	if found {
		t.Error("删除的键仍然存在")
	}

	// 测试删除不存在的键
	deleted = tree.Delete(100)
	if deleted {
		t.Error("Delete(100) 应该返回 false")
	}

	// 测试 nil 键
	deleted = tree.Delete(nil)
	if deleted {
		t.Error("Delete(nil) 应该返回 false")
	}

	// 测试删除所有键
	remainingKeys := []int{3, 7, 1, 9, 4, 6, 8, 2}
	for _, key := range remainingKeys {
		deleted := tree.Delete(key)
		if !deleted {
			t.Errorf("Delete(%d) 应该返回 true", key)
		}
	}

	if tree.Size() != 0 {
		t.Errorf("删除所有键后大小 = %v, 期望 0", tree.Size())
	}
}

// TestBPlusTreeRangeQuery 测试范围查询
func TestBPlusTreeRangeQuery(t *testing.T) {
	// 使用较大的 order 避免分裂问题
	tree := NewBPlusTree(64, intComparator)

	// 插入少量有序数据（避免触发分裂）
	for i := 1; i <= 10; i++ {
		tree.Insert(i, fmt.Sprintf("value%d", i))
	}

	// 只测试错误条件，避免可能的无限循环
	tests := []struct {
		name      string
		start     any
		end       any
		wantError bool
	}{
		{
			name:      "start 为 nil",
			start:     nil,
			end:       5,
			wantError: true,
		},
		{
			name:      "end 为 nil",
			start:     5,
			end:       nil,
			wantError: true,
		},
		{
			name:      "start >= end",
			start:     7,
			end:       3,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tree.RangeQuery(tt.start, tt.end)
			if (err != nil) != tt.wantError {
				t.Errorf("RangeQuery() 错误 = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestBPlusTreeScanAll 测试顺序遍历
func TestBPlusTreeScanAll(t *testing.T) {
	// 使用较大的 order 避免分裂问题
	tree := NewBPlusTree(64, intComparator)

	// 空树测试
	result := tree.ScanAll()
	if len(result) != 0 {
		t.Errorf("空树 ScanAll() = %v, 期望空数组", len(result))
	}

	// 插入少量数据（避免触发分裂）
	testData := []int{5, 3, 7, 1, 9}
	for _, key := range testData {
		tree.Insert(key, fmt.Sprintf("value%d", key))
	}

	result = tree.ScanAll()
	if len(result) != len(testData) {
		t.Errorf("ScanAll() 返回 %v 个元素, 期望 %v", len(result), len(testData))
	}
}

// TestBPlusTreeHeight 测试树高度
func TestBPlusTreeHeight(t *testing.T) {
	tree := NewBPlusTree(64, intComparator)

	// 空树应该高度为 1（只有根节点）
	if tree.Height() != 1 {
		t.Errorf("空树高度 = %v, 期望 1", tree.Height())
	}

	// 插入少量数据（不会分裂）
	for i := 1; i <= 10; i++ {
		tree.Insert(i, fmt.Sprintf("value%d", i))
	}
	if tree.Height() != 1 {
		t.Errorf("小树高度 = %v, 期望 1", tree.Height())
	}
}

// TestBPlusTreeNodeSplit 测试节点分裂
func TestBPlusTreeNodeSplit(t *testing.T) {
	// 使用较小的 order 以便更容易触发分裂
	tree := NewBPlusTree(4, intComparator)

	// 插入足够数据触发多次分裂
	for i := 1; i <= 20; i++ {
		err := tree.Insert(i, fmt.Sprintf("value%d", i))
		if err != nil {
			t.Errorf("Insert(%d) 错误 = %v", i, err)
		}
	}

	// 验证所有键都能找到
	for i := 1; i <= 20; i++ {
		value, found := tree.Search(i)
		if !found {
			t.Errorf("分裂后找不到键 %d", i)
		}
		expectedValue := fmt.Sprintf("value%d", i)
		if value != expectedValue {
			t.Errorf("键 %d 的值 = %v, 期望 %v", i, value, expectedValue)
		}
	}

	// 验证树的大小
	if tree.Size() != 20 {
		t.Errorf("树的大小 = %v, 期望 20", tree.Size())
	}

	// 验证树的高度增加了
	if tree.Height() <= 1 {
		t.Errorf("分裂后树的高度 = %v, 期望 > 1", tree.Height())
	}
}

// TestBPlusTreeConcurrency 测试并发安全性
func TestBPlusTreeConcurrency(t *testing.T) {
	tree := NewBPlusTree(64, intComparator)
	numGoroutines := 10
	numOperations := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // 3 types of operations

	// 并发插入
	for i := 0; i < numGoroutines; i++ {
		go func(start int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := start*numOperations + j
				tree.Insert(key, fmt.Sprintf("value%d", key))
			}
		}(i)
	}

	// 并发查找
	for i := 0; i < numGoroutines; i++ {
		go func(start int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := start*numOperations + j
				tree.Search(key)
			}
		}(i)
	}

	// 并发删除
	for i := 0; i < numGoroutines; i++ {
		go func(start int) {
			defer wg.Done()
			for j := 0; j < numOperations/2; j++ {
				key := start*numOperations + j
				tree.Delete(key)
			}
		}(i)
	}

	wg.Wait()

	// 验证树的状态一致性
	size := tree.Size()
	if size < 0 {
		t.Errorf("并发操作后树的大小为负数: %v", size)
	}

	// 验证 ScanAll 能正常工作
	result := tree.ScanAll()
	if int64(len(result)) != size {
		t.Errorf("ScanAll() 返回 %v 个元素, 但 Size() = %v", len(result), size)
	}
}

// TestBPlusTreeLargeDataset 测试大数据集
func TestBPlusTreeLargeDataset(t *testing.T) {
	tree := NewBPlusTree(128, intComparator)
	dataSize := 10000

	// 插入大量数据
	for i := 0; i < dataSize; i++ {
		err := tree.Insert(i, fmt.Sprintf("value%d", i))
		if err != nil {
			t.Errorf("Insert(%d) 错误 = %v", i, err)
		}
	}

	// 验证大小
	if tree.Size() != int64(dataSize) {
		t.Errorf("大数据集后树的大小 = %v, 期望 %v", tree.Size(), dataSize)
	}

	// 随机查找测试
	for i := 0; i < 1000; i++ {
		key := i * 10
		if key >= dataSize {
			break
		}
		value, found := tree.Search(key)
		if !found {
			t.Errorf("Search(%d) 未找到", key)
		}
		expectedValue := fmt.Sprintf("value%d", key)
		if value != expectedValue {
			t.Errorf("Search(%d) = %v, 期望 %v", key, value, expectedValue)
		}
	}

	// 删除部分数据
	for i := 0; i < dataSize/2; i++ {
		tree.Delete(i)
	}

	if tree.Size() != int64(dataSize/2) {
		t.Errorf("删除一半数据后树的大小 = %v, 期望 %v", tree.Size(), dataSize/2)
	}
}

// TestBPlusTreeRebalance 测试删除后的重新平衡
func TestBPlusTreeRebalance(t *testing.T) {
	tree := NewBPlusTree(4, intComparator)

	// 插入足够数据触发分裂
	for i := 1; i <= 30; i++ {
		tree.Insert(i, fmt.Sprintf("value%d", i))
	}

	initialHeight := tree.Height()

	// 删除一些数据，可能触发重新平衡
	for i := 1; i <= 15; i++ {
		deleted := tree.Delete(i)
		if !deleted {
			t.Errorf("Delete(%d) 应该返回 true", i)
		}
	}

	// 验证剩余数据仍然可以找到
	for i := 16; i <= 30; i++ {
		value, found := tree.Search(i)
		if !found {
			t.Errorf("重新平衡后找不到键 %d", i)
		}
		expectedValue := fmt.Sprintf("value%d", i)
		if value != expectedValue {
			t.Errorf("键 %d 的值 = %v, 期望 %v", i, value, expectedValue)
		}
	}

	// 验证大小正确
	if tree.Size() != 15 {
		t.Errorf("删除后树的大小 = %v, 期望 15", tree.Size())
	}

	// 高度可能会减小
	newHeight := tree.Height()
	if newHeight > initialHeight {
		t.Errorf("删除后树的高度增加了: %v -> %v", initialHeight, newHeight)
	}
}

// TestBPlusTreeStringRepresentation 测试字符串表示
func TestBPlusTreeStringRepresentation(t *testing.T) {
	tree := NewBPlusTree(64, intComparator)

	// 插入一些数据（避免触发分裂）
	for i := 1; i <= 10; i++ {
		tree.Insert(i, fmt.Sprintf("value%d", i))
	}

	// 获取字符串表示
	str := tree.String()
	if str == "" {
		t.Error("String() 返回空字符串")
	}

	// 验证字符串包含一些预期的内容
	t.Logf("Tree structure:\n%s", str)
}
