package main

import (
	"fmt"
	"os"

	"github.com/WAY29/onefinger/core"
	"github.com/WAY29/onefinger/structs"
	"github.com/WAY29/onefinger/utils"

	cli "github.com/jawher/mow.cli"
	_ "github.com/panjf2000/ants"
)

var (
	app         *cli.Cli
	__version__ = "1.0.2"
)

func main() {
	app = cli.App("onefinger", "Simple website fingerprinting tool")

	// 定义选项
	var (
		target      = app.StringsOpt("t target", make([]string, 0), "Target url")
		targetFiles = app.StringsOpt("tf", make([]string, 0), "Target url file")
		threads     = app.IntOpt("threads", 10, "Thread number")
		timeout     = app.IntOpt("timeout", 20, "Request timeout")
	)
	// 定义版本
	app.Version("v version", "onefinger 1.1.0")
	// 命令行描述
	app.Spec = "-v | (-t=<target> | --tf=<targetFile>)... [--threads=<threads>] [--timeout=<timeout>]"

	app.Action = func() {
		targetsSlice := make([]string, 0)
		if len(*target) != 0 {
			targetsSlice = append(targetsSlice, *target...)
		}
		for _, targetFile := range *targetFiles {
			if utils.Exists(targetFile) && utils.IsFile(targetFile) {
				err, lineSlice := utils.ReadFileAsLine(targetFile)
				if err != nil {
					utils.OptionsError("Targetfile handle error: "+err.Error(), 2)
				}
				targetsSlice = append(targetsSlice, lineSlice...)
			}
		}

		if len(targetsSlice) == 0 {
			utils.OptionsError("Target is empty", 2)
		}

		// 目标切片去重
		targetsSlice = utils.RemoveDuplicateElement(targetsSlice)

		// 主要运行逻辑
		fmt.Printf("%v\n", targetsSlice)
		os.Exit(0)

		// 最大线程数为目标数量
		if *threads > len(targetsSlice) {
			*threads = len(targetsSlice)
		}

		core.Initialize(targetsSlice, *threads, structs.RequestOptions{
			Timeout: *timeout,
		})
		core.Start(targetsSlice)
		core.Wait()
		core.End()
	}

	app.Run(os.Args)
}
