# Makefile for 数据结构项目

.PHONY: help test bench clean example fmt vet

# 默认目标
help:
	@echo "可用命令:"
	@echo "  test      - 运行所有测试"
	@echo "  bench     - 运行基准测试"
	@echo  "  example   - 运行示例代码"
	@echo "  fmt       - 格式化代码"
	@echo "  vet       - 静态分析"
	@echo "  clean     - 清理构建文件"

# 运行所有测试
test:
	go test -v ./pkg/datastructures/...

# 运行基准测试
bench:
	go test -bench=. -benchmem ./pkg/datastructures/...

# 运行性能对比测试
bench-compare:
	go test -v -run TestPerformanceComparison ./pkg/datastructures/...

# 运行正确性测试
test-correctness:
	go test -v -run TestCorrectness ./pkg/datastructures/...

# 运行示例代码
example:
	go run examples/example_usage.go

# 格式化代码
fmt:
	go fmt ./...

# 静态分析
vet:
	go vet ./...

# 构建
build:
	go build -o bin/datastructures ./examples/example_usage.go

# 清理构建文件
clean:
	rm -rf bin/
	go clean ./...

# 运行特定数据结构的测试
test-bplus:
	go test -v -run TestBPlus ./pkg/datastructures/...

test-skiplist:
	go test -v -run TestSkipList ./pkg/datastructures/...

test-hash:
	go test -v -run TestExtendibleHash ./pkg/datastructures/...

test-bloom:
	go test -v -run TestBloom ./pkg/datastructures/...

test-merkle:
	go test -v -run TestMerkle ./pkg/datastructures/...

# 运行所有基准测试并生成报告
bench-all:
	go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof ./pkg/datastructures/...
	@echo "性能报告已生成: cpu.prof, mem.prof"

# 分析性能报告
analyze: bench-all
	go tool pprof bin/datastructures cpu.prof
	go tool pprof bin/datastructures mem.prof

# 检查代码覆盖率
coverage:
	go test -coverprofile=coverage.out ./pkg/datastructures/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 生成测试报告
test-report:
	go test -v -coverprofile=coverage.out ./pkg/datastructures/...
	go tool cover -func=coverage.out

# 一键运行完整测试套件
all: fmt vet test bench coverage
	@echo "所有测试完成!"

# 快速验证（运行基本测试）
verify: fmt vet test-correctness
	@echo "快速验证完成!"

# 安装依赖
deps:
	go mod tidy
	go mod download
