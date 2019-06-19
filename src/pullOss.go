package worker

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"os"
	"strings"
)

// Given a list of files, we will download them
func Download(nameMap map[string]string, c *Config) (map[string]string) {

	errors := map[string]string{}
	sess := getSession(c)
	downloader := s3manager.NewDownloader(sess)

	for item,mtime := range nameMap {

		file, err := os.Create(item)
		if err != nil {
			errors[item] = fmt.Sprintf("Unable to open file %v", err)
			continue
		}

		// Download this tar.gz file
		_, err = downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(c.BucketName),
				Key:    aws.String(item),
			})

		if err != nil {
			errors[item] = fmt.Sprintf("Unable to download item %v", err)
			continue
		}

		// Decompress it to Local path
		decompressedPath,err := decompress(item,c.LocalPath)
		if err!=nil {
			errors[item] = fmt.Sprintf("Failed in decomressing,reasons:%v",err)
			continue
		}

		// Remove remove the tar.gz file
		if err := os.Remove(item) ; err != nil {
			errors[item] = fmt.Sprintf("Failed in removing tar.gz,reasons:%v",err)
			continue
		}

		// Eventually change its mtime, same as it is on server
		if err := setMTime(decompressedPath,mtime);err!=nil{
			errors[item] = err.Error()
			continue
		}
	}

	return errors
}

func decompress(tarFile, dest string) (filename string,err error) {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return "",err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return "",err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return "",err
			}
		}
		filename = dest + hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return "",err
		}
		io.Copy(file, tr)
	}
	return filename,nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}