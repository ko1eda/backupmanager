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

// CreateAccessKeyForUser creates an access key for the given user name
func (i *IAMService) CreateAccessKeyForUser(u *iam.User) *iam.AccessKey {
	result, err := i.iam.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: u.UserName,
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Println("CreateAccessKeyError", aerr.Error())
		}
	}

	return result.AccessKey
}

// CreateLimitedAccessBucketPolicy creates a new policy with permissions scoped to a single bucket
func (i *IAMService) CreateLimitedAccessBucketPolicy(bucket string) *iam.Policy {
	type statemententry struct {
		Effect   string
		Action   []string
		Resource []string
	}

	type policydocument struct {
		Version   string
		Statement []statemententry
	}

	policy := policydocument{
		Version: "2012-10-17",
		Statement: []statemententry{
			statemententry{
				Effect: "Allow",
				Action: []string{
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

	res, err := i.iam.CreatePolicy(&iam.CreatePolicyInput{
		PolicyDocument: aws.String(string(json)),
		PolicyName:     aws.String(bucket + "-limited-access-policy"),
	})

	if err != nil {
		log.Fatal("Policy creation error", err)
	}

	return res.Policy
}

// AttachPolicyToUser attaches a policy to a given iam user
func (i *IAMService) AttachPolicyToUser(p *iam.Policy, u *iam.User) {
	_, err := i.iam.AttachUserPolicy(&iam.AttachUserPolicyInput{
		PolicyArn: p.Arn,
		UserName:  u.UserName,
	})

	if err != nil {
		log.Fatal("AttachPolicyError", err)
	}
}

// policy := struct {
// 	Version   string
// 	Statement []struct{
// 		Effect   string
// 		Action   []string
// 		Resource []string
// 	}
// }{
// 	Version:   "2012-10-17",
// 	Statement:[]struct{
// 		Effect   string
// 		Action   []string
// 		Resource []string
// 	}{
// 		{
// 			Effect: "Allow",
// 			Action: []string{
// 				"s3:DeleteObject",
// 				"s3:GetObject",
// 				"s3:HeadBucket",
// 				"s3:ListObjects",
// 				"s3:PutObject",
// 				"s3:ListBucket",
// 				"s3:GetBucketLocation",
// 				"s3:ListBucketMultipartUploads",
// 			},
// 			Resource: []string{
// 				"arn:aws:s3:::" + bucket,
// 				"arn:aws:s3:::" + bucket + "/*",
// 			},
// 		},
// 		{
// 			Effect: "Allow",
// 			Action: []string{
// 				"s3:*",
// 			},
// 			Resource: []string{
// 				"arn:aws:s3:::" + bucket,
// 				"arn:aws:s3:::" + bucket + "/*",
// 			},
// 		},
// 	},
// }
