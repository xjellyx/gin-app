.DEFAULT_GOAL := build_all


fmt:
	for pkg in ${PACKAGES}; do \
		go fmt $$pkg; \
	done;

# 输出二进制自己修改
build:
	# 构建可以运行的目录
	sh scripts/build.sh
# 删除编译结果
clean:
	rm -rf ./bin/*

# 生成swag文档
swag:
	 swag init   -g cmd/main.go
