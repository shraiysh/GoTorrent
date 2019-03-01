# GoTorrent
 [![CircleCI](https://circleci.com/gh/IITH-SBJoshi/concurrency-8/tree/master.svg?style=svg&circle-token=88f8e60508e4f98f339d7b395c228c6f309c2564)](https://circleci.com/gh/IITH-SBJoshi/concurrency-8/tree/master) [![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/IITH-SBJoshi/concurrency-8/blob/master/LICENSE) [![GoDoc](https://godoc.org/github.com/narqo/go-badge?status.svg)](https://iith-sbjoshi.github.io/concurrency-8/pkg/github.com/concurrency-8)
 
Bit-torrent client implementation as part of CS2443 (Principles of Programming Language II ) by Prof. Saurabh Joshi .

* [Description](#description)
* [Documentation](#Documentation)
* [Setup](#Setup)
* [Usage](#Usage)
* [Guidelines for contribution](#Guidelines-for-contribution)
* [Resources and References](#Resources-and-References)

## Description
1. **Objective**
	- Get familair with writing concurrent programs.
	- Using software technologies like Continous Integration , Unit Testing , Documentation .
2. **Features**
	- Downloading multiple torrent files concurrently.
	- Fetching Peer lists from both HTTP and UDP Trackers.
	- Fetching pieces of blocks concurrently from Peers.
	- Enabling Resume capabilities on abrupt termination.
	- Generating detailed log files for debugging.
	- A command line interface for managing.
3. **Team**
	- Shraiysh Gupta (CS17BTECH11050)
	- Puneet Mangla (CS17BTECH11029)
	- Lingam Sai Ramana Reddy (Cs17BTECH11022)
	- Hitesh (MA17BTECH11004)

## Documentation
1. You can refer to [this](https://iith-sbjoshi.github.io/concurrency-8/pkg/github.com/concurrency-8) to see the documentation generated for the master branch.

## Setup
1. **Installing Golang**
	- Follow [this](https://golang.org/doc/install) link **OR**
	- Run ```sudo apt-get install golang```
	- Set the environment variables `GOPATH` and `GOBIN` as follows :
		- ```GOPATH="$HOME/go"```
		- ```GOBIN="$GOPATH/bin"```
		- ```PATH=$PATH:$GOBIN```
2. **Building**
	
	Get [dep](https://github.com/golang/dep) for installing the dependencies
	```
	$ cd $GOPATH/src/github.com # Come to the appropriate directory
	$ git clone https://github.com/IITH-SBJoshi/concurrency-8.git
	$ cd concurrency-8/
	$ dep ensure		# Get the dependencies
	$ ./build.sh    	# To check if all tests passes
	```
## Usage
1. **Downloading**
	- ```go run main.go --files File1 File2 File3 -v -d ../../```
2. **Flags**

| __Flag Name__ | __Description__ | __Default__ |
|-------------|------------|------------|
| ```--files [path] [path] ...``` |  List of Torrent Files | empty |
| ```--download -d```  | Specify the download path for downloading the files.| "" |
| ```--rescap -rc```  | True if pause and resume feature is needed. False otherwise. | false |
| ```--resume -r```  | True to resume partially downloaded files. | false |
| ```--help```  | Print this help message and exit. |- |
| ```--verbose -v```  | True if misc output is required. False otherwise. | false |

## Guidelines for contribution :
1. Take open issues and ask for assignment in comment section.
2. **Working on seperate branch**
	- Clone the repository : ```git clone https://github.com/IITH-SBJoshi/concurrency-8.git```
	- Create a issue specific branch in cloned repository : ```git checkout -b issue#<issue number>```
	- Run the code by following the steps above
	- You can now start working on your current branch
3. **Testing the changes**
	- Run the test cases if any: ```go test <test file>.go```
	- Check the linting (Install [golint](https://github.com/golang/lint), if not already installed): ```golint <file_name>```
	- Run `./build.sh` to check if the build passes

	Note: If running `goreportcard-cli -v` shows errors in the files that are in the `vendor/` directory, ignore those issues. The TRAVIS build will take care that those files are not checked.
4. **Commiting the changes**
	- Update ```.gitignore``` if there is any need .
	- To add changes in your working directory : ```git add .```
	- Commit your changes : ```git commit -m "<message>"```
	- Follow a simple commit message guideline eg . ``` Fix <issue_id> : <small description> Author@<your name>```
5. **Pushing the changes**
	- Get current master: `git fetch origin master`
	- Merge master with your branch: `git merge master`
	- Push your changes : ```git push origin <your branch name>:<your branch name>```
	- Make sure that ```Travis CI build``` is passed.
6. **Generating Pull requests :**
	- [Generate a pull request](https://help.github.com/articles/about-pull-requests/) from your ```branch``` branch to ```master``` branch.
	- Give the PR and apt title, and mention `Fixes #<issue_number>` in the comment to link it with the issue.
	- Don't close the issue by your own.
7. **Commenting your Code**
	- Include your comments directly preceding an object for GoDoc to document it.
	- Indent pre-formatted comments.
	- Refer to the [Guidelines](https://blog.golang.org/godoc-documenting-go-code) for more info on commenting.

### Resources and References
- [Bittorent Specifications](http://jonas.nitro.dk/bittorrent/bittorrent-rfc.html)
- [Bittorent Specifications](http://www.bittorrent.org/beps/bep_0003.html)
- [Bittorrent in Javascript](https://allenkim67.github.io/programming/2016/05/04/how-to-make-your-own-bittorrent-client.html)
- [Bittorrent in C#](https://www.seanjoflynn.com/research/bittorrent.html)
- [Complete Bittorrent in Go](https://github.com/jackpal/Taipei-Torrent)
- [Network Programming in Go](https://ipfs.io/ipfs/QmfYeDhGH9bZzihBUDEQbCbTc5k5FZKURMUoUvfmc27BwL/index.html)
- [Concurrency in Go](https://github.com/golang/go/wiki/LearnConcurrency)
- [Go advanced testing tips & tricks](https://medium.com/@povilasve/go-advanced-tips-tricks-a872503ac859)
- [Travis CI tutorial](https://docs.travis-ci.com/user/tutorial/)
