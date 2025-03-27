#!/bin/bash
# 运行所有测试
echo "运行单元测试..."
go test -v ./internal/... -configDir=./config

## 运行指定测试
#go test -v -run TestChatFlow ./internal/handler/

# 生成测试覆盖率报告
echo "生成覆盖率报告..."
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out