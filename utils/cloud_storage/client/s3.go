package client

import (
	"go-server-base/global"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Client struct {
	scType string
	bucket string
	Sess   session.Session
}

func NewS3Client() (*S3Client, error) {
	accessKey := global.CONF.S3.AccessKeyId
	secretKey := global.CONF.S3.AccessKeySecret
	endpoint := global.CONF.S3.Endpoint
	region := global.CONF.S3.Region
	bucket := global.CONF.S3.BucketName
	scType := global.CONF.S3.ScType
	if len(scType) == 0 {
		scType = "Standard"
	}
	mode := global.CONF.S3.Mode
	if len(mode) == 0 {
		mode = "virtual hosted"
	}
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(mode == "path"),
	})
	if err != nil {
		return nil, err
	}
	return &S3Client{scType: scType, bucket: bucket, Sess: *sess}, nil
}

func (s S3Client) ListBuckets() ([]interface{}, error) {
	var result []interface{}
	svc := s3.New(&s.Sess)
	res, err := svc.ListBuckets(nil)
	if err != nil {
		return nil, err
	}
	for _, b := range res.Buckets {
		result = append(result, b.Name)
	}
	return result, nil
}

func (s S3Client) Exist(path string) (bool, error) {
	svc := s3.New(&s.Sess)
	if _, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: &s.bucket,
		Key:    &path,
	}); err != nil {
		if aerr, ok := err.(awserr.RequestFailure); ok {
			if aerr.StatusCode() == 404 {
				return false, nil
			}
		} else {
			return false, aerr
		}
	}
	return true, nil
}

func (s *S3Client) Size(path string) (int64, error) {
	svc := s3.New(&s.Sess)
	file, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &path,
	})
	if err != nil {
		return 0, err
	}
	return *file.ContentLength, nil
}

func (s S3Client) Delete(path string) (bool, error) {
	svc := s3.New(&s.Sess)
	if _, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(path)}); err != nil {
		return false, err
	}
	if err := svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (s S3Client) Upload(src, target string) (bool, error) {
	fileInfo, err := os.Stat(src)
	if err != nil {
		return false, err
	}
	file, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer file.Close()

	uploader := s3manager.NewUploader(&s.Sess)
	if fileInfo.Size() > s3manager.MaxUploadParts*s3manager.DefaultUploadPartSize {
		uploader.PartSize = fileInfo.Size() / (s3manager.MaxUploadParts - 1)
	}
	if _, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:       aws.String(s.bucket),
		Key:          aws.String(target),
		Body:         file,
		StorageClass: &s.scType,
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (s S3Client) Download(src, target string) (bool, error) {
	if _, err := os.Stat(target); err != nil {
		if os.IsNotExist(err) {
			os.Remove(target)
		} else {
			return false, err
		}
	}
	file, err := os.Create(target)
	if err != nil {
		return false, err
	}
	defer file.Close()
	downloader := s3manager.NewDownloader(&s.Sess)
	if _, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(src),
	}); err != nil {
		os.Remove(target)
		return false, err
	}
	return true, nil
}

func (s *S3Client) ListObjects(prefix string) ([]string, error) {
	svc := s3.New(&s.Sess)
	var result []string
	outputs, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: &s.bucket,
		Prefix: &prefix,
	})
	if err != nil {
		return result, err
	}
	for _, item := range outputs.Contents {
		result = append(result, *item.Key)
	}
	return result, nil
}

func (s S3Client) GetTimeUploadUrl() string {
	year := strconv.Itoa(time.Now().Year())
	month := strconv.Itoa(int(time.Now().Month()))
	day := strconv.Itoa(time.Now().Day())
	return path.Join(s.bucket, year, month, day)
}
