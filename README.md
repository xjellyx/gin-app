#go-project-layout

# 项目目录结构
[参考](https://github.com/golang-standards/project-layout/blob/master/README_zh.md)

# 项目构建
```shell
# 执行 make build ,会自动生产run目录，该目录是服务运行所需的文件
make build
cd run && docker-compose up -d

```