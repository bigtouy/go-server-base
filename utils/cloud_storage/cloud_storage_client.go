package cloud_storage

import (
	"go-server-base/constant"
	"go-server-base/utils/cloud_storage/client"
)

type CloudStorageClient interface {
	ListBuckets() ([]interface{}, error)
	//ListObjects(prefix string) ([]string, error)
	Exist(path string) (bool, error)
	Delete(path string) (bool, error)
	Upload(src, target string) (bool, error)
	Download(src, target string) (bool, error)
	GetTimeUploadUrl() string
	//Size(path string) (int64, error)
}

func NewCloudStorageClient(backupType string) (CloudStorageClient, error) {
	switch backupType {
	//case constant.Local:
	//	return client.NewLocalClient(vars)
	case constant.S3:
		return client.NewS3Client()
	case constant.OSS:
		return client.NewOssClient()
	default:
		return nil, constant.ErrNotSupportType
	}
}
