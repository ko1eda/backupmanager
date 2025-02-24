package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	"github.com/ko1eda/backupmanager/http"
	"github.com/ko1eda/backupmanager/smtp"
	"github.com/ko1eda/backupmanager/wasabi"
)

func main() {
	port := flag.String("p", "8080", "Set the port the server will run on")
	dir := flag.String("d", ".", "Set the directory where log files will be stored. Defaults to the current working directory")
	flag.Parse()

	// File to store logs
	f, err := os.OpenFile(makePath(*dir), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("FileOpenError: %v", err)
	}
	defer f.Close()

	// log to stderr and file
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

	// creates a new aws session (this is thread safe)
	s := session.New(config)

	// create a new s3 client
	s3Client := wasabi.NewS3Service(s, &aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String(os.Getenv("WASABI_S3_ENDPOINT")),
	})

	// create a new iam client
	iamClient := wasabi.NewIAMService(s, &aws.Config{
		Endpoint: aws.String(os.Getenv("WASABI_IAM_ENDPOINT")),
	})

	// creates a mailer with credentials
	mailer := smtp.NewMailer(os.Getenv("MAILER_ADDRESS"), os.Getenv("MAILER_FROM_ADDRESS"), smtp.WithCredentials(
		os.Getenv("MAILER_USERNAME"),
		os.Getenv("MAILER_PASSWORD"),
	))

	// validator is used to validate the secret keys on incoming requests
	validator := http.NewValidator(os.Getenv("REQUEST_SECRET_KEY"))

	// set up the server with all dependencies
	srvr := http.NewServer(http.WithAddress(*port))
	srvr.IAMService = iamClient
	srvr.S3Service = s3Client
	srvr.Mailer = mailer
	srvr.Validator = validator

	srvr.Open()
	defer srvr.Close()

	// we listen forever OR until we get a signal from the os to shut it down
	c := make(chan os.Signal, 1)
	// interupt syscall.SIGINT
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
}

// makePath makes a path to store log files
func makePath(dir string) string {
	dir = filepath.Clean(dir)

	if err := os.MkdirAll(filepath.Join(dir, "logs"), 0770); err != nil {
		log.Fatalf("CreateStorageDirError: %v", err)
	}

	return filepath.Join(dir, "logs", "application.log")
}
