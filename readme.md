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

REQUEST_SECRET_KEY=""

MAILER_ADDRESS=""
MAILER_USERNAME=""
MAILER_PASSWORD=""
MAILER_FROM_ADDRESS=""

```
3. run go build from the root directory of this project
4. run the executable created from the go build command


### SDK Code examples
+ https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/common-examples.html
### Indiviual sdk docs
+ https://docs.aws.amazon.com/sdk-for-go/api/service/iam/#IAM.AttachUserPolicy 
### aws iface interfaces and testing 
+ https://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3iface/