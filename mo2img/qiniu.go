package mo2img

import (
	"os"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

var accessKey = os.Getenv("qiniuak")
var secretKey = os.Getenv("qiniusk")
var mac = qbox.NewMac(accessKey, secretKey)

// GenerateUploadToken generate qiniu upload token
func GenerateUploadToken(saveKey string) (token string) {
	var putPolicy = storage.PutPolicy{
		Scope:        "mo2",
		Expires:      7200,
		FsizeLimit:   1024 * 1024 * 20,
		SaveKey:      saveKey,
		ForceSaveKey: true,
	}
	return putPolicy.UploadToken(mac)
}
