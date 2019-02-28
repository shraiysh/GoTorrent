package main

import (
	"fmt"
	"os"

	"github.com/concurrency-8/torrent"
)

func main() {
	var errormsg = `Usage of ./main:
	--help
		  Print this help message and exit.
	--r
		  True if pause and resume feature is needed. False otherwise.
	--v
		  True if misc output is required. False otherwise.
	--files [path] [path] ...
		  List of Torrent Files
	Sample input:
		  ./main --files File1 File2 File3 -v`
	l := len(os.Args)
	if l == 1 {
		fmt.Println(errormsg)
		return
	}
	if os.Args[1] == "--help" {
		fmt.Println(errormsg)
		return
	}
	files := make([]string, 0)
	filesflag := false
	resumeflag := false
	verboseflag := false
	for i := 1; i < l; i++ {
		arg := os.Args[i]
		if filesflag == true && arg[0] != '-' {
			files = append(files, arg)
		} else if filesflag == true {
			if arg == "--r" || arg == "-r" {
				resumeflag = true
			} else if arg == "--v" || arg == "-v" {
				verboseflag = true
			}
		}
		if arg == "--files" {
			filesflag = true
		}
	}
	ports := make([]int, len(files))
	//start peer ports from 20000. There's actually no restriction on the port numbers.
	//According to the specification, intially tracker UDP has ports from 6881-6889
	//Note that some of these ports may already be in use by the OS, in that case, Download fails.
	ports[0] = 20000
	for i, file := range files {
		if ports[i] == 0 {
			ports[i] = ports[i-1] + 1
		}
		go torrent.DownloadFromFile(file, ports[i])
	}
	if resumeflag {
		//TODO allow pause and resume capability here.
	}
	if verboseflag {
		//TODO print log to stdout here.
	}
}
