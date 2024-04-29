package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-app/api/v1/route"
	"gin-app/internal/bootstrap"
	"github.com/spf13/pflag"
)

func main() {
	pflag.Parse()
	app, err := bootstrap.App(config)
	if err != nil {
		log.Fatalln(err)
	}
	defer app.Close()
	srv := route.Setup(app, time.Second*3)
	go func() {
		log.Println("Server is running on port:", app.Conf.HTTPort)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigterm
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
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
