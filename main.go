package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	bhttp "github.com/ko1eda/backupmanager/http"
	"github.com/ko1eda/backupmanager/wasabi"
)

func main() {
	// Parse incoming requsts
	// Log all requests and errors
	// Get hashed query parameter and decrpyt it
	// create new IAM user
	// create new S3 bucket
	// create new IAM Policy
	// return credentials to user

	// process env variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// load config from env
	// bucket := aws.String("test1.securedatatransit.com")
	// key := aws.String("wasabi-testobject")

	// create base s3 config
	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("WASABI_ACCESS_KEY_ID"),
			os.Getenv("WASABI_SECRET_ACCESS_KEY"),
			"",
		),
		Endpoint: aws.String(os.Getenv("WASABI_ENDPOINT")),
		Region:   aws.String(os.Getenv("WASABI_REGION")),
	}

	s := session.New(config)

	s3Client := wasabi.NewS3Service(s, &aws.Config{
		S3ForcePathStyle: aws.Bool(true),
	})

	log.Println("Listening on port 8080...")

	http.Handle("/", bhttp.NewWasabiHandler(s3Client))
	http.ListenAndServe(":8080", nil)
}
