package client

import (
	"context"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"go-server-base/global"
	"os"
	"path"
	"strconv"
	"time"

	osssdk "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type OssClient struct {
	scType    string
	bucketStr string
	client    osssdk.Client
	config    osssdk.Config
}

func NewOssClient() (*OssClient, error) {
	endpoint := global.CONF.Oss.Endpoint
	accessKey := global.CONF.Oss.AccessKeyId
	secretKey := global.CONF.Oss.AccessKeySecret
	bucketStr := global.CONF.Oss.BucketName
	scType := global.CONF.Oss.ScType
	if len(scType) == 0 {
		scType = "Standard"
	}
	// 填充 A K
	provider := credentials.NewStaticCredentialsProvider(accessKey, secretKey)
	cfg := osssdk.LoadDefaultConfig()
	cfg.WithCredentialsProvider(provider)
	cfg.WithRegion(global.CONF.Oss.Region)
	cfg.WithEndpoint(endpoint)
	cfg.WithConnectTimeout(20 * time.Second)
	cfg.WithReadWriteTimeout(60 * time.Second)
	//cfg.WithDisableSSL(true)
	cfg.WithLogLevel(osssdk.LogInfo)

	client := osssdk.NewClient(cfg)

	return &OssClient{scType: scType, bucketStr: bucketStr, client: *client, config: *cfg}, nil
}

func (o OssClient) ListBuckets() ([]interface{}, error) {
	response, err := o.client.ListBuckets(context.TODO(), &osssdk.ListBucketsRequest{})
	if err != nil {
		return nil, err
	}
	var result []interface{}
	for _, bucket := range response.Buckets {
		result = append(result, bucket.Name)
	}
	return result, err
}

func (o OssClient) Exist(path string) (bool, error) {
	exist, err := o.client.IsObjectExist(context.TODO(), o.bucketStr, path)
	if err != nil {
		return false, err
	}
	return exist, nil
}

//func (o OssClient) Size(path string) (int64, error) {
//	bucket, err := o.client.Bucket(o.bucketStr)
//	if err != nil {
//		return 0, err
//	}
//	lor, err := bucket.ListObjectsV2(osssdk.Prefix(path))
//	if err != nil {
//		return 0, err
//	}
//	if len(lor.Objects) == 0 {
//		return 0, fmt.Errorf("no such file %s", path)
//	}
//	return lor.Objects[0].Size, nil
//}

func (o OssClient) Delete(path string) (bool, error) {
	_, err := o.client.DeleteObject(context.TODO(), &osssdk.DeleteObjectRequest{
		Bucket: osssdk.Ptr(o.bucketStr),
		Key:    osssdk.Ptr(path),
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (o OssClient) Upload(src, target string) (bool, error) {

	// 打开文件
	file, err := os.Open(src)
	if err != nil {
		return false, fmt.Errorf("无法打开文件: %v", err)
	}

	defer file.Close()

	result, err := o.client.PutObject(context.TODO(), &osssdk.PutObjectRequest{
		Bucket:       osssdk.Ptr(o.bucketStr),
		Key:          osssdk.Ptr(target),
		Body:         file,
		StorageClass: osssdk.StorageClassType(o.scType),
	})
	if err != nil {
		return false, err
	}
	fmt.Sprintf("https://%s.%s/%s", o.bucketStr, o.config.Endpoint, target)
	global.LOG.Info(result.ResultCommon)
	return true, nil
}

func (o OssClient) Download(src, target string) (bool, error) {

	d := o.client.NewDownloader()
	_, err := d.DownloadFile(context.TODO(),
		&osssdk.GetObjectRequest{
			Bucket: osssdk.Ptr(o.bucketStr),
			Key:    osssdk.Ptr(src),
		},
		target,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (o OssClient) GetTimeUploadUrl() string {
	year := strconv.Itoa(time.Now().Year())
	month := strconv.Itoa(int(time.Now().Month()))
	day := strconv.Itoa(time.Now().Day())
	return path.Join(o.bucketStr, year, month, day)
}

//func (o *OssClient) ListObjects(prefix string) ([]string, error) {
//	bucket, err := o.client.Bucket(o.bucketStr)
//	if err != nil {
//		return nil, err
//	}
//	lor, err := bucket.ListObjectsV2(osssdk.Prefix(prefix))
//	if err != nil {
//		return nil, err
//	}
//	var result []string
//	for _, obj := range lor.Objects {
//		result = append(result, obj.Key)
//	}
//	return result, nil
//}
