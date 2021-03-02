package mo2utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"mo2/mo2img"
	"os"
	"path"
	"strings"

	"github.com/qiniu/go-sdk/v7/storage"
)

// UploadAllFiles upload all files in a dir to qiniu cdn
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

// IsEnvRelease check if GIN_MODE is release
func IsEnvRelease() (release bool) {
	return os.Getenv("GIN_MODE") == "release"
}

// UploadCDN upload frontend files to cdn
func UploadCDN() {
	if IsEnvRelease() {
		rootPath := "./dist"
		uploadRootPath := "dist"
		UploadAllFiles(rootPath, uploadRootPath)
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
