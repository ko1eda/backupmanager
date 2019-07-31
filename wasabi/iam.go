package wasabi

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

// IAMService represents an IAM user
type IAMService struct {
	iam iamiface.IAMAPI
}

// NewIAMService returns a newly configured s3 service
func NewIAMService(sess *session.Session, configs ...*aws.Config) *IAMService {
	return &IAMService{
		iam: iam.New(sess, configs...),
	}
}

// CreateUser creates a new IAM user if it does not already exist for the given name
// All users created through this application will have a prefix of /ops/
func (i *IAMService) CreateUser(name string) (*iam.User, error) {
	if res, err := i.iam.GetUser(&iam.GetUserInput{UserName: &name}); err == nil {
		return res.User, nil
	}

	res, err := i.iam.CreateUser(&iam.CreateUserInput{UserName: aws.String(name)})

	if err != nil {
		// if there is any error other than entity already exists we have a problem
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != iam.ErrCodeEntityAlreadyExistsException {
				return nil, aerr
			}
		}
	}

	return res.User, nil
}

// CreateAccessKeyForUser creates an access key for the given user name
func (i *IAMService) CreateAccessKeyForUser(u *iam.User) (*iam.AccessKey, error) {
	result, err := i.iam.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: u.UserName,
	})

	if err != nil {
		return nil, err
	}

	return result.AccessKey, nil
}

// CreateLimitedAccessBucketPolicy creates a new policy with permissions scoped to a single bucket
// First we do a check to ensure that the policy exists, and if it does we return that.
// All policies created through this application have a path prefix of /opts/
func (i *IAMService) CreateLimitedAccessBucketPolicy(bucket string) (*iam.Policy, error) {
	prefix := "/" + bucket + "/"

	// If the policy already exists return that
	if res, err := i.iam.ListPolicies(&iam.ListPoliciesInput{PathPrefix: aws.String(prefix)}); err == nil {
		for _, pol := range res.Policies {
			if *pol.PolicyName == bucket+"-limited-access-policy" {
				return pol, nil
			}
		}
	}

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
		return nil, err
	}

	res, err := i.iam.CreatePolicy(&iam.CreatePolicyInput{
		PolicyDocument: aws.String(string(json)),
		PolicyName:     aws.String(bucket + "-limited-access-policy"),
		Path:           aws.String(prefix),
	})

	if err != nil {
		return nil, err
	}

	return res.Policy, nil
}

// AttachPolicyToUser attaches a policy to a given iam user
func (i *IAMService) AttachPolicyToUser(p *iam.Policy, u *iam.User) error {
	_, err := i.iam.AttachUserPolicy(&iam.AttachUserPolicyInput{
		PolicyArn: p.Arn,
		UserName:  u.UserName,
	})

	if err != nil {
		return err
	}

	return nil
}
