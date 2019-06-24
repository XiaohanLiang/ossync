package worker

import (
	"fmt"
	"strings"
)

/*
	For each file remote, we would inspect if they exists locally

	IF they exists, and last modify date matches, then pass
	IF they do not exist, just download it from OSS
	IF they exists, and last modify date differs, download, by decompressing it, latest will rewrite the old

*/
func Diff(local,remote map[string]string)(map[string]string){

	ret := map[string]string{}
	trimmedLocal := trimFileName(local)
	trimmedRemote := trimFileName(remote)

	for file := range trimmedRemote {

		// If local file doesn't exists or last modify time doesn't match
		// We will download it later on
		if trimmedLocal[file] != trimmedRemote[file] {
			ret[file+".tar.gz"] = trimmedRemote[file]
		}

	}
	return ret
}

func validateFiles(fileName string) bool {

	if !strings.HasSuffix(fileName,".tar.gz"){
		fmt.Printf("[OSS-File] We've got a file name={%v}, it's not end with .tar.gz" +
			"while it is supposed to, we have skipped it",fileName)
		return false
	}

	return true
}

func trimFileName(l map[string]string)map[string]string{

	ret := map[string]string{}
	for k,v := range l {
		tk := strings.TrimSuffix(k,".tar.gz")
		ret[tk] = v
	}
	return ret
}