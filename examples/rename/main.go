package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wxnacy/go-bdpan"
)

func main() {
	fmt.Println(os.Args)
	// if len(os.Args) != 3 {
	// fmt.Printf("Usage: %s <path> <newname>\n", os.Args[0])
	// return
	// }
	// path := os.Args[1]
	// newname := os.Args[2]
	path := "/apps/bdpan/test123_20250926_164903.go"
	newname := "testwefs.go"

	// 从环境变量获取accessToken
	accessToken := os.Getenv("BDPAN_ACCESS_TOKEN")
	if accessToken == "" {
		fmt.Println("请设置BDPAN_ACCESS_TOKEN环境变量")
		os.Exit(1)
	}

	info, err := bdpan.RenameFiles(accessToken, bdpan.NewFileManager(path, "", newname))
	if err != nil {
		log.Fatalf("RenameFiles failed: %v\n", err)
	}
	fmt.Printf("Rename task created with info: %+v\n", info)
	fmt.Printf("Rename task created with info: %+v\n", info.Error())
}
