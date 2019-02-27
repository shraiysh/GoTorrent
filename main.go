package main

import (
	"os"
	"fmt"
)

func main() {
	var errormsg = `Usage of ./main:
	--r
		  True if pause and resume feature is needed. False otherwise.
	--v
		  True if misc output is required. False otherwise.
	--files [path] [path] ...
		  List of Torrent Files`
	l := len(os.Args)
	if l == 1{
		fmt.Println(errormsg)
		return
	}
	files := make([]string, 1)
	var filesflag = false
	var resumeflag = false
	var verboseflag = false
	for i:=1;i<l;i++{
		arg := os.Args[i]
		if filesflag==true && arg[0]!='-' {
			files = append(files, arg)
		} else if filesflag==true{
			if arg=="--r" || arg=="-r" {
				resumeflag = true
			} else if arg=="--v" || arg=="-v"{
				verboseflag = true
			}
		}
		if arg == "--files"{
			filesflag=true
		}
	}
	


}
