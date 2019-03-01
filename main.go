package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/concurrency-8/args"
	"github.com/concurrency-8/torrent"
	"github.com/sethgrid/multibar"
)

func main() {
	var wait sync.WaitGroup
	downloadpath := ""
	var errormsg = `Usage of ./main:
	--help
		  Print this help message and exit.
	--download -d
		  Specify the download path for downloading the files.
	--rescap -rc
		  True if pause and resume feature is needed. False otherwise.
	--resume -r
		  True to resume partially downloaded files.
	--verbose -v
		  True if misc output is required. False otherwise.
	--files [path] [path] ...
		  List of Torrent Files
	Sample input:
		  ./concurrency-8 --files File1 File2 File3 -v -d ../../`
	l := len(os.Args)
	if l == 1 || os.Args[1] == "--help" {
		fmt.Println(errormsg)
		return
	}
	files := make([]string, 0)
	filesflag := false
	resumeflag := false
	rcflag := false
	verboseflag := false
	for i := 1; i < l; i++ {
		arg := os.Args[i]
		if filesflag == true && arg[0] != '-' {
			files = append(files, arg)
		} else {
			filesflag = false
			if arg == "--resume" || arg == "-r" {
				resumeflag = true
			} else if arg == "--verbose" || arg == "-v" {
				verboseflag = true
			} else if arg == "--rescap" || arg == "-rc" {
				rcflag = true
			} else if arg == "--download" || arg == "-d" {
				downloadpath = os.Args[i+1]
			}
		}
		if arg == "--files" {
			filesflag = true
		}
	}
	//Parsing CLI arguments complete.
	args.ARGS.FilePath = files
	args.ARGS.Resume = resumeflag
	args.ARGS.Verbose = verboseflag
	args.ARGS.DownloadPath = downloadpath
	args.ARGS.ResumeCapability = rcflag
	fmt.Printf("%v\n", args.ARGS)
	wait.Add(len(files))
	ports := make([]int, len(files))
	//start peer ports from 20000. There's actually no restriction on the port numbers.
	//According to the specification, initially tracker UDP has ports from 6881-6889
	//Note that some of these ports may already be in use by the OS, in that case, Download fails.

	progressbars, _ := multibar.New()

	ports[0] = 20000
	for i, file := range files {
		if ports[i] == 0 {
			ports[i] = ports[i-1] + 1
		}
		bar := progressbars.MakeBar(100, file)
		go func(file string, port int) {
			torrent.DownloadFromFile(file, port, &bar)
			defer wait.Done()
		}(file, ports[i])

	}

	go progressbars.Listen()

	wait.Wait()
}
