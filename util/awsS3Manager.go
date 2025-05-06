package util

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var cfg aws.Config
var s3client *s3.Client

type BucketBasics struct {
	S3Client *s3.Client
}

func init() {
	region := viper.GetString("AWS_REGION")
	var err error
	cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))

	if err != nil {
		log.Error(err)
	}
	s3client = s3.NewFromConfig(cfg)
}

func ListObjects(bucket string) ([]types.Object, error) {
	input := &s3.ListObjectsInput{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int32(50),
	}

	result, err := s3client.ListObjects(context.TODO(), input)
	var s3objects []types.Object
	if err != nil {
		log.Errorf("Couldn't list objects in buckets for your account. Here's why: %v\n", err)
	} else {
		s3objects = result.Contents
	}
	return s3objects, err
}

func BucketExists(bucketName string) (bool, error) {
	_, err := s3client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Errorf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				log.Errorf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	} else {
		log.Infof("Bucket %v exists and you already own it.", bucketName)
	}

	return exists, err
}

func DownloadFile(bucketName string, objectKey string, fileName string) error {
	result, err := s3client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Errorf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Errorf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Errorf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}
