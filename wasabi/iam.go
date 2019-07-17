package wasabi

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// IAMService represents an IAM user
type IAMService struct {
	iam *iam.IAM
}

// NewIAMService returns a newly configured s3 service
func NewIAMService(sess *session.Session, configs ...*aws.Config) *IAMService {
	return &IAMService{
		iam.New(sess, configs...),
	}
}

// CreateUser creates a new IAM user
func (i *IAMService) CreateUser(name string) {
	res, err := i.iam.CreateUser(&iam.CreateUserInput{UserName: &name})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Println("Called from Iam service create user")
			log.Println(aerr.Code())
			log.Println(aerr.Error())
			log.Println(aerr.OrigErr())
			log.Fatal(aerr.Message())
		}
	}

	log.Printf("%+v", res)
}
