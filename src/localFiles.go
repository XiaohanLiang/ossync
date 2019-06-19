package worker

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func LocalFileMap(c *Config)(map[string]string,error){

	fileMap := map[string]string{}
	files, err := ioutil.ReadDir(c.LocalPath)
	if err != nil {
		return fileMap,err
	}

	for _, f := range files {

		fileName := f.Name()
		fileLastModify := f.ModTime().Unix()
		fileMap[fileName] = strconv.FormatInt(fileLastModify,10)
	}
	return fileMap,nil
}

func setMTime(file string,mtime string) error {

	i, _ := strconv.ParseInt(mtime, 10, 64)
	tm := time.Unix(i, 0)

	err := os.Chtimes(file,time.Now(), tm)
	if err != nil {
		return fmt.Errorf("Cannot modify its mtime,reasons:%v",err)
	}

	return nil
}
