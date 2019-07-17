package wasabi

import (
	"encoding/json"
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
		iam: iam.New(sess, configs...),
	}
}

// CreateUser creates a new IAM user
func (i *IAMService) CreateUser(name string) *iam.User {
	res, err := i.iam.CreateUser(&iam.CreateUserInput{UserName: &name})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Println("Called from Iam service create user")
			log.Println(aerr.Code())
			log.Println(aerr.Error())
			// log.Println(aerr.OrigErr())
			// log.Fatal(aerr.Message())
		}
	}

	return res.User
}

// CreateSingleUserBucketPolicy creates a new policy with permissions scoped to a single bucket
func (i *IAMService) CreateSingleUserBucketPolicy(bucket string) {
	policy := policydocument{
		Version: "2012-10-17",
		Statement: []statemententry{
			statemententry{
				Effect: "Allow",
				Action: []string{
					"s3:DeleteObject",
					"s3:GetObject",
					"s3:HeadBucket",
					"s3:ListObjects",
					"s3:PutObject",
					"s3:ListBucket",
					"s3:GetBucketLocation",
					"s3:ListBucketMultipartUploads",
				},
				Resource: []string{
					"arn:aws:s3:::" + bucket,
					"arn:aws:s3:::" + bucket + "/*",
				},
			},
			statemententry{
				Effect: "Allow",
				Action: []string{
					"s3:*",
				},
				Resource: []string{
					"arn:aws:s3:::" + bucket,
					"arn:aws:s3:::" + bucket + "/*",
				},
			},
		},
	}

	json, err := json.Marshal(policy)

	if err != nil {
		log.Fatal("unmarshall error", err)
	}

	_, err = i.iam.CreatePolicy(&iam.CreatePolicyInput{
		PolicyDocument: aws.String(string(json)),
		PolicyName:     aws.String(bucket + "-limited-access-policy"),
	})

	if err != nil {
		log.Fatal("Policy creation error", err)
	}
}

type policydocument struct {
	Version   string
	Statement []statemententry
}

type statemententry struct {
	Effect   string
	Action   []string
	Resource []string
}
