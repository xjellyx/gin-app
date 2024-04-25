package main

import (
	"fmt"
	"gin-app/api/v1/route"
	"gin-app/internal/bootstrap"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	pflag.Parse()
	app, err := bootstrap.App(config)
	if err != nil {
		log.Fatalln(err)
	}
	defer app.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		route.Setup(app, time.Second*3)
		wg.Done()
	}()
	wg.Wait()
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigterm
	}()

}

var (
	config string
)

func init() {
	pflag.StringVar(&config, "c", "config/config.yaml", "choose config file.")
	pflag.BoolVar(&ShowVersion, "v", false, "version")
	// 配置文件路径
	pflag.Int("HTTP_PORT", 8818, "")
	// 数据库
	pflag.String("DB_DRIVER", "postgresql", "")
	pflag.String("DB_DSN", "postgres://postgres:123456@localhost:5432/public?sslmode=disable&TimeZone=Asia/Shanghai", "")
	pflag.Bool("DB_AUTO_MIGRATE", false, "")
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "-version") {
		fmt.Println(version())
		os.Exit(1)
	}
}

// 版本信息
var (
	// BuildDate date string of when build was performed filled in by -X compile flag
	BuildDate string
	// LatestCommit date string of when build was performed filled in by -X compile flag
	LatestCommit string
	// BuildNumber date string of when build was performed filled in by -X compile flag
	BuildNumber string
	// BuiltOnIP date string of when build was performed filled in by -X compile flag
	BuiltOnIP string
	// BuiltOnOs date string of when build was performed filled in by -X compile flag
	BuiltOnOs string
	// RuntimeVer date string of when build was performed filled in by -X compile flag
	RuntimeVer string
	// Branch git branch
	Branch string
	// CommitCnt ...
	CommitCnt string
	// ShowVersion 展示版本号
	ShowVersion bool
)

func version() string {
	v := fmt.Sprintf(`
---------------------------------------------------------
BUILT_ON_IP       %s
BUILT_ON_OS       %s
DATE              %s
LATEST_COMMIT     %s
BRANCH            %s
COMMIT_CNT        %s
BUILD_NUMBER      %s
RUNTIME_VER       %s
---------------------------------------------------------
	`, BuiltOnIP, BuiltOnOs, BuildDate, LatestCommit, Branch, CommitCnt, BuildNumber, RuntimeVer)

	return v
}
