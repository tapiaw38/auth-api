package aws

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/private/protocol/rest"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Config is the aws s3 configuration
type S3Config struct {
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSBucket          string
}

// S3Client is the aws s3 client
type S3Client struct {
	Sess   *session.Session
	config *S3Config
}

// NewSession creates a new aws session
func NewSession(config *S3Config) *S3Client {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AWSRegion),
		Credentials: credentials.NewStaticCredentials(
			config.AWSAccessKeyID,
			config.AWSSecretAccessKey,
			""), // token can be left blank for now
	})
	if err != nil {
		log.Println("Failed to create new session", err)
	}

	return &S3Client{
		Sess:   sess,
		config: config,
	}
}

// UploadFileToS3 saves a file to aws bucket and returns the url to // the file and an error if there's any
func (t *S3Client) UploadFileToS3(file multipart.File, fileHeader *multipart.FileHeader, id string) (string, error) {
	// get the file size and read
	// the file content into a buffer
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	// create a unique file name for the file
	tempFileName := "pictures/" + strings.Split(fileHeader.Filename, ".")[0] + id + filepath.Ext(fileHeader.Filename)

	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	_, err := s3.New(t.Sess).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(t.config.AWSBucket),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		return "", err
	}

	return tempFileName, err
}

// GenerateUrl generates a url to the file in s3
func (t *S3Client) GenerateUrl(keyName string) string {
	req, _ := s3.New(t.Sess).GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(t.config.AWSBucket),
		Key:    aws.String(keyName),
	})
	rest.Build(req)
	urlStr := req.HTTPRequest.URL.String()

	return urlStr
}
