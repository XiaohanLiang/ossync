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
	"path"
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

func decompress(srcFilePath string, destDirPath string) (ret string,err error){

	fmt.Println("UnTarGzing " + srcFilePath + "...")
	os.Mkdir(destDirPath, os.ModePerm)

	fr, err := os.Open(srcFilePath)
	if err != nil {
		return "",fmt.Errorf("Failed in opening src file, reasons:%v",err)
	}
	defer fr.Close()

	// 使用Gzip工具打开压缩文件
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return "",fmt.Errorf("Failed in opening with Gzip tool, reasons:%v",err)
	}

	// 使用Tar工具打开包
	tr := tar.NewReader(gr)

	folderName := true

	for {

		hdr, err := tr.Next()
		if folderName {
			ret = destDirPath + "/" + hdr.Name
			folderName = false
		}
		// EOF出现说明文件已经读完
		if err == io.EOF {
			break
		}

		fmt.Println("解压文件 -> file..." + hdr.Name)

		// 被装在包里的除了实际的文件还有文件夹, 但是文件夹不在解压范围内
		// 所以我们需要确保, 必须是文件,我们再开始解压, 写入
		if hdr.Typeflag != tar.TypeDir {

			os.MkdirAll(destDirPath+"/"+path.Dir(hdr.Name), os.ModePerm)
			fw, err := os.Create(destDirPath + "/" + hdr.Name)

			if err != nil {
				return "",fmt.Errorf("Failed in creating path,reasons:%v",err)
			}

			_, err = io.Copy(fw, tr)
			if err != nil {
				return "", fmt.Errorf("Failed in copying files,reasons:%v",err)
			}

			// 创建文件必须关闭,否则分分钟 too many open files
			if e := fw.Close(); e != nil {
				fmt.Printf("Failed in closing created files,reasons:%v",e.Error())
			}
		}
	}
	fmt.Println("Well done!")
	return
}