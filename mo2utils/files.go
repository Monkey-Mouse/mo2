package mo2utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Monkey-Mouse/mo2/mo2img"

	"github.com/qiniu/go-sdk/v7/storage"
)

// ProcessAllFiles use handler function to process each file recursively under srcPath
func ProcessAllFiles(srcPath string, uploadRootPath string, processHandler func(parameter ...string)) {
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		fmt.Println("Directory created")
	} else {
		processAllFiles(srcPath, srcPath, uploadRootPath, processHandler)
	}
	return
}

func processAllFiles(curPath string, rootPath string, uploadRootPath string, processHandler func(parameter ...string)) {

	files, err := ioutil.ReadDir(curPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		srcNewPath := path.Join(curPath, f.Name())
		relativePath := strings.TrimPrefix(srcNewPath, rootPath)
		uploadRelativePath := path.Join(uploadRootPath, relativePath)
		if f.IsDir() {
			processAllFiles(srcNewPath, rootPath, uploadRootPath, processHandler)
		} else {
			processHandler(srcNewPath, rootPath, uploadRelativePath)
		}
	}
}

// IsEnvRelease check if GIN_MODE is release
func IsEnvRelease() (release bool) {
	return os.Getenv("GIN_MODE") == "release"
}

// UploadCDN upload frontend files to cdn
func UploadCDN() {
	if IsEnvRelease() {
		rootPath := "./dist"
		uploadRootPath := "dist"
		ProcessAllFiles(rootPath, uploadRootPath, func(parameter ...string) {
			qiniuFile(parameter[0], parameter[2])
		})
	}
}

func qiniuFile(srcPath string, dstPath string) {
	localFile := srcPath
	key := dstPath
	upToken := mo2img.GenerateOverwriteToken(localFile, dstPath)
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
