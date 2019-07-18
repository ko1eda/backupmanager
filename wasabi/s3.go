package wasabi

import (
	"log"

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

// CreateBucket returns a newly configured s3 service
func (s3s *S3Service) CreateBucket(name string) {
	res, err := s3s.s3.CreateBucket(&s3.CreateBucketInput{Bucket: &name})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == s3.ErrCodeBucketAlreadyExists {
			log.Println("BucketAlreadyExists", err)

			return
		}

		log.Println("CreateBucketError", err)
	}

	log.Println(*res.Location)

	// TODO switch this because wasabi doesn't seem to log bucket exists error
	// if err != nil {
	// 	if aerr, ok := err.(awserr.Error); ok {
	// 		switch aerr.Code() {
	// 		case s3.ErrCodeBucketAlreadyExists:
	// 			fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
	// 		case s3.ErrCodeBucketAlreadyOwnedByYou:
	// 			fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
	// 		default:
	// 			fmt.Println(aerr.Error())
	// 		}
	// 	} else {
	// 		fmt.Println(err.Error())
	// 	}
	// }
}
