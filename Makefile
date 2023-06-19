.PHONY: local test cover cover-func cover-html clean

GO=go1.20.4

local:
	$(GO) build -ldflags '-w -s' -o bin/util

# 执行单元测试
test:
	$(GO) test ./...

# 执行benchmark测试
bench:
	$(GO) test -bench . ./... -run=none

# 统计覆盖率
cover:
	$(GO) test ./... -coverprofile cover.profile

# 打开浏览器显示覆盖统计信息
cover-html:
	$(GO) tool cover -html=cover.profile

# 显示函数覆盖统计
cover-func:
	$(GO) tool cover -func=cover.profile

# 清理
clean:
	rm -rf ./bin \
	rm cover.profile
