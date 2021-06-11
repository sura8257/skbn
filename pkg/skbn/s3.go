package skbn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sura8257/skbn/pkg/utils"
)

// Get aws config
func awsConfig() (aws.Config, error) {
	region := "us-east-2"

	if rg := os.Getenv("AWS_REGION"); rg != "" {
		region = rg
	}

	log.Printf("Using AWS_REGION: %s", region)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return cfg, err
	}
	return cfg, err
}

/*
// Copy copies files from src to dst
func copyS3ToS3(src, dst string) error {

	fmt.Println("copyS3ToS3: "+ src + "to" + dst)

	srcPrefix, srcPath := utils.SplitInTwo(src, "://")
	dstPrefix, dstPath := utils.SplitInTwo(dst, "://")

	srcPathSplit := strings.Split(srcPath, "/")
	srcBucket, _ := initS3Variables(srcPathSplit)

	dstPathSplit := strings.Split(dstPath, "/")
	dstBucket, _ := initS3Variables(dstPathSplit)


	return nil
	}
*/

// Download a file from s3 bucket
func copyS3ToFile(src, dst string, parallel int, bufferSize int64) error {

	log.Printf("Downloading file: %s ", src)
	_, srcPath := utils.SplitInTwo(src, "://")

	srcPathSplit := strings.Split(srcPath, "/")

	if err := validateS3Path(srcPathSplit); err != nil {
		return err
	}
	srcBucket, s3Path := initS3Variables(srcPathSplit)

	log.Printf("Bucket: %s File: %s", srcBucket, s3Path)

	cfg, err := awsConfig()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	log.Printf("Concurrency: %d PartSize: %d", parallel, uint64(bufferSize))

	client := s3.NewFromConfig(cfg)

	downloader := manager.NewDownloader(client, func(d *manager.Downloader) {
		d.Concurrency = parallel
		d.PartSize = bufferSize
	})

	_, err = downloader.Download(context.TODO(), f, &s3.GetObjectInput{
		Bucket: aws.String(srcBucket),
		Key:    aws.String(s3Path),
	})

	if err != nil {
		return err
	}

	err = utils.CheckFileStat(dst)
	if err != nil {
		return err
	}

	log.Printf("Successfully downloaded file: %s", dst)

	return nil
}

// Upload a file to s3 bucket
func copyFileToS3(src, dst string, parallel int, bufferSize int64) error {

	log.Printf("Uploading file: %s ", src)

	err := utils.CheckFileStat(src)
	if err != nil {
		return err
	}

	_, dstPath := utils.SplitInTwo(dst, "://")

	dstPathSplit := strings.Split(dstPath, "/")
	dstBucket, s3Path := initS3Variables(dstPathSplit)

	log.Printf("Bucket: %s File: %s", dstBucket, s3Path)

	cfg, err := awsConfig()
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	file, err := os.Open(src)

	if err != nil {
		fmt.Println("Unable to open file " + src)
		return err
	}

	defer file.Close()

	log.Printf("Concurrency: %d PartSize: %d", parallel, uint64(bufferSize))

	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.Concurrency = parallel
		u.PartSize = bufferSize
		u.LeavePartsOnError = false
	})

	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(dstBucket),
		Key:    aws.String(s3Path),
		Body:   file,
	})

	if err != nil {
		return err
	}

	log.Printf("Successfully uploaded")

	return nil
}

// Check S3 path
func validateS3Path(pathSplit []string) error {
	if len(pathSplit) >= 1 {
		return nil
	}
	return errors.New("illegal s3 path")
}

// Parse S3Uri
func initS3Variables(split []string) (string, string) {
	bucket := split[0]
	path := filepath.Join(split[1:]...)

	return bucket, path
}
