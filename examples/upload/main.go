package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/wxnacy/go-bdpan"
	"github.com/wxnacy/go-tools"
)

const (
	// 分片大小，4MB，普通用户的最大分片大小
	ChunkSize = 4 * 1024 * 1024
)

// https://pan.baidu.com/union/doc/3ksg0s9ye
func main() {
	// 从命令行参数获取本地文件路径和百度网盘路径
	var from string
	var to string
	flag.StringVar(&from, "from", "", "本地文件路径")
	flag.StringVar(&to, "to", "", "百度网盘路径")
	flag.Parse()

	// 检查参数
	if from == "" || to == "" {
		fmt.Println("用法: go run main.go --from 本地文件路径 --to 百度网盘路径")
		os.Exit(1)
	}

	// 从环境变量获取accessToken
	accessToken := os.Getenv("BDPAN_ACCESS_TOKEN")
	if accessToken == "" {
		fmt.Println("请设置BDPAN_ACCESS_TOKEN环境变量")
		os.Exit(1)
	}

	// 上传文件
	if err := uploadFile(accessToken, from, to); err != nil {
		fmt.Printf("上传失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("文件 %s 上传成功，保存路径为 %s\n", from, to)
}

// uploadFile 实现文件上传的完整流程
func uploadFile(accessToken, localFilePath, remoteFilePath string) error {
	// 1. 打开本地文件
	file, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 2. 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 3. 计算文件大小
	fileSize := fileInfo.Size()
	fileMD5, err := tools.Md5File(localFilePath)
	if err != nil {
		return err
	}
	fmt.Printf("文件md5: %s\n", fileMD5)

	// 4. 计算文件的MD5分块列表
	blockList, err := calculateBlockList(file, fileSize)
	if err != nil {
		return fmt.Errorf("计算文件MD5失败: %w", err)
	}

	// 5. 预上传
	preCreateReq := bdpan.NewPreCreateFileReq(remoteFilePath, int32(fileSize), blockList)
	preCreateReq.ContentMD5 = fileMD5
	preCreateRes, err := bdpan.PreCreateFile(accessToken, preCreateReq)
	if err != nil {
		return fmt.Errorf("预上传失败: %w", err)
	}

	if preCreateRes.Errno != 0 {
		return fmt.Errorf("预上传失败，错误码: %d", preCreateRes.Errno)
	}

	fmt.Printf("预上传成功，uploadid: %s\n", preCreateRes.Uploadid)

	// 6. 分片上传
	// 如果文件小于等于4MB，只需要上传一个分片
	// 如果文件大于4MB，需要按照4MB大小分片上传
	var remoteBlockList []string

	if fileSize <= ChunkSize {
		// 小文件上传
		if _, err := file.Seek(0, 0); err != nil {
			return fmt.Errorf("文件指针重置失败: %w", err)
		}

		uploadPartReq := bdpan.NewUploadFilePartReq(remoteFilePath, preCreateRes.Uploadid, 0)
		uploadPartReq.File = file

		uploadPartRes, err := bdpan.UploadFilePart(accessToken, uploadPartReq)
		if err != nil {
			return fmt.Errorf("分片上传失败: %w", err)
		}

		if uploadPartRes.Errno != 0 {
			return fmt.Errorf("分片上传失败，错误码: %d", uploadPartRes.Errno)
		}

		remoteBlockList = append(remoteBlockList, uploadPartRes.Md5)
		fmt.Printf("分片 0 上传成功，md5: %s\n", uploadPartRes.Md5)
	} else {
		// 大文件分片上传
		chunkCount := int(fileSize / ChunkSize)
		if fileSize%ChunkSize != 0 {
			chunkCount++
		}

		for i := range chunkCount {
			if _, err := file.Seek(int64(i)*ChunkSize, 0); err != nil {
				return fmt.Errorf("文件指针定位失败: %w", err)
			}

			// 创建临时文件存储分片数据
			tempFile, err := os.CreateTemp("", "upload_chunk_*")
			if err != nil {
				return fmt.Errorf("创建临时文件失败: %w", err)
			}
			tempFilePath := tempFile.Name()
			defer os.Remove(tempFilePath)

			// 读取分片数据
			writer := bufio.NewWriter(tempFile)
			r := io.LimitReader(file, ChunkSize)
			if _, err := io.Copy(writer, r); err != nil {
				tempFile.Close()
				return fmt.Errorf("写入分片数据失败: %w", err)
			}
			writer.Flush()
			tempFile.Close()

			// 重新打开临时文件用于上传
			tempFile, err = os.Open(tempFilePath)
			if err != nil {
				return fmt.Errorf("打开临时文件失败: %w", err)
			}
			defer tempFile.Close()

			// 上传分片
			uploadPartReq := bdpan.NewUploadFilePartReq(remoteFilePath, preCreateRes.Uploadid, i)
			uploadPartReq.File = tempFile

			uploadPartRes, err := bdpan.UploadFilePart(accessToken, uploadPartReq)
			if err != nil {
				tempFile.Close()
				return fmt.Errorf("分片 %d 上传失败: %w", i, err)
			}
			tempFile.Close()

			if uploadPartRes.Errno != 0 {
				return fmt.Errorf("分片 %d 上传失败，错误码: %d", i, uploadPartRes.Errno)
			}

			remoteBlockList = append(remoteBlockList, uploadPartRes.Md5)
			fmt.Printf("分片 %d/%d 上传成功，md5: %s path: %s\n", i+1, chunkCount, uploadPartRes.Md5, tempFilePath)
		}
	}

	// 7. 创建文件
	createFileReq := bdpan.NewCreateFileReq(remoteFilePath, int32(fileSize), 0, preCreateRes.Uploadid, remoteBlockList)
	createFileRes, err := bdpan.CreateFile(accessToken, createFileReq)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}

	if createFileRes.Errno != 0 {
		return fmt.Errorf("创建文件失败，错误码: %d", createFileRes.Errno)
	}

	fmt.Printf("文件创建成功，fs_id: %d\n", createFileRes.FSId)
	fmt.Printf("文件创建成功，md5: %s\n", createFileRes.Md5)

	file_res, err := bdpan.GetFileInfo(accessToken, bdpan.NewGetFileInfoReq(createFileRes.FSId))
	if err != nil {
		return fmt.Errorf("获取文件失败: %w", err)
	}
	fmt.Printf("文件名称: %s\n", file_res.GetFilename())
	fmt.Printf("文件md5: %s\n", file_res.MD5)
	return nil
}

// calculateBlockList 计算文件的MD5分块列表
func calculateBlockList(file *os.File, fileSize int64) ([]string, error) {
	blockList := make([]string, 0)

	// 如果文件小于等于4MB，只需要计算一个MD5
	if fileSize <= ChunkSize {
		hasher := md5.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return nil, err
		}
		md5Str := hex.EncodeToString(hasher.Sum(nil))
		blockList = append(blockList, md5Str)
	} else {
		// 如果文件大于4MB，需要按照4MB大小分片计算MD5
		chunkCount := int(fileSize / ChunkSize)
		if fileSize%ChunkSize != 0 {
			chunkCount++
		}

		for i := 0; i < chunkCount; i++ {
			offset, err := file.Seek(int64(i)*ChunkSize, 0)
			if err != nil {
				return nil, err
			}
			if offset != int64(i)*ChunkSize {
				return nil, fmt.Errorf("failed to seek to expected position: got %d, want %d", offset, int64(i)*ChunkSize)
			}

			hasher := md5.New()
			r := io.LimitReader(file, ChunkSize)
			if _, err := io.Copy(hasher, r); err != nil {
				return nil, err
			}

			md5Str := hex.EncodeToString(hasher.Sum(nil))
			blockList = append(blockList, md5Str)
		}
	}

	// 重置文件指针
	offset, err := file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	if offset != 0 {
		return nil, fmt.Errorf("failed to reset file pointer: got offset %d", offset)
	}

	return blockList, nil
}
