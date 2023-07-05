.DEFAULT_GOAL := build_all

# 自动获取系统环境变量
export DATE := $(shell date +%Y%m%d-%H:%M:%S)
export LATEST_COMMIT := $(shell git log --pretty=format:'%h' -n 1)
export BRANCH := $(shell git branch |grep -v "no branch"| grep \*|cut -d ' ' -f2)
export BUILT_ON_IP := $(shell [ $$(uname) = Linux ] && hostname -i || hostname )
export RUNTIME_VER := $(shell go version)
export BUILT_ON_OS=$(shell uname -a)
ifeq ($(BRANCH),)
BRANCH := master
endif


export COMMIT_CNT := $(shell git rev-list HEAD | wc -l | sed 's/ //g' )
export BUILD_NUMBER := ${BRANCH}-${COMMIT_CNT}
export COMPILE_LDFLAGS='-s -w \
                          -X "main.BuildDate=${DATE}" \
                          -X "main.LatestCommit=${LATEST_COMMIT}" \
                          -X "main.BuildNumber=${BUILD_NUMBER}" \
                          -X "main.BuiltOnIP=${BUILT_ON_IP}" \
                          -X "main.BuiltOnOs=${BUILT_ON_OS}" \
                          -X "main.Branch=${BRANCH}" \
                          -X "main.CommitCnt=${COMMIT_CNT}" \
                          -X "main.RuntimeVer=${RUNTIME_VER}" '

fmt:
	for pkg in ${PACKAGES}; do \
		go fmt $$pkg; \
	done;

# 输出二进制自己修改
build_all:
	# 构建二进制
	echo $COMPILE_LDFLAGS
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=https://proxy.golang.com.cn,direct go build -o bin/gin-app -ldflags=${COMPILE_LDFLAGS} ./cmd/...
# 删除编译结果
clean:
	rm -rf ./bin/*

# 生成swag文档
swag:
	 swag init --pd --parseDepth=3
