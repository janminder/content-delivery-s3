package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type FileService interface {
	GetFile(bucket string, key string) (*os.File, error)
}

type fileService struct {
	conf *viper.Viper
}

func NewFileService(c *viper.Viper) FileService {
	return &fileService{
		conf: c,
	}
}

func (s *fileService) GetFile(bucket string, key string) (*os.File, error) {

	log.Debug("Get file from ", bucket, " with name ", key)

	// Configure to use MinIO Server
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(s.conf.GetString("s3.accessKey"), s.conf.GetString("s3.secret"), ""),
		Endpoint:         aws.String(s.conf.GetString("s3.protocol") + "://" + s.conf.GetString("s3.host") + ":" + s.conf.GetString("s3.port")),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess := session.New(s3Config)

	path := s.conf.GetString("s3.tmpDir")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	file, err := os.Create(path + key)
	defer s.closeFile(file)

	if err != nil {
		fmt.Println("Failed to create file", err)
	}

	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

	if err != nil {
		fmt.Println("Failed to download file", err)
		return nil, err
	} else {
		log.Info("Downloaded file ", file.Name() , " ", numBytes, " bytes")
		return file, nil
	}
}

func (s *fileService) closeFile(f *os.File) {
	log.Debug("closing ", f.Name())
	err := f.Close()
	if err != nil{
		log.Error("failed to close the tmp file ", f.Name(), " ", err)
	}
}

