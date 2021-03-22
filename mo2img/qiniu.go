package mo2img

import (
	"fmt"
	"os"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

var accessKey = os.Getenv("qiniuak")
var secretKey = os.Getenv("qiniusk")
var mac = qbox.NewMac(accessKey, secretKey)

const bucket = "mo2"

// GenerateUploadToken generate qiniu upload token
func GenerateUploadToken(saveKey string) (token string) {
	var putPolicy = storage.PutPolicy{
		Scope:        bucket,
		Expires:      7200,
		FsizeLimit:   1024 * 1024 * 20,
		SaveKey:      saveKey,
		ForceSaveKey: true,
	}
	return putPolicy.UploadToken(mac)
}

// GenerateOverwriteToken generate qiniu upload token
func GenerateOverwriteToken(saveKey string, overwriteKey string) (token string) {
	var putPolicy = storage.PutPolicy{
		Scope:        fmt.Sprintf("%s:%s", bucket, overwriteKey),
		Expires:      7200,
		FsizeLimit:   1024 * 1024 * 20,
		SaveKey:      saveKey,
		ForceSaveKey: true,
	}
	return putPolicy.UploadToken(mac)
}
