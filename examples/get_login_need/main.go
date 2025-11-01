package main

import (
	"fmt"
	"os"

	"github.com/wxnacy/go-bdpan"
)

func main() {
	appKey := os.Getenv("BDPAN_APP_KEY")
	secretKey := os.Getenv("BDPAN_SECRET_KEY")
	codeRes, _ := bdpan.GetDeviceCode(appKey, "basic,netdisk")
	fmt.Printf("GetDeviceCode res: %#v\n", codeRes)
	tokenRes, err := bdpan.GetDeviceToken(appKey, secretKey, codeRes.DeviceCode)
	fmt.Printf("GetDeviceToken res: %#v\n", tokenRes)
	fmt.Printf("GetDeviceToken err: %#v\n", err)
}
