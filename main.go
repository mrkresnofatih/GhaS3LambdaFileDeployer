package main

import (
	"fmt"
	"os"
	"strings"

	aws "github.com/aws/aws-sdk-go/aws"
	sess "github.com/aws/aws-sdk-go/aws/session"
	lambda "github.com/aws/aws-sdk-go/service/lambda"
	s3manager "github.com/aws/aws-sdk-go/service/s3/s3manager"
	xid "github.com/rs/xid"
)

func main() {
	// Handling Inputs from ENV
	filePathInRepo := os.Getenv("FILE_PATH")
	awsKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	bucketAddress := os.Getenv("BUCKET_ADDRESS")
	fileName := os.Getenv("FILE_NAME")
	lambdaFunc := os.Getenv("LAMBDA_FUNC")

	// Check If All Params Are Available
	isAllParamsAvailable := getIsAllParamsAvailable(
		filePathInRepo,
		awsKey,
		awsSecret,
		awsRegion,
		bucketAddress,
		fileName,
		lambdaFunc,
	)
	if !isAllParamsAvailable {
		fmt.Println("[FAIL]: Input incomplete!")
		return
	}
	fmt.Println("[SCCS]: All Params Available!")

	// Create AWS Session
	sess, err := sess.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		fmt.Println("[FAIL]: Failed to start new S3 Manager Session: ", err)
		return
	}
	fmt.Println("[SCCS]: Created AWS Session!")

	// Create Uploader Instance
	uploader := s3manager.NewUploader(sess)
	f, er := os.Open(filePathInRepo)
	if er != nil {
		fmt.Println("[FAIL]: Failed to open file: ", er)
		return
	}
	fmt.Println("[SCCS]: Uploader Instance Created!")

	// Upload The File to S3
	versionedFilename := generateVersionedFilename(fileName)
	_, e := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketAddress),
		Key:    aws.String(versionedFilename),
		Body:   f,
	})
	if e != nil {
		fmt.Println("[FAIL]: Failed to upload file: ", e)
		return
	}
	fmt.Println("[SCCS]: File Uploaded to S3!")

	// Create Lambda Client Instance
	lsvc := lambda.New(sess)

	// Create UpdateFunctionCodeInput
	req := new(lambda.UpdateFunctionCodeInput)
	req.SetFunctionName(lambdaFunc)
	req.SetS3Bucket(bucketAddress)
	req.SetS3Key(versionedFilename)

	// Validate UpdateFunctionCodeInput
	req_err := req.Validate()
	if req_err != nil {
		fmt.Println("[FAIL]: UpdateFunctionCodeInput Validation Failed: ", req_err)
		return
	}
	fmt.Println("[SCCS]: UpdateFunctionCodeInput Validation Succeeded!")

	// Update Lambda Function
	_, errr := lsvc.UpdateFunctionCode(req)
	if errr != nil {
		fmt.Println("[FAIL] Failed to Update Lambda Function Code", errr)
		return
	}
	fmt.Println("[SCCS]: Target Lambda Function Updated!")
}

// Generate Filename with Versioning Suffix
func generateVersionedFilename(filename string) string {
	filenames := strings.Split(filename, ".")
	name := filenames[0]
	ext := filenames[1]
	version := xid.New().String()
	finalName := name + "-" + version + "." + ext
	return finalName
}

// Validate Inputs
func getIsAllParamsAvailable(
	input1 string,
	input2 string,
	input3 string,
	input4 string,
	input5 string,
	input6 string,
	input7 string,
) bool {
	isAllAvailable := (input1 != "")
	isAllAvailable = (isAllAvailable && (input2 != ""))
	isAllAvailable = (isAllAvailable && (input3 != ""))
	isAllAvailable = (isAllAvailable && (input4 != ""))
	isAllAvailable = (isAllAvailable && (input5 != ""))
	isAllAvailable = (isAllAvailable && (input6 != ""))
	isAllAvailable = (isAllAvailable && (input7 != ""))

	return isAllAvailable
}
