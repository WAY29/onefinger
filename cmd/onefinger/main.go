package main

import (
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
		target = app.StringsOpt("t target", make([]string, 0), "Target url")
		// target     = app.StringOpt("t target", "", "Target url")
		targetFile = app.StringOpt("tf", "", "Target url file")
		threads    = app.IntOpt("threads", 10, "Thread number")
		timeout    = app.IntOpt("timeout", 20, "Request timeout")
	)
	// 定义版本
	app.Version("v version", "onefinger 1.0.0")
	// 命令行描述
	app.Spec = "[-v] | (-t=<target>... | --tf=<targetFile>) [--threads=<threads>] [--timeout=<timeout>]"

	app.Action = func() {
		targetsSlice := make([]string, 0)
		if len(*target) != 0 {
			targetsSlice = append(targetsSlice, *target...)
		} else {
			if utils.IsFile(*targetFile) {
				err, lineSlice := utils.ReadFileAsLine(*targetFile)
				if err != nil {
					utils.OptionsError("Targetfile handle error: "+err.Error(), 2)
				}
				targetsSlice = append(targetsSlice, lineSlice...)
			} else {
				utils.OptionsError("Target and targetFile empty", 2)
			}
		}

		// 主要运行逻辑

		// 最大线程数为目标数量
		if *threads > len(targetsSlice) {
			*threads = len(targetsSlice)
		}

		// todo remove: 输出测试
		// for _, t := range targetsSlice {
		// 	fmt.Println("target: " + t)
		// }
		// fmt.Println("thread number: " + strconv.Itoa(*threads))

		core.Initialize(targetsSlice, *threads, structs.RequestOptions{
			Timeout: *timeout,
		})
		core.Start(targetsSlice)
		core.Wait()
		core.End()
	}

	app.Run(os.Args)
}
