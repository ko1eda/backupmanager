package wasabi

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Service represents a wrapped s3 object
type S3Service struct {
	s3 *s3.S3
}

// NewS3Service returns a newly configured s3 service
func NewS3Service(sess *session.Session, configs ...*aws.Config) *S3Service {
	return &S3Service{
		s3.New(sess, configs...),
	}
}

// CreateBucket returns an newly created S3 bucket
func (s3s *S3Service) CreateBucket(name string) (bucketname string, err error) {
	res, err := s3s.s3.CreateBucket(&s3.CreateBucketInput{Bucket: &name})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() != s3.ErrCodeBucketAlreadyExists {
			return "", err
		}

		log.Println("WARNING_BucketAlreadyExists")

		return name, nil
	}

	return strings.TrimPrefix(*res.Location, "/"), nil
}
