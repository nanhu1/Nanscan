package main

import (
	"flag"
	"fmt"
	"nanscan/fileutil"
	"nanscan/flag_new"
	"nanscan/json_core"
	"nanscan/request"
	"sync"
)

func main() {
	flag_new.Banner()

	var (
		Url     string
		File    string
		Threads int
		WG      sync.WaitGroup
	)

	fmt.Println("(-u=<targetUrl> | -f=<target File> | -t=<threads>)")

	flag.StringVar(&Url, "u", "", "输入url")
	flag.StringVar(&File, "f", "", "文件内为url")
	flag.IntVar(&Threads, "t", 5, "线程默认为5")

	flag.Parse()

	targetsSlice := make([]string, 0)
	if Threads > len(targetsSlice) {
		Threads = len(targetsSlice)
	}

	if Url != "" && File == "" {
		req, _ := request.Reqdata(Url)
		json_core.Fetchbody(req)
	} else if Url == "" && File != "" {
		filename, _ := fileutil.ReadFile(File)
		targetsSlice = filename

		for i := 0; i < len(filename); i++ {
			WG.Add(1)
			go func(url string) {
				// 使用defer, 表示函数完成时将等待组值减1
				defer WG.Done()
				req, _ := request.Reqdata(url)
				json_core.Fetchbody(req)

			}(filename[i])
		}
		WG.Wait()
	}
}
