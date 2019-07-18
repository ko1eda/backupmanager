package main

import (
	"os"
	"os/signal"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	"github.com/ko1eda/backupmanager/http"
	"github.com/ko1eda/backupmanager/log"
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
	logger, closefn := log.New()
	defer closefn()

	// process env variables
	if err := godotenv.Load(); err != nil {
		logger.Log("Error loading .env file")
		os.Exit(1)
	}

	// create base s3 config
	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("WASABI_ACCESS_KEY_ID"),
			os.Getenv("WASABI_SECRET_ACCESS_KEY"),
			"",
		),
		Region: aws.String(os.Getenv("WASABI_REGION")),
	}

	s := session.New(config)

	s3Client := wasabi.NewS3Service(s, &aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String(os.Getenv("WASABI_S3_ENDPOINT")),
	})

	iamClient := wasabi.NewIAMService(s, &aws.Config{
		Endpoint: aws.String(os.Getenv("WASABI_IAM_ENDPOINT")),
	})

	srvr := http.NewServer("8080", s3Client, iamClient, logger)
	srvr.Open()
	defer srvr.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// log.Println("Listening on port 8080...")

	// http.Handle("/", bhttp.NewWasabiHandler(s3Client, iamClient))

	// http.ListenAndServe(":8080", nil)
}
