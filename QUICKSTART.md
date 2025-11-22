# 快速开始指南

## 🚀 快速验证

运行演示脚本查看所有数据结构的工作效果：

```bash
cd BTreeVsHashMap
go run examples/demo.go
```

期望输出：
```
=== 数据结构演示 ===

--- B+树演示 ---
✓ 插入 6 个键值对
✓ 等值查询: Search(100) = value_100
✓ 范围查询: RangeQuery(50, 200) 返回 3 条记录
✓ 树高度: 1
✓ 数据量: 6

--- 跳表演示 ---
✓ 插入 5 个键值对
✓ 等值查询: Search(100) = value_100
✓ 范围查询: RangeQuery(50, 200) 返回 3 条记录
✓ 跳表层数: 4
✓ 数据量: 5

--- 可扩展哈希演示 ---
✓ 插入 6 个键值对
✓ 等值查询: Search(100) = value_100
✓ 桶使用统计: 平均=3.50, 最大=4, 最小=2, 满桶数=3
✓ 全局深度: 2
✓ 数据量: 6

--- 布隆过滤器演示 ---
添加元素:
  - 添加键 100
  - 添加键 200
  - 添加键 50
  - 添加键 150
  - 添加键 300

检查元素存在性:
  - 键 100: 可能存在
  - 键 150: 可能存在
  - 键 999: 一定不存在
✓ 当前假阳性率: 0.0000%
✓ 元素数量: 5

--- 默克尔树演示 ---
✓ 根哈希: 033815f6539db95258e0a46d16ae7f8b5d9637c4bdb16b87b3b91b6a5fb18465

验证数据完整性:
  - 文件 file1: ✓ 有效
  - 文件 file2: ✓ 有效
  - 文件 file3: ✓ 有效
  - 文件 file4: ✓ 有效
✓ 完整性证明长度: 2
✓ 数据块数量: 4
✓ 树高度: 3
```

## 📚 学习路径

### 1. 入门指南
- 📖 阅读 [README.md](README.md) 了解所有数据结构
- 🎬 运行 `examples/demo.go` 查看效果
- 💡 查看 [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) 了解项目结构

### 2. 深入学习
- 📝 学习每个数据结构的实现 (`pkg/datastructures/*.go`)
- 🧪 运行基准测试：`go test -bench=. -benchmem ./pkg/datastructures/`
- 🔍 查看完整示例：`examples/example_usage.go`

### 3. 实际应用
- ⚙️ 在你的项目中导入所需的数据结构
- 📊 根据查询模式选择合适的数据结构
- 🎯 参考最佳实践指南

## 📋 项目文件说明

```
BTreeVsHashMap/
├── README.md                    # 完整文档
├── PROJECT_SUMMARY.md          # 项目总结
├── IMPLEMENTATION_REPORT.md    # 实现报告
├── QUICKSTART.md               # 本文件
├── go.mod                      # Go模块
├── Makefile                    # 构建命令
│
├── pkg/datastructures/         # 核心实现
│   ├── bplus_tree.go          # B+树
│   ├── skip_list.go           # 跳表
│   ├── merkle_tree.go         # 默克尔树
│   ├── extendible_hash.go     # 可扩展哈希
│   ├── bloom_filter.go        # 布隆过滤器
│   └── benchmark_test.go      # 测试
│
└── examples/                   # 示例
    ├── example_usage.go       # 完整示例
    └── demo.go                # 演示
```

## 🎯 数据结构选择指南

### 需要范围查询？
- ✅ 是 → 选择 **B+树** 或 **跳表**
- ❌ 否 → 进入下一步

### 磁盘存储还是内存？
- 💾 磁盘 → 选择 **B+树** 或 **可扩展哈希**
- 🧠 内存 → 选择 **跳表**

### 主要是等值查询？
- ✅ 是 → 选择 **可扩展哈希** 或 **跳表**
- ❌ 否 → 选择 **B+树**

### 需要验证数据完整性？
- ✅ 是 → 选择 **默克尔树**

### 需要快速判断存在性（允许误判）？
- ✅ 是 → 选择 **布隆过滤器**

### 同时需要等值查询和范围查询？
- ✅ 是 → 选择 **跳表**（最佳平衡）

## 🛠️ 常用命令

```bash
# 运行演示
go run examples/demo.go

# 运行所有测试
make test
go test -v ./pkg/datastructures/

# 运行基准测试
make bench
go test -bench=. -benchmem ./pkg/datastructures/

# 查看代码覆盖率
make coverage

# 格式化代码
make fmt

# 构建示例
make build
```

## 📞 获取帮助

- 📖 查看 README.md 获取详细文档
- 🧪 运行测试了解功能
- 📊 运行基准测试了解性能
- 💡 查看示例代码学习用法

---

**开始探索数据结构的世界吧！** 🚀
