## Useful links
### Setting up the project
0. Make sure go runtime is installed on whatever machine you are using to build this project
1. create a .env file (it should look like the block below but make sure to add keys)

```
WASABI_ACCESS_KEY_ID=
WASABI_SECRET_ACCESS_KEY=
WASABI_S3_ENDPOINT=https://s3.wasabisys.com
WASABI_IAM_ENDPOINT=https://iam.wasabisys.com
WASABI_REGION="us-east-1"
```
3. run go build from the root directory of this project


### SDK Code examples
+ https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/common-examples.html
### Indiviual sdk docs
+ https://docs.aws.amazon.com/sdk-for-go/api/service/iam/#IAM.AttachUserPolicy 
### aws iface interfaces and testing 
+ https://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3iface/


## Building the logger
+ Path info
    + https://stackoverflow.com/questions/47261719/how-can-i-resolve-a-relative-path-to-absolute-path-in-golang
    + https://stackoverflow.com/questions/28448543/how-to-create-nested-directories-using-mkdir-in-golang
+ Open/Close file info
    + https://codereview.stackexchange.com/questions/59692/deferred-log-file-close
    + https://www.joeshaw.org/dont-defer-close-on-writable-files/


## Nested anonymous structs in go
```
policy := struct {
	Version   string
	Statement []struct{
		Effect   string
		Action   []string
		Resource []string
	}
}{
	Version:   "2012-10-17",
	Statement:[]struct{
		Effect   string
		Action   []string
		Resource []string
	}{
		{
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
		{
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

```