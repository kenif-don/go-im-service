package util

import (
	"bytes"
	"context"
	"go-im-service/src/configs/conf"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/configs/log"
	"time"

	"git.sr.ht/~jamesponddotco/bunnystorage-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Upload2S3(filename string, data []byte) (string, *utils.Error) {
	sess, e := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(conf.Conf.Aws.Id, conf.Conf.Aws.Secret, ""),
		Endpoint:         aws.String(conf.Conf.Aws.Endpoint),
		Region:           aws.String(conf.Conf.Aws.Region),
		S3ForcePathStyle: aws.Bool(true),
	})
	if e != nil {
		log.Debug(e)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	uploader := s3.New(sess)
	_, e = uploader.PutObjectWithContext(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(conf.Conf.Aws.Bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(data),
		ACL:    aws.String("public-read"),
	})
	if e != nil {
		log.Error(e)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	// 获取预览URL
	return "https://" + conf.Conf.Aws.Endpoint + "/" + conf.Conf.Aws.Bucket + "/" + filename, nil
}
func Upload2Bunny(filename string, data []byte) (string, *utils.Error) {
	cfg := &bunnystorage.Config{
		StorageZone: conf.Conf.Aws.Id,
		Key:         conf.Conf.Aws.Secret,
		Endpoint:    bunnystorage.EndpointSingapore,
		Timeout:     60 * 60 * time.Second,
	}
	client, e := bunnystorage.NewClient(cfg)
	if e != nil {
		log.Error(e)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	_, e = client.Upload(context.Background(), "", filename, "", bytes.NewReader(data))
	if e != nil {
		log.Error(e)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	// 获取预览URL
	return conf.Conf.Aws.Endpoint + "/" + filename, nil
}
