package wasabi

import (
	"fmt"

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

	// s3Client := wasabi.NewS3Service(newSession, &aws.Config{
	// 	S3ForcePathStyle: aws.Bool(true),
	// })

	// fmt.Printf("%+v", output)

	// Upload a new object "wasabi-testobject" with the string "Wasabi Hot storage"
	// _, err = s3Client.PutObject(&s3.PutObjectInput{
	// 	Body:   strings.NewReader("wasabi hot storage"),
	// 	Bucket: bucket,
	// 	Key:    key,
	// })
	// if err != nil {
	// 	fmt.Printf("Failed to upload object %s/%s, %s\n", *bucket, *key, err.Error())
	// 	return
	// }
	// fmt.Printf("Successfully uploaded key %s\n", *key)

	//Get Object
	// _, err = s3Client.GetObject(&s3.GetObjectInput{
	// 	Bucket: bucket,
	// 	Key:    key,
	// })
	// if err != nil {
	// 	fmt.Println("Failed to download file", err)
	// 	return
	// }
	// fmt.Printf("Successfully Downloaded key %s\n", *key)
}

// CreateBucket returns a newly configured s3 service
func (s3s *S3Service) CreateBucket(name string) {

	_, err := s3s.s3.CreateBucket(&s3.CreateBucketInput{Bucket: &name})

	// TODO switch this because wasabi doesn't seem to log bucket exists error
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}
}
