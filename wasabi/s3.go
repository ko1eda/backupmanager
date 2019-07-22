package wasabi

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// S3Service represents a wrapped s3 object
type S3Service struct {
	s3 s3iface.S3API
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

		return name, nil
	}

	return strings.TrimPrefix(*res.Location, "/"), nil
}
