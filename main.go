package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	"github.com/ko1eda/backupmanager/http"
	"github.com/ko1eda/backupmanager/wasabi"
)

func main() {
	// Parse incoming requsts
	// Log all requests and errors
	// Get hashed query parameter and decrpyt it

	// File to store logs
	f, err := os.OpenFile(makePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("FileOpenError: %v", err)
	}

	defer f.Close()

	mw := io.MultiWriter(os.Stderr, f)
	log.SetOutput(mw)

	// process env variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("ENVLoadError: ", err)
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

	srvr := http.NewServer(s3Client, iamClient)
	srvr.Open()
	defer srvr.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

// makePath makes a path to store log files on the current
// working directory
func makePath() string {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatalf("GetWorkingDirError: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(wd, "storage", "logs"), 0770); err != nil {
		log.Fatalf("CreateStorageDirError: %v", err)
	}

	return filepath.Join(wd, "storage", "logs", "application.log")
}
