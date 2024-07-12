package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	bucket = "tendercall-db"
	region = "us-east-1" // Replace with your AWS region
)

// HandleFileUpload function declared

func HandleFileUpload(w http.ResponseWriter, r *http.Request) {

	// Check if the request method is POST

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieves the uploaded file named "file" from the HTTP request

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Determines the file name and path where the file will be stored locally

	fileName := header.Filename
	filePath := filepath.Join("./upload", fileName)
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file on server", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// io.Copy Copies the contents of the uploaded file to the newly created file on the server.
	// os.Create Creates a new file on the server's filesystem to store the uploaded content.

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to copy file", http.StatusInternalServerError)
		return
	}

	// UploadFileToS3 to upload the local file to AWS S3 bucket in the specified region.
	url, err := UploadFileToS3(bucket, region, filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload file to S3: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s uploaded successfully to bucket %s\n", fileName, bucket)
	fmt.Fprintf(w, "Public URL: %s\n", url)
}

// HandleGetURL function declaration

func HandleGetURL(w http.ResponseWriter, r *http.Request) {

	// Check if the request method is GET

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieves the value of the query parameter named "file" from the URL of the HTTP request.

	fileName := r.URL.Query().Get("file")

	if fileName == "" {
		http.Error(w, "Missing 'file' parameter", http.StatusBadRequest)
		return
	}

	// onstructs the public URL for accessing the file stored in AWS S3

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, fileName)
	fmt.Fprintf(w, "Public URL for %s: %s\n", fileName, url)
}

// UploadFileToS3 function declaration

func UploadFileToS3(bucket, region, filePath string) (string, error) {

	// To open the file specified by filePath using os.Open.
	// If the file does not exist, os.Open returns an error.

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q: %v", filePath, err)
	}
	defer file.Close()

	// Retrieves information about the file using file.Stat()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}

	/*
		Initialize AWS session with credentials and configuration
		aws.Config is a struct that holds the configuration for the AWS SDK for Go.
		session.NewSession initializes a new AWS session with the specified configuration..
		Region specifies the AWS region where AWS service requests are sent.
		Credentials allows specifying AWS credentials used to authenticate requests.
	*/

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			"AKIAYS2NVN4MBSHP33FF",                     // replace with your access key ID
			"aILySGhiQAB7SaFnqozcRZe1MhZ0zNODLof2Alr4", // replace with your secret access key
			""), // optional token, leave blank if not using
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	// s3.New(sess) creates a new S3 service client using the AWS session (sess) created in the previous step.

	svc := s3.New(sess)

	/*
		Upload file to s3
		fmt.Sprintf constructs a unique object key (imageKey) for the image file in the S3 bucket.
		svc.PutObject uploads the image (fileBytes) to the specified S3 bucket ("your-bucket-name") with the specified object key (imageKey).
		bytes.NewReader(fileBytes) sets the content of the object to be uploaded.
	*/

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(filepath.Base(filePath)),
		Body:          file,
		ContentLength: aws.Int64(fileInfo.Size()),
		ContentType:   aws.String("image/jpeg"), // adjust content type as needed
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// // Generate public URL
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, filepath.Base(filePath))

	return url, nil
}
