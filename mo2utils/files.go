package mo2utils

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storage"
	"io/ioutil"
	"log"
	"mo2/mo2img"
	"os"
	"path"
	"strings"
)

func CopyAllFiles(srcPath string, tgtPath string) {
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		//os.Mkdir(rootPath, 0755)
		fmt.Println("Directory created")
	} else {
		os.Mkdir(tgtPath, 0755)
		copyAllFiles(srcPath, tgtPath)
		fmt.Println("Directory already exists")
	}
	return

}
func copyAllFiles(curPath string, tgtPath string) {
	files, err := ioutil.ReadDir(curPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
		srcNewPath := path.Join(curPath, f.Name())
		tgtNewPath := path.Join(tgtPath, f.Name())
		if f.IsDir() {
			//mkdir
			os.Mkdir(tgtNewPath, 0755)
			copyAllFiles(srcNewPath, tgtNewPath)
		} else {
			// Part 1: open input file.
			inputFile, _ := os.Open(srcNewPath)

			// Part 2: call ReadAll to get contents of input file.
			data, _ := ioutil.ReadAll(inputFile)

			// Part 3: write data to copy file.
			ioutil.WriteFile(tgtNewPath, data, 0)
			fmt.Println("DONE")

		}
	}
}

func UploadAllFiles(srcPath string, uploadRootPath string) {
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		fmt.Println("Directory created")
	} else {
		uploadAllFiles(srcPath, srcPath, uploadRootPath)
	}
	return
}

func uploadAllFiles(curPath string, rootPath string, uploadRootPath string) {

	files, err := ioutil.ReadDir(curPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		srcNewPath := path.Join(curPath, f.Name())
		relativePath := strings.TrimPrefix(srcNewPath, rootPath)
		uploadRelativePath := path.Join(uploadRootPath, relativePath)

		if f.IsDir() {
			uploadAllFiles(srcNewPath, rootPath, uploadRootPath)
		} else {
			qiniuFile(srcNewPath, uploadRelativePath)
			//fmt.Println(srcNewPath)
			//fmt.Println(uploadRelativePath)
		}
	}
}

func UploadCDN() {
	if ginMode := os.Getenv("GIN_MODE"); ginMode == "release" {
		rootPath := "./dist"
		uploadRootPath := "dist"
		UploadAllFiles(rootPath, uploadRootPath)
	}
}

func qiniuFile(srcPath string, dstPath string) {
	localFile := srcPath
	key := dstPath
	upToken := mo2img.GenerateUploadToken(localFile)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	/*	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}*/
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret.Key, ret.Hash)
}
