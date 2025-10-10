package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wxnacy/go-bdpan"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "", "百度网盘路径")
	flag.Parse()

	if path == "" {
		fmt.Println("用法: go run main.go --path 百度网盘路径")
		os.Exit(1)
	}

	accessToken := os.Getenv("BDPAN_ACCESS_TOKEN")
	if accessToken == "" {
		fmt.Println("请设置BDPAN_ACCESS_TOKEN环境变量")
		os.Exit(1)
	}

	req := &bdpan.GetFileListAllReq{
		Path:      path,
		Recursion: 1, // Recursive
		Limit:     1000,
		Start:     0,
	}

	var allFiles []*bdpan.FileInfo

	for {
		res, err := bdpan.GetFileListAll(accessToken, req)
		if err != nil {
			fmt.Printf("获取文件列表失败: %v\n", err)
			os.Exit(1)
		}

		if res.Errno != 0 {
			fmt.Printf("获取文件列表API错误: %s (%d)\n", res.Error(), res.Errno)
			os.Exit(1)
		}

		allFiles = append(allFiles, res.List...)

		if res.HasMore == 0 {
			break
		}
		req.Start = res.Cursor
	}

	for _, file := range allFiles {
		fmt.Printf("fsid: %20d IsDir: %v, name: %s\n", file.FSID, file.IsDir(), file.Path)
	}
	fmt.Printf("文件列表 (总数: %d):\n", len(allFiles))
}
