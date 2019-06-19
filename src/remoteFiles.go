package worker

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"strconv"
)

type Config struct {

	EndPoint   string
	Region     string
	AccessKey  string
	SecretKey  string
	BucketName string
	LocalPath  string

}

/*
    This function return a map, consists of
	map[fileName] = fileLastModifiedTime
*/

func RemoteFileMap(c *Config) (map[string]string,error) {

	fileMap := map[string]string{}
	sess := getSession(c)
	svc := s3.New(sess)

	resp,err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(c.BucketName),
	})
	if err != nil {
		return fileMap,err
	}

	for _,item := range resp.Contents {

		lastModified := item.LastModified.Unix()
		fileMap[*item.Key] = strconv.FormatInt(lastModified,10)
	}

	return fileMap,nil
}

func getSession(c *Config) *session.Session {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
			Region:      aws.String(c.Region),
			Endpoint:    aws.String(c.EndPoint),
		},
	}))

	return sess
}
