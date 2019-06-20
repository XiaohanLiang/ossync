package main

import (
	"fmt"
	"github.com/ossync/src"
	"github.com/urfave/cli"
	"os"
	"strings"
)

var (
	config  *worker.Config
)

func main() {

	app := cli.NewApp()
	app.Name = "ossync"
	app.Usage = "Helps you synchronizing files from OSS"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "e,endpoint",
			Usage: "Endpoint of your oss bucket",
		},
		cli.StringFlag{
			Name:  "r,region",
			Usage: "Region of your oss bucket",
		},
		cli.StringFlag{
			Name:  "a,accesskey",
			Usage: "AccessKey of your account",
		},
		cli.StringFlag{
			Name:  "s,secretkey",
			Usage: "SecretKey of your account",
		},
		cli.StringFlag{
			Name:  "b,bucket",
			Usage: "Name of your bucket",
		},
		cli.StringFlag{
			Name:  "p,path",
			Usage: "Local path, place you would like to deploy files",
		},
	}
	app.Run(os.Args)
}

func run(c *cli.Context) error {

	if c.NumFlags() < 6 {
		cli.ShowAppHelp(c)
		return cli.NewExitError("[Parameter Missing] Please set them all", -1)
	}

	path := c.String("path")
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	config = &worker.Config{
		EndPoint:   c.String("endpoint"),
		Region:     c.String("region"),
		AccessKey:  c.String("accesskey"),
		SecretKey:  c.String("secretkey"),
		BucketName: c.String("bucket"),
		LocalPath:  path,
	}

	remoteFileMap, err := worker.RemoteFileMap(config)
	if err != nil {
		return cli.NewExitError("[RemoteError] Cannot read remote OSS:"+err.Error(), -1)
	}

	fmt.Printf("[*RemoteFileMap]: %v \n", remoteFileMap)

	localFileMap, err := worker.LocalFileMap(config)
	if err != nil {
		return cli.NewExitError("[LocalError] Cannot read local files:"+err.Error(), -1)
	}

	fmt.Printf("[*LocalFileMap]: %v \n", localFileMap)

	fileList := worker.Diff(localFileMap, remoteFileMap)
	fmt.Printf("[*DownloadList] %v \n", fileList)

	errors := worker.Download(fileList, config)
	if len(errors) > 0 {
		printerror("[DownloadError] Some files cannot be downloaded from OSS -> \n")
		for k, v := range errors {
			fmt.Printf("\t\t FileName: %v, Fail Reason: %v \n", k, v)
		}
	} else {
		printline("Download complete :) \n")
	}

	return nil
}

func printline(s string) {
	green := string([]byte{27, 91, 51, 50, 109})
	reset := string([]byte{27, 91, 48, 109})
	fmt.Printf("%v%v%v", green,s,reset)
}
func printerror(s string) {
	reset := string([]byte{27, 91, 48, 109})
	red := string([]byte{27, 91, 51, 49, 109})
	fmt.Printf("%v%v%v",red,s,reset)
}
