package wasabi

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// IAMService represents an IAM user
type IAMService struct {
	iam *iam.IAM
}

// NewIAMService returns a newly configured s3 service
func NewIAMService(sess *session.Session, config *aws.Config) *IAMService {
	return &IAMService{
		iam.New(sess, config),
	}
}

// CreateUser creates a new IAM user
func (i *IAMService) CreateUser(name string) {
	_, err := i.iam.CreateUser(&iam.CreateUserInput{UserName: &name})

	if err != nil {
		log.Fatal(err)
	}
}
