package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-app/api/v1/route"
	"gin-app/internal/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {
	flag.Parse()
	app, err := bootstrap.App(config)
	if err != nil {
		log.Fatalln(err)
	}
	defer app.Close()
	g := gin.Default()
	route.Setup(app.Conf, time.Second*3, app.Database, app.Log, g)
	g.Run(":8080")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigterm
	}()

}

var (
	config string
)

func init() {
	flag.StringVar(&config, "c", "config.yaml", "choose config file.")
	flag.BoolVar(&ShowVersion, "v", false, "version")

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
