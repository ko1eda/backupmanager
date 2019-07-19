## Useful links
### SDK Code examples
+ https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/common-examples.html
### Indiviual sdk docs
+ https://docs.aws.amazon.com/sdk-for-go/api/service/iam/#IAM.AttachUserPolicy 


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