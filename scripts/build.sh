#!/bin/bash

# 自动获取系统环境变量
export DATE=$( date +%Y%m%d-%H:%M:%S)
export LATEST_COMMIT=$( git log --pretty=format:'%h' -n 1)
export BRANCH=$( git branch |grep -v "no branch"| grep \*|cut -d ' ' -f2)
export BUILT_ON_IP=$( [ $(uname) = Linux ] && hostname -i || hostname )
export RUNTIME_VER=$( go version)
export BUILT_ON_OS=$( uname -a)
if [ ${#BRANCH} -eq 0 ]; then
  BRANCH=master
fi

export COMMIT_CNT=$( git rev-list HEAD | wc -l | sed 's/ //g' )
export BUILD_NUMBER=${BRANCH}-${COMMIT_CNT}
# shellcheck disable=SC2089
# shellcheck disable=SC2016
export COMPILE_LDFLAGS='-s -w
                          -X "main.BuildDate='${DATE}'"
                          -X "main.LatestCommit='${LATEST_COMMIT}'"
                          -X "main.BuildNumber='${BUILD_NUMBER}'"
                          -X "main.BuiltOnIP='${BUILT_ON_IP}'"
                          -X "main.BuiltOnOs='${BUILT_ON_OS}'"
                          -X "main.Branch='${BRANCH}'"
                          -X "main.CommitCnt='${COMMIT_CNT}'"
                          -X "main.RuntimeVer='${RUNTIME_VER}'" '
# 构建可执行二进制
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=https://proxy.golang.com.cn,direct go build -o bin/gin-app -ldflags="${COMPILE_LDFLAGS}" ./cmd/...
# 重要提示,镜像名称在新建项目时必须自己修改!!!
# 构建docker镜像
docker build . -t gin-app/server
mkdir -p run/
# 保存docker镜像
docker save gin-app/server -o run/gin-app-server.tar.gz
# 把准备好的数据copy到run目录，方便之间迁移部署
cp -r bin run/
rsync -avuz config run/
rsync -avuz docker-compose.yaml run/
rsync -avuz Dockerfile run/
echo "#!/bin/bash
docker build . -t gin-app/server
" > run/build.sh
echo  "#!/bin/bash
docker load -i gin-app-server.tar.gz
" > run/load.sh