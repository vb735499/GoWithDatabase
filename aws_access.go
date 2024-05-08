package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// - PUT, COPY, POST, or LIST	| 0.005  USD	/ every 1000+ requests.
// - GET, SELECT				| 0.0004 USD	/ every 1000+ requests.

type BucketBasics struct {
	S3Client *s3.Client
}

// Upload(Create)
func (basics BucketBasics) _UploadFile(bucketName string, objectKey string, filePath string) error {
	file, err := os.Open(filePath)

	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", filePath, err)
	} else {
		defer file.Close()
		_, err = basics.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
				filePath, bucketName, objectKey, err)
		}
	}
	return err
}

// Read
// ListObjects lists the objects in a bucket.
func (basics BucketBasics) _ListObjects(bucketName string) ([]types.Object, error) {
	result, err := basics.S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	var contents []types.Object
	if err != nil {
		log.Printf("Couldn't list objects in bucket %v. Here's why: %v\n", bucketName, err)
	} else {
		contents = result.Contents
	}
	return contents, err
}

// Download(fetch)
func (basics BucketBasics) _DownloadFile(bucketName string, objectKey string, filePath string) error {
	result, err := basics.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", filePath, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}

// Delete(Delete)
func (basics BucketBasics) _DeleteObjects(bucketName string, objectKeys []string) error {
	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	output, err := basics.S3Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{Objects: objectIds},
	})
	if err != nil {
		log.Printf("Couldn't delete objects from bucket %v. Here's why: %v\n", bucketName, err)
	} else {
		log.Printf("Deleted %v objects.\n", len(output.Deleted))
	}
	return err
}

/*
Describe:
  - Upload file to the AWS S3 bucket.

parameters:
@param bucketName: AWS S3 bucket name.
@param memberId: User account id.
@param filePath: The upload file path.
*/
func (basics BucketBasics) UploadFile(bucketName string, memberId string, filePath string) {
	filePaths := strings.Split(filePath, "/")
	fileName := filePaths[len(filePaths)-1]
	objectKey := "imgs/" + memberId + "/" + fileName
	log.Println("FileName" + fileName)
	basics._UploadFile(bucketName, objectKey, filePath)
}

/*
Describe:
  - Qeury All of files from AWS S3 bucket.

parameters:
@param bucketName: AWS S3 bucket name.

return:
@param _files: it's key-pair datatype. _files[username][0] = image
*/
func (basics BucketBasics) QueryAll(bucketName string) map[string][]string {
	output, _ := basics._ListObjects(bucketName)
	_files := make(map[string][]string)

	log.Println("first page results:")
	for _, object := range output {
		slice_s := strings.Split(aws.ToString(object.Key), "/")
		if slice_s[len(slice_s)-1] == "" {
			continue
		}
		prefix_uri := "https://s3-ap-southeast-2.amazonaws.com/" + bucketName + "/imgs/" + slice_s[1] + "/"
		_files[slice_s[1]] = append(_files[slice_s[1]], prefix_uri+slice_s[len(slice_s)-1])
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
	return _files
}

/*
Describe:
  - Download file from AWS S3 bucket to web serber.

parameters:
@param bucketName: AWS S3 bucket name.
@param memberId: User account id.
@param fileName: Download file to the {fileName} path.
*/
func (basics BucketBasics) DownloadFile(bucketName string, memberId string, fileName string) {
	objectKey := "imgs/" + memberId + "/" + fileName
	basics._DownloadFile(bucketName, objectKey, fileName)
}

/*
Describe:
  - Delete file on AWS S3 bucket.

parameters:
@param bucketName: AWS S3 bucket name.
@param memberId: User account id.
@param objectKeys: a list of filepath you want to delete.
*/
func (basics BucketBasics) DeleteObjects(bucketName string, memberId string, objectKeys []string) {
	memberFolder := "imgs/" + memberId + "/"
	var reDirectoryObjectsKeys []string
	for _, key := range objectKeys {
		reDirectoryObjectsKeys = append(reDirectoryObjectsKeys, memberFolder+key)
	}
	basics._DeleteObjects(bucketName, reDirectoryObjectsKeys)
}

/*
Describe:
  - Create a client used in access the AWS S3.

return:
@param BucketBasics: the client can access the AWS S3.
*/
func getClient() BucketBasics {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	return BucketBasics{
		S3Client: client,
	}
}
