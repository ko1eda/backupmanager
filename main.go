package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

func main() {
	// Parse incoming requsts
	// Log all requests and errors
	// Get hashed query parameter and decrpyt it
	// create new IAM user
	// create new S3 bucket
	// create new IAM Policy
	// return credentials to user

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// load config from env
	bucket := aws.String("test1.securedatatransit.com")

	// key := aws.String("wasabi-testobject")

	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("WASABI_ACCESS_KEY_ID"),
			os.Getenv("WASABI_SECRET_ACCESS_KEY"),
			"",
		),
		Endpoint:         aws.String(os.Getenv("WASABI_ENDPOINT")),
		Region:           aws.String(os.Getenv("WASABI_REGION")),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession := session.New(config)

	s3Client := s3.New(newSession)

	output, err := s3Client.CreateBucket(&s3.CreateBucketInput{Bucket: bucket})

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

	fmt.Printf("%+v", output)

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
